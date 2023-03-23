package interfaces

import (
	"context"

	"github.com/atrian/go-notify-customer/internal/dto"

	"github.com/google/uuid"
)

// StatService Интерфейс сервиса хранения данных о бизнес событиях
// (прим.: новая запись, перенос записи, отмена записи). Ответственность: обслуживание REST CRUD
type StatService interface {
	// BaseService Общий сервисный интерфейс с методами Start и Stop
	BaseService

	// All FindByPersonUUID FindByEventUUID - выдача статистики в разрезах
	All(ctx context.Context) []dto.Stat
	Store(ctx context.Context, stat dto.Stat) error
	FindByPersonUUID(ctx context.Context, personUUID uuid.UUID) ([]dto.Stat, error)
	FindByNotificationId(ctx context.Context, notificationUUID uuid.UUID) ([]dto.Stat, error)
}
