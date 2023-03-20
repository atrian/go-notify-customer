package workers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/atrian/go-notify-customer/internal/dto"
	"github.com/atrian/go-notify-customer/internal/interfaces"
	"github.com/atrian/go-notify-customer/internal/workers/channelServices"
)

var (
	_ interfaces.Worker = (*ChannelWorker)(nil)
)

type ChannelWorker struct {
	ctx          context.Context
	mu           sync.Locker
	config       config
	services     map[string]channelService
	sendStatChan chan<- dto.Stat
	client       interfaces.AmpqClient
	logger       interfaces.Logger
}

// loadServices загрузки стандартных сервисов отправки.
// Защищено через sync.Locker
func (c *ChannelWorker) loadServices() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.services = map[string]channelService{
		"sms":  channelServices.NewTwilio(c.ctx, c.config),
		"mail": channelServices.NewMail(c.ctx, c.config),
	}
}

// ReloadService для добавления новых сервисов отправки на лету или для подмены
// сервисов в тестах на заглушки. Защищено через sync.Locker
func (c *ChannelWorker) ReloadService(overwrite string, service channelService) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.services == nil {
		c.services = make(map[string]channelService)
	}

	c.services[overwrite] = service
}

func NewChannelWorker(ctx context.Context, conf config, client interfaces.AmpqClient, sendStatChan chan<- dto.Stat, logger interfaces.Logger) *ChannelWorker {
	w := ChannelWorker{
		mu:           &sync.Mutex{},
		config:       conf,
		client:       client,
		sendStatChan: sendStatChan,
		ctx:          ctx,
		logger:       logger,
	}

	w.loadServices()

	return &w
}

type channelService interface {
	SendMessage(message string, destination string) error
}

type config interface {
	GetAmpqDSN() string
	GetNotificationQueue() string
	GetFailedWorksQueue() string
	mailConfig
	twilioConfig
}

type mailConfig interface {
	GetMailSenderAddress() string
	GetMailSMTPHost() string
	GetMailLogin() string
	GetMailPassword() string
	GetMailMessageTheme() string
}

type twilioConfig interface {
	GetTwilioAccountSid() string
	GetTwilioAuthToken() string
	GetTwilioSenderPhone() string
}

// Start потребляет очередь consumeQueue, восстанавливает объект dto.Message из json
// и отправляет его в ChannelWorker.Send
func (c *ChannelWorker) Start(consumeQueue string, successQueue string, failQueue string) {
	c.client.MigrateDurableQueues(consumeQueue, successQueue, failQueue)

	msgs, err := c.client.Consume(consumeQueue)
	if err != nil {
		c.logger.Error("Can't consume message queue", err)
	}

	go func() {
		for d := range msgs {
			var message dto.Message
			jsonErr := json.Unmarshal(d.Body, &message)

			if jsonErr != nil {
				c.logger.Error("ChannelWorker start JSON unmarshall err", jsonErr)
			}

			c.logger.Info(fmt.Sprintf("Received a message from BUS notificationUUID:%v personUUID: %v, text: %v", message.NotificationUUID, message.PersonUUID, message.Text))

			c.Send(message)
		}
	}()

	<-c.ctx.Done()
}

// Send принимает сообщение в формате dto.Message и отправляет его в нужный сервис
// в случае ошибки пишет в канал статистики через ChannelWorker.sendStat
func (c *ChannelWorker) Send(message dto.Message) {
	service, exist := c.services[message.Channel]

	if !exist {
		c.logger.Error(fmt.Sprintf("Bad channel: %v for notificationUUID:%v", message.Channel, message.NotificationUUID), errors.New("not exist"))
		c.sendStat(message, dto.BadChannel)
		return
	}

	err := service.SendMessage(message.Text, message.DestinationAddress)

	if err != nil {
		c.logger.Error(fmt.Sprintf("External sender error for notificationUUID:%v", message.NotificationUUID), err)
		c.sendStat(message, dto.Failed)
		return
	}

	c.logger.Info(fmt.Sprintf("Notification SENT notificationUUID:%v", message.NotificationUUID))
	c.sendStat(message, dto.Sent)
}

// Stop остановка воркера
func (c *ChannelWorker) Stop() {
	c.logger.Info("Channel worker stopped")
}

// sendStat отправка статистики в формате dto.Stat в канал sendStatChan
// канал слушает stat.Service
func (c *ChannelWorker) sendStat(message dto.Message, status dto.StatStatus) {
	// отправляем статистику
	c.sendStatChan <- dto.Stat{
		PersonUUID:       message.PersonUUID,
		NotificationUUID: message.NotificationUUID,
		CreatedAt:        time.Now().Format(time.RFC3339),
		Status:           status,
	}
}
