package interfaces

import (
	"context"

	"github.com/atrian/go-notify-customer/internal/dto"

	"github.com/google/uuid"
)

// EventService Интерфейс сервиса хранения данных о бизнес событиях
// (прим.: новая запись, перенос записи, отмена записи). Ответственность: обслуживание REST CRUD
type EventService interface {
	// BaseService Общий сервисный интерфейс с методами Start и Stop
	BaseService

	// All Store StoreBatch Update FindBy DeleteById - CRUD
	All(ctx context.Context) []dto.Event
	Store(ctx context.Context, event dto.Event) (dto.Event, error)
	StoreBatch(ctx context.Context, events []dto.Event) ([]dto.Event, error)
	Update(ctx context.Context, event dto.Event) (dto.Event, error)
	FindById(ctx context.Context, eventUUID uuid.UUID) (dto.Event, error)
	DeleteById(ctx context.Context, eventUUID uuid.UUID) error
}
