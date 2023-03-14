package interfaces

import (
	"context"

	"github.com/atrian/go-notify-customer/internal/services/stat/entity"
)

// StatService Интерфейс сервиса хранения данных о бизнес событиях
// (прим.: новая запись, перенос записи, отмена записи). Ответственность: обслуживание REST CRUD
type StatService interface {
	// BaseService Общий сервисный интерфейс с методами Start и Stop
	BaseService

	// All FindByPersonUUID FindByEventUUID - выдача статистики в разрезах
	All(ctx context.Context) []entity.Stat
	Store(statChan chan entity.Stat) error // Store сохраняет статистику переданную через канал
	FindByPersonUUID(ctx context.Context, orderId string) ([]entity.Stat, error)
	FindByEventUUID(ctx context.Context, orderId string) ([]entity.Stat, error)
}