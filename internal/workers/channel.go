package workers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/atrian/go-notify-customer/internal/workers/channelServices"
	"time"

	"github.com/atrian/go-notify-customer/internal/dto"
	"github.com/atrian/go-notify-customer/internal/interfaces"
)

var _ interfaces.Worker = (*ChannelWorker)(nil)

type ChannelWorker struct {
	services     map[string]channelService
	config       config
	client       interfaces.AmpqClient
	sendStatChan chan<- dto.Stat
	ctx          context.Context
	logger       interfaces.Logger
}

func (c *ChannelWorker) loadServices() {
	c.services = map[string]channelService{
		"sms":  channelServices.NewTwilio(),
		"mail": channelServices.NewMail(),
	}
}

func (c *ChannelWorker) ReloadService(overwrite string, service channelService) {
	if c.services == nil {
		c.services = make(map[string]channelService)
	}

	c.services[overwrite] = service
}

func NewChannelWorker(conf config, client interfaces.AmpqClient, sendStatChan chan<- dto.Stat, ctx context.Context, logger interfaces.Logger) *ChannelWorker {
	w := ChannelWorker{
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
	GetMailConfig()   // TODO
	GetTwilioConfig() // TODO
}

func (c *ChannelWorker) Start(consumeQueue string, successQueue string, failQueue string) {
	c.client.MigrateDurableQueues(consumeQueue)

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
