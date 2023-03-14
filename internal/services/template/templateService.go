package template

import (
	"context"

	"github.com/google/uuid"

	"github.com/atrian/go-notify-customer/internal/services/template/entity"
)

type Service struct {
}

func (s Service) Start() {
	//TODO implement me
	panic("implement me")
}

func (s Service) Stop() {
	//TODO implement me
	panic("implement me")
}

func (s Service) All(ctx context.Context) []entity.Template {
	//TODO implement me
	panic("implement me")
}

func (s Service) Store(ctx context.Context, template entity.Template) (entity.Template, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) StoreBatch(ctx context.Context, templates []entity.Template) ([]entity.Template, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) Update(ctx context.Context, template entity.Template) (entity.Template, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) FindById(ctx context.Context, templateUUID uuid.UUID) (entity.Template, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) DeleteById(ctx context.Context, templateUUID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
