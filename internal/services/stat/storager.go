package stat

import (
	"context"

	"github.com/google/uuid"

	"github.com/atrian/go-notify-customer/internal/dto"
)

// Storager интерфейс хранилища сервиса stat
type Storager interface {
	// All возвращает все записи
	All(ctx context.Context) ([]dto.Stat, error)
	// Store созраняет запись в хранилище
	Store(ctx context.Context, stat dto.Stat) error
	// GetByNotificationId возвращает записи по uuid уведомления
	GetByNotificationId(ctx context.Context, notificationUUID uuid.UUID) ([]dto.Stat, error)
	// GetByPersonId возвращает записи по uuid получателя уведомления
	GetByPersonId(ctx context.Context, personUUID uuid.UUID) ([]dto.Stat, error)
}
