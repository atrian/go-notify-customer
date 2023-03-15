package event

import (
	"context"

	"github.com/google/uuid"

	"github.com/atrian/go-notify-customer/internal/services/event/entity"
)

type Storager interface {
	All(ctx context.Context) ([]entity.Event, error)
	Store(ctx context.Context, event entity.Event) error
	Update(ctx context.Context, event entity.Event) error
	GetById(ctx context.Context, eventUUID uuid.UUID) (entity.Event, error)
	DeleteById(ctx context.Context, eventUUID uuid.UUID) error
}
