package interfaces

import (
	"context"

	"github.com/atrian/go-notify-customer/internal/dto"

	"github.com/google/uuid"
)

// TemplateService Интерфейс сервиса хранения данных о шаблонах уведомлений
// Ответственность: обслуживание REST CRUD
type TemplateService interface {
	// BaseService Общий сервисный интерфейс с методами Start и Stop
	BaseService

	// All Store StoreBatch Update FindBy DeleteById - CRUD
	All(ctx context.Context) []dto.Template
	Store(ctx context.Context, template dto.Template) (dto.Template, error)
	StoreBatch(ctx context.Context, templates []dto.Template) ([]dto.Template, error)
	Update(ctx context.Context, template dto.Template) (dto.Template, error)
	FindById(ctx context.Context, templateUUID uuid.UUID) (dto.Template, error)
	FindByEventId(ctx context.Context, eventUUID uuid.UUID) (dto.Template, error)
	DeleteById(ctx context.Context, templateUUID uuid.UUID) error
}
