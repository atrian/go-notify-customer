package stat

import (
	"context"

	"github.com/atrian/go-notify-customer/internal/services/stat/entity"

	"github.com/google/uuid"
)

type Storager interface {
	All(ctx context.Context) ([]entity.Stat, error)
	Store(ctx context.Context, stat entity.Stat) error
	GetByNotificationId(ctx context.Context, notificationUUID uuid.UUID) (entity.Stat, error)
	GetByPersonId(ctx context.Context, personUUID uuid.UUID) ([]entity.Stat, error)
}
