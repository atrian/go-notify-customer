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
	prepareTemplate(template string, replaces []dto.MessageParam) string
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
	contacts, err := d.services.getContacts(d.ctx, notification.PersonUUIDs)
	if err != nil {
		// TODO err handle
	}

	// запрос бизнес события
	event, err := d.services.getEvent(d.ctx, notification.EventUUID)
	if err != nil {
		// TODO err handle
	}

	// запрос шаблона для события
	templates, err := d.services.getTemplates(d.ctx, notification.EventUUID)
	if err != nil {
		// TODO err handle
	}

	// Формируем доступные шаблоны - делаем подстановки параметров в текст
	// структура preparedTemplates [тип_канала]текст_с_подстановками
	preparedTemplates := make(map[string]string, len(templates))

	for _, template := range templates {
		preparedTemplates[template.ChannelType] = d.services.prepareTemplate(template.Body, notification.MessageParams)
	}

	// для каждого канала в котором должно быть уведомление
	for _, notificationChannel := range event.NotificationChannels {

		// проверяем что есть шаблон
		template, exist := preparedTemplates[notificationChannel]
		if !exist {
			d.logger.Info(fmt.Sprintf("Template does not exist for channel: %v, event: %v", notificationChannel, event.EventUUID))
			// если шаблона нет слать нечего, пропускаем канал
			continue
		}

		// Для каждого пользователя берем нужный контакт
		for _, contact := range contacts {
			// выбор контакта для канала
			relatedContact, cErr := contactLocator(notificationChannel, contact)
			if cErr != nil {
				// TODO не нашли контакт человека для этого канала, пишем лог
				continue
			}

			// и добавляем сообщение в слайс на отправку
			messages = append(messages, dto.Message{
				PersonUUID:         uuid.UUID{},
				Text:               template,
				Channel:            notificationChannel,
				DestinationAddress: relatedContact.Destination,
			})
		}
	}

	return messages
}

// contactLocator Выбирает адрес назначения (телефон, емейл, и пр) для определенного канала
func contactLocator(notificationChannel string, contacts dto.PersonContacts) (dto.Contact, error) {
	for _, contact := range contacts.Contacts {
		if notificationChannel == contact.Channel {
			return contact, nil
		}
	}

	return dto.Contact{}, fmt.Errorf("NotFound")
}
