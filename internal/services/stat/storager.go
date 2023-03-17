package stat

import (
	"context"

	"github.com/atrian/go-notify-customer/internal/dto"

	"github.com/google/uuid"
)

type Storager interface {
	All(ctx context.Context) ([]dto.Stat, error)
	Store(ctx context.Context, stat dto.Stat) error
	GetByNotificationId(ctx context.Context, notificationUUID uuid.UUID) (dto.Stat, error)
	GetByPersonId(ctx context.Context, personUUID uuid.UUID) ([]dto.Stat, error)
}
