package event

import (
	"context"

	"github.com/google/uuid"

	"github.com/atrian/go-notify-customer/internal/services/event/entity"
)

type Service struct {
	storage Storager
}

func (e Service) Start() {
	// подключаем in memory хранилище
	e.storage = NewMemoryStorage()
}

func (e Service) Stop() {
	//TODO implement me
	panic("implement me")
}

func (e Service) All(ctx context.Context) []entity.Event {
	events, err := e.storage.All(ctx)

	if err != nil {
		// TODO log err
	}

	return events
}

func (e Service) Store(ctx context.Context, event entity.Event) (entity.Event, error) {
	event.EventUUID = uuid.New()

	err := e.storage.Store(ctx, event)
	if err != nil {
		return entity.Event{}, err
	}

	return event, nil
}

func (e Service) StoreBatch(ctx context.Context, events []entity.Event) ([]entity.Event, error) {
	for i := 0; i < len(events); i++ {
		events[i].EventUUID = uuid.New()
		err := e.storage.Store(ctx, events[i])
		if err != nil {
			// TODO err handling
			// Удалить все добавленные и вернуть ошибку?
		}
	}

	return events, nil
}

func (e Service) Update(ctx context.Context, event entity.Event) (entity.Event, error) {
	err := e.storage.Store(ctx, event)
	if err != nil {
		return entity.Event{}, err
	}

	return event, nil
}

func (e Service) FindById(ctx context.Context, eventUUID uuid.UUID) (entity.Event, error) {
	return e.storage.GetById(ctx, eventUUID)
}

func (e Service) DeleteById(ctx context.Context, eventUUID uuid.UUID) error {
	return e.storage.DeleteById(ctx, eventUUID)
}
