package event

import (
	"context"

	"github.com/google/uuid"

	"github.com/atrian/go-notify-customer/internal/dto"
)

// Storager интерфейс хранилища сервиса event
type Storager interface {
	// All возвращает все записи
	All(ctx context.Context) ([]dto.Event, error)
	// Store созраняет запись в хранилище
	Store(ctx context.Context, event dto.Event) error
	// Update обновляет запись в хранилище
	Update(ctx context.Context, event dto.Event) error
	// GetById возвращает запись по uuid сущности
	GetById(ctx context.Context, eventUUID uuid.UUID) (dto.Event, error)
	// DeleteById удаляет запись по uuid сущности
	DeleteById(ctx context.Context, eventUUID uuid.UUID) error
}
