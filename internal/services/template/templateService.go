package template

import (
	"context"

	"github.com/google/uuid"

	"github.com/atrian/go-notify-customer/internal/dto"
)

type Service struct {
	storage Storager
}

func New() *Service {
	s := Service{
		storage: NewMemoryStorage(),
	}
	return &s
}

func (s Service) Start() {
	// log?
}

func (s Service) Stop() {
	//TODO implement me
	panic("implement me")
}

func (s Service) All(ctx context.Context) []dto.Template {
	templates, err := s.storage.All(ctx)

	if err != nil {
		// TODO log err
	}

	if templates == nil {
		return []dto.Template{}
	}

	return templates
}

func (s Service) Store(ctx context.Context, template dto.Template) (dto.Template, error) {
	template.TemplateUUID = uuid.New()

	err := s.storage.Store(ctx, template)
	if err != nil {
		return dto.Template{}, err
	}

	return template, nil
}

func (s Service) StoreBatch(ctx context.Context, templates []dto.Template) ([]dto.Template, error) {
	for i := 0; i < len(templates); i++ {
		templates[i].TemplateUUID = uuid.New()
		err := s.storage.Store(ctx, templates[i])
		if err != nil {
			// TODO err handling
			// Удалить все добавленные и вернуть ошибку?
		}
	}

	return templates, nil
}

func (s Service) Update(ctx context.Context, template dto.Template) (dto.Template, error) {
	err := s.storage.Store(ctx, template)
	if err != nil {
		return dto.Template{}, err
	}

	return template, nil
}

func (s Service) FindById(ctx context.Context, templateUUID uuid.UUID) (dto.Template, error) {
	return s.storage.GetById(ctx, templateUUID)
}

func (s Service) FindByEventId(ctx context.Context, eventUUID uuid.UUID) ([]dto.Template, error) {
	return s.storage.GetByEventId(ctx, eventUUID)
}

func (s Service) DeleteById(ctx context.Context, templateUUID uuid.UUID) error {
	return s.storage.DeleteById(ctx, templateUUID)
}
