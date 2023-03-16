package template

import (
	"context"
	"github.com/atrian/go-notify-customer/internal/services/template/entity"

	"github.com/google/uuid"
)

type Storager interface {
	All(ctx context.Context) ([]entity.Template, error)
	Store(ctx context.Context, template entity.Template) error
	Update(ctx context.Context, template entity.Template) error
	GetById(ctx context.Context, templateUUID uuid.UUID) (entity.Template, error)
	GetByEventId(ctx context.Context, eventUUID uuid.UUID) (entity.Template, error)
	DeleteById(ctx context.Context, templateUUID uuid.UUID) error
}
