// Package event Сервис работы с бизнес событиями для нужд
// REST CRUD и формирования сообщений
// при отправке уведомления во входящих данных указывается EventUUID
// по которому находится шаблон сообщения для разных каналов доставки
//
// Формат передачи между слоями приложения dto.Event
// Формат входящего json см. dto.IncomingEvent
//
//	type IncomingNotification struct {
//			EventUUID     uuid.UUID      `json:"event_uuid"`               // uuid бизнес события
//			PersonUUIDs   []uuid.UUID    `json:"person_uuids"`             //
//			MessageParams []MessageParam `json:"message_params,omitempty"` //
//			Priority      uint           `json:"priority,omitempty"`       //
//		}
//
// Подробнее см. dto.IncomingNotification
package event

import (
	"context"

	"github.com/google/uuid"

	"github.com/atrian/go-notify-customer/internal/dto"
	"github.com/atrian/go-notify-customer/internal/interfaces"
)

// Service структура сервиса бизнес событий содержит in-mem хранилище с интерфейсом:
//
//	type Storager interface {
//			All(ctx context.Context) ([]dto.Event, error)
//			Store(ctx context.Context, event dto.Event) error
//			Update(ctx context.Context, event dto.Event) error
//			GetById(ctx context.Context, eventUUID uuid.UUID) (dto.Event, error)
//			DeleteById(ctx context.Context, eventUUID uuid.UUID) error
//		}
//
// и логгер с интерфейсом interfaces.Logger
type Service struct {
	storage Storager
	logger  interfaces.Logger
}

// New при создании требует логгер удовлетворяющий интерфейсу interfaces.Logger
// В приложении используется реализация с Zap
func New(logger interfaces.Logger) *Service {
	s := Service{
		logger:  logger,
		storage: NewMemoryStorage(),
	}
	return &s
}

// Start стартовые процедуры для сервиса
func (e Service) Start(ctx context.Context) {
	e.logger.Info("Event service started")
}

// Stop завершение работы сервиса grace shutdown
func (e Service) Stop() {
	e.logger.Info("Event service stopped")
}

// All возвращает все хранящиеся бизнес события слайсом dto.Event
func (e Service) All(ctx context.Context) []dto.Event {
	events, err := e.storage.All(ctx)

	if err != nil {
		e.logger.Error("Event service storage.All err", err)
	}

	if events == nil {
		return []dto.Event{}
	}

	return events
}

// Store созраняет dto.Event в хранилище. Событию присваивается UUID
func (e Service) Store(ctx context.Context, event dto.Event) (dto.Event, error) {
	event.EventUUID = uuid.New()

	err := e.storage.Store(ctx, event)
	if err != nil {
		e.logger.Error("Event service storage.Store err", err)
		return dto.Event{}, err
	}

	return event, nil
}

// StoreBatch массовое сохранение событий в хранилище. Событиям присваивается UUID
//
// BUG(Timur) если при возникновении отдельного события возникнут ошибки весь батч не будет отклонен
// поведение будет исправлено в следующих версиях
func (e Service) StoreBatch(ctx context.Context, events []dto.Event) ([]dto.Event, error) {
	for i := 0; i < len(events); i++ {
		events[i].EventUUID = uuid.New()
		err := e.storage.Store(ctx, events[i])
		if err != nil {
			e.logger.Error("Event service storage.Store err", err)
		}
	}

	return events, nil
}

// Update обновляет бизнес событие в хранилище
func (e Service) Update(ctx context.Context, event dto.Event) (dto.Event, error) {
	err := e.storage.Update(ctx, event)
	if err != nil {
		e.logger.Error("Event service storage.Update err", err)
		return dto.Event{}, err
	}

	return event, nil
}

// FindById поиск события по uuid
func (e Service) FindById(ctx context.Context, eventUUID uuid.UUID) (dto.Event, error) {
	return e.storage.GetById(ctx, eventUUID)
}

// DeleteById удаление события по uuid
func (e Service) DeleteById(ctx context.Context, eventUUID uuid.UUID) error {
	return e.storage.DeleteById(ctx, eventUUID)
}
