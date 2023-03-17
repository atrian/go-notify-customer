package notificationDispatcher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/atrian/go-notify-customer/internal/dto"
	"github.com/atrian/go-notify-customer/internal/interfaces"
	"github.com/atrian/go-notify-customer/pkg/logger"
)

type serviceGateway interface {
	getContacts(ctx context.Context, personUUIDs []uuid.UUID) ([]dto.PersonContacts, error)
	getTemplates(ctx context.Context, eventUUID uuid.UUID) ([]dto.Template, error)
	getEvent(ctx context.Context, eventUuid uuid.UUID) (dto.Event, error)
	parseTemplate(template string, replaces map[string]string) string
}

type dispatcherConfig interface {
	GetMessageDispatchQueue() string
}

type Dispatcher struct {
	ctx              context.Context
	notificationChan <-chan dto.Notification
	config           dispatcherConfig
	services         serviceGateway
	ampqClient       interfaces.AmpqClient
	logger           logger.Logger
}

func New(
	ctx context.Context,
	notificationChan chan dto.Notification,
	config dispatcherConfig,
	serviceGateway serviceGateway,
	ampqClient interfaces.AmpqClient,
	logger logger.Logger) *Dispatcher {

	d := Dispatcher{
		ctx:              ctx,
		notificationChan: notificationChan,
		config:           config,
		services:         serviceGateway,
		ampqClient:       ampqClient,
		logger:           logger,
	}

	return &d
}

// Start стартовые операции для notificationDispatcher - ampq миграция,
// запуск прослушивания канала
func (d Dispatcher) Start() {
	// миграция AMPQ очередей
	d.ampqClient.MigrateDurableQueues(d.config.GetMessageDispatchQueue())

	// слушаем канал с уведомлениями, строим сообщения и отправляем на исполнение
	go d.listenInputChannel(d.ctx, d.notificationChan)
}

func (d Dispatcher) Stop() {
	d.ampqClient.Stop()
}

func (d Dispatcher) dispatch(message dto.Message) error {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		d.logger.Error("Message JSON marshal failed", err)
		return err
	}

	err = d.ampqClient.Publish(d.config.GetMessageDispatchQueue(), jsonMessage)
	if err != nil {
		d.logger.Error("Message JSON marshal failed", err)
		return err
	}

	infoMessage := fmt.Sprintf("Message dispatched for %v", message.PersonUUID.String())
	d.logger.Info(infoMessage)

	return nil
}

// listenInputChannel слушает канал с уведомлениями dto.Notification
// и размещает Dispatch сообщения dto.Message для дальнейшей отправки
func (d Dispatcher) listenInputChannel(ctx context.Context, input <-chan dto.Notification) {
	for {
		select {
		case notification := <-input:
			messages := d.buildMessages(notification)

			for i := 0; i < len(messages); i++ {
				dErr := d.dispatch(messages[i])

				if dErr != nil {
					d.logger.Error("Message dispatch error", dErr)
				}
			}

		case <-ctx.Done():
			return
		default:
			// do nothing
		}
	}
}

// buildMessages формирует клиентские сообщения из уведомления, шаблона и контактов
func (d Dispatcher) buildMessages(notification dto.Notification) []dto.Message {
	var messages []dto.Message

	// запрос контактов
	contacts, _ := d.services.getContacts(d.ctx, notification.PersonUUIDs)
	// запрос бизнес события
	event, _ := d.services.getEvent(d.ctx, notification.EventUUID)
	// запрос шаблона для события
	templates, _ := d.services.getTemplates(d.ctx, notification.EventUUID)

	_, _, _ = contacts, event, templates
	// сформировать сообщение
	// добавить его в слайс
	// вернуть слайс

	return messages
}
