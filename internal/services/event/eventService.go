package event

import (
	"context"

	"github.com/atrian/go-notify-customer/internal/dto"

	"github.com/google/uuid"
)

type Service struct {
	storage Storager
}

func New() *Service {
	s := Service{
		storage: NewMemoryStorage(),
	}
	return &s
}

func (e Service) Start() {
	// log?
}

func (e Service) Stop() {
	//TODO implement me
	panic("implement me")
}

func (e Service) All(ctx context.Context) []dto.Event {
	events, err := e.storage.All(ctx)

	if err != nil {
		// TODO log err
	}

	return events
}

func (e Service) Store(ctx context.Context, event dto.Event) (dto.Event, error) {
	event.EventUUID = uuid.New()

	err := e.storage.Store(ctx, event)
	if err != nil {
		return dto.Event{}, err
	}

	return event, nil
}

func (e Service) StoreBatch(ctx context.Context, events []dto.Event) ([]dto.Event, error) {
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

func (e Service) Update(ctx context.Context, event dto.Event) (dto.Event, error) {
	err := e.storage.Store(ctx, event)
	if err != nil {
		return dto.Event{}, err
	}

	return event, nil
}

func (e Service) FindById(ctx context.Context, eventUUID uuid.UUID) (dto.Event, error) {
	return e.storage.GetById(ctx, eventUUID)
}

func (e Service) DeleteById(ctx context.Context, eventUUID uuid.UUID) error {
	return e.storage.DeleteById(ctx, eventUUID)
}
