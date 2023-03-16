package notificationDispatcher

import (
	"context"
	"github.com/atrian/go-notify-customer/config"
	"github.com/atrian/go-notify-customer/internal/interfaces"
	"github.com/atrian/go-notify-customer/internal/services/notificationDispatcher/entity"
	notifyEntity "github.com/atrian/go-notify-customer/internal/services/notify/entity"
	"github.com/atrian/go-notify-customer/pkg/ampq"
	"github.com/atrian/go-notify-customer/pkg/logger"
)

type Dispatcher struct {
	ctx              context.Context
	notificationChan chan notifyEntity.Notification
	config           config.QueueConfig
	ampqClient       interfaces.AmpqClient
	logger           logger.Logger
}

func New(ctx context.Context, config config.Config, notificationChan chan notifyEntity.Notification, logger logger.Logger) *Dispatcher {
	d := Dispatcher{
		ctx:              ctx,
		notificationChan: notificationChan,
		config:           config.NotificationQueue,
		ampqClient:       ampq.New(config.AMPQDSN, logger),
		logger:           logger,
	}

	return &d
}

func (d Dispatcher) Start() {
	// запуск ampq миграций сервиса notificationDispatcher
	d.ampqClient.MigrateDurableQueues(d.config.DispatchQueue, d.config.ListenQueue)

	go func(ctx context.Context, input <-chan notifyEntity.Notification) {
		for {
			select {
			case notification := <-input:
				messages := d.buildMessages(notification)
				
				for i := 0; i < len(messages); i++ {
					err := d.Dispatch(messages[i])

					if err != nil {
						d.logger.Error("Message dispatch error", err)
					}
				}

			case <-ctx.Done():
				// TODO shutdown
				return

			default:
				// do nothing
			}
		}
	}(d.ctx, d.notificationChan)
}

func (d Dispatcher) Stop() {
	d.ampqClient.Stop()
}

func (d Dispatcher) Dispatch(message entity.Message) error {
	//TODO implement me
	panic("implement me")
}

func (d Dispatcher) buildMessages(notification notifyEntity.Notification) []entity.Message {
	var messages []entity.Message

	// запросить контакты
	// запросить шаблон
	// запросить бизнес событие
	// сформировать сообщение
	// добавить его в слайс
	// вернуть слайс

	return messages
}
