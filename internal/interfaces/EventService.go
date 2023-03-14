package interfaces

import (
	"context"

	"github.com/google/uuid"

	"github.com/atrian/go-notify-customer/internal/services/event/entity"
)

// EventService Интерфейс сервиса хранения данных о бизнес событиях
// (прим.: новая запись, перенос записи, отмена записи). Ответственность: обслуживание REST CRUD
type EventService interface {
	// BaseService Общий сервисный интерфейс с методами Start и Stop
	BaseService

	// All Store StoreBatch Update FindBy DeleteById - CRUD
	All(ctx context.Context) []entity.Event
	Store(ctx context.Context, event entity.Event) (entity.Event, error)
	StoreBatch(ctx context.Context, events []entity.Event) ([]entity.Event, error)
	Update(ctx context.Context, event entity.Event) (entity.Event, error)
	FindById(ctx context.Context, eventUUID uuid.UUID) (entity.Event, error)
	DeleteById(ctx context.Context, eventUUID uuid.UUID) error
}
