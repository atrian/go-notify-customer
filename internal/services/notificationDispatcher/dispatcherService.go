// Package notificationDispatcher диспетчер подготовки уведомлений к отправке
// обращается к другим сервисам зя дополнительной информацией по уведомлению через фасад serviceGateway
// результат работы отправляет во внешнюю очередь (RabbitMQ)
package notificationDispatcher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/atrian/go-notify-customer/internal/dto"
	"github.com/atrian/go-notify-customer/internal/interfaces"
)

// serviceGateway интерфейс сервисного фасада, ограничение на досупные методы автономных сервисов.
type serviceGateway interface {
	// getContacts запрос контактов во внешнем защищенном vault. gRPC
	getContacts(ctx context.Context, personUUIDs []uuid.UUID) ([]dto.PersonContacts, error)
	// getTemplates запрос шаблона сообщения для бизнес события
	getTemplates(ctx context.Context, eventUUID uuid.UUID) ([]dto.Template, error)
	// getEvent запрос деталей бизнес события
	getEvent(ctx context.Context, eventUuid uuid.UUID) (dto.Event, error)
	// prepareTemplate выполнение именованных подстановок в шаблоне сообщения
	prepareTemplate(template string, replaces []dto.MessageParam) string
}

// dispatcherConfig интерфейс кинфигурации доступной сервису notificationDispatcher
type dispatcherConfig interface {
	GetAmpqDSN() string
	GetNotificationQueue() string
}

// Dispatcher содержит канал по котороку получает входящие уведомлений
// конфигурацию, фасад с нужными для работы сервисами,
// ampq клиент с интерфейсом interfaces.AmpqClient
// и логгер с интерфейсом interfaces.Logger
type Dispatcher struct {
	notificationChan <-chan dto.Notification
	config           dispatcherConfig
	services         serviceGateway
	ampqClient       interfaces.AmpqClient
	logger           interfaces.Logger
}

func New(
	notificationChan chan dto.Notification,
	config dispatcherConfig,
	serviceGateway serviceGateway,
	ampqClient interfaces.AmpqClient,
	logger interfaces.Logger) *Dispatcher {

	d := Dispatcher{
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
func (d Dispatcher) Start(ctx context.Context) {
	d.logger.Info("Notification dispatcher started")

	// Подключаем ampq клиент
	err := d.ampqClient.Connect(d.config.GetAmpqDSN())
	if err != nil {
		d.logger.Error("Dispatcher ampqClient.Connect err", err)
	}
	// миграция AMPQ очередей
	d.ampqClient.MigrateDurableQueues(d.config.GetNotificationQueue())

	// слушаем канал с уведомлениями, строим сообщения и отправляем на исполнение
	go d.listenInputChannel(ctx, d.notificationChan)
}

// Stop корректное завершение работы
func (d Dispatcher) Stop() {
	d.ampqClient.Stop()
	d.logger.Info("Notification dispatcher stopped")
}

// dispatch отправка готового сообщения в очередь для исполнения
// channel worker'ом.
func (d Dispatcher) dispatch(message dto.Message) error {
	// Готовим json и публикуем в очередь
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		d.logger.Error("Message JSON marshal failed", err)
		return err
	}

	err = d.ampqClient.Publish(d.config.GetNotificationQueue(), jsonMessage)
	if err != nil {
		d.logger.Error("ampqClient.Publish error", err)
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
			// выдаем UUID уведомлению
			notification.NotificationUUID = uuid.New()

			d.logger.Debug(fmt.Sprintf("Notification received: %v", notification))

			// собираем сообщения
			messages := d.buildMessages(ctx, notification)

			// размещаем сообщения во внешнюю очередь на отправку
			for i := 0; i < len(messages); i++ {
				dErr := d.dispatch(messages[i])

				if dErr != nil {
					d.logger.Error("Message dispatch error", dErr)
				}
			}

		case <-ctx.Done(): // отбой по контексту
			return
		default:
			// do nothing
		}
	}
}

// buildMessages формирует клиентские сообщения из уведомления, шаблона и контактов
func (d Dispatcher) buildMessages(ctx context.Context, notification dto.Notification) []dto.Message {
	var messages []dto.Message

	// запрос контактов
	contacts, err := d.services.getContacts(ctx, notification.PersonUUIDs)
	if err != nil {
		d.logger.Error("Dispatcher getContacts err", err)
	}

	// запрос бизнес события
	event, err := d.services.getEvent(ctx, notification.EventUUID)
	if err != nil {
		d.logger.Error("Dispatcher getEvent err", err)
	}

	// запрос шаблона для события
	templates, err := d.services.getTemplates(ctx, notification.EventUUID)
	if err != nil {
		d.logger.Error("Dispatcher getTemplates err", err)
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
				d.logger.Error("Dispatcher contactLocator err", cErr)
				continue
			}

			d.logger.Debug("Message prepared")

			// и добавляем сообщение в слайс на отправку
			messages = append(messages, dto.Message{
				PersonUUID:         contact.PersonUUID,
				NotificationUUID:   notification.NotificationUUID,
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
