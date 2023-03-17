package event

import (
	"context"

	"github.com/atrian/go-notify-customer/internal/dto"

	"github.com/google/uuid"
)

type Storager interface {
	All(ctx context.Context) ([]dto.Event, error)
	Store(ctx context.Context, event dto.Event) error
	Update(ctx context.Context, event dto.Event) error
	GetById(ctx context.Context, eventUUID uuid.UUID) (dto.Event, error)
	DeleteById(ctx context.Context, eventUUID uuid.UUID) error
}
