package event

import (
	"context"

	"github.com/google/uuid"

	"github.com/atrian/go-notify-customer/internal/services/event/entity"
)

type Service struct {
}

func (e Service) Start() {
	//TODO implement me
	panic("implement me")
}

func (e Service) Stop() {
	//TODO implement me
	panic("implement me")
}

func (e Service) All(ctx context.Context) []entity.Event {
	//TODO implement me
	panic("implement me")
}

func (e Service) Store(ctx context.Context, event entity.Event) (entity.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (e Service) StoreBatch(ctx context.Context, events []entity.Event) ([]entity.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (e Service) Update(ctx context.Context, event entity.Event) (entity.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (e Service) FindById(ctx context.Context, eventUUID uuid.UUID) (entity.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (e Service) DeleteById(ctx context.Context, eventUUID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
