package interfaces

import (
	"context"

	"github.com/google/uuid"

	"github.com/atrian/go-notify-customer/internal/services/template/entity"
)

// TemplateService Интерфейс сервиса хранения данных о шаблонах уведомлений
// Ответственность: обслуживание REST CRUD
type TemplateService interface {
	// BaseService Общий сервисный интерфейс с методами Start и Stop
	BaseService

	// All Store StoreBatch Update FindBy DeleteById - CRUD
	All(ctx context.Context) []entity.Template
	Store(ctx context.Context, template entity.Template) (entity.Template, error)
	StoreBatch(ctx context.Context, templates []entity.Template) ([]entity.Template, error)
	Update(ctx context.Context, template entity.Template) (entity.Template, error)
	FindById(ctx context.Context, templateUUID uuid.UUID) (entity.Template, error)
	FindByEventId(ctx context.Context, eventUUID uuid.UUID) (entity.Template, error)
	DeleteById(ctx context.Context, templateUUID uuid.UUID) error
}
