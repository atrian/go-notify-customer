package template

import (
	"context"

	"github.com/google/uuid"

	"github.com/atrian/go-notify-customer/internal/dto"
)

// Storager интерфейс хранилища сервиса template
type Storager interface {
	// All возвращает все записи
	All(ctx context.Context) ([]dto.Template, error)
	// Store созраняет запись в хранилище
	Store(ctx context.Context, template dto.Template) error
	// Update обновляет запись в хранилище
	Update(ctx context.Context, template dto.Template) error
	// GetById возвращает запись по uuid сущности
	GetById(ctx context.Context, templateUUID uuid.UUID) (dto.Template, error)
	// GetByEventId возвращает записи по uuid бизнес события
	GetByEventId(ctx context.Context, eventUUID uuid.UUID) ([]dto.Template, error)
	// DeleteById удаляет запись по uuid сущности
	DeleteById(ctx context.Context, templateUUID uuid.UUID) error
}
