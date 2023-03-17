package template

import (
	"context"

	"github.com/google/uuid"

	"github.com/atrian/go-notify-customer/internal/dto"
)

type Storager interface {
	All(ctx context.Context) ([]dto.Template, error)
	Store(ctx context.Context, template dto.Template) error
	Update(ctx context.Context, template dto.Template) error
	GetById(ctx context.Context, templateUUID uuid.UUID) (dto.Template, error)
	GetByEventId(ctx context.Context, eventUUID uuid.UUID) ([]dto.Template, error)
	DeleteById(ctx context.Context, templateUUID uuid.UUID) error
}
