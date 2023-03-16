package template

import (
	"context"

	"github.com/google/uuid"

	"github.com/atrian/go-notify-customer/internal/services/template/entity"
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

func (s Service) All(ctx context.Context) []entity.Template {
	templates, err := s.storage.All(ctx)

	if err != nil {
		// TODO log err
	}

	return templates
}

func (s Service) Store(ctx context.Context, template entity.Template) (entity.Template, error) {
	template.EventUUID = uuid.New()

	err := s.storage.Store(ctx, template)
	if err != nil {
		return entity.Template{}, err
	}

	return template, nil
}

func (s Service) StoreBatch(ctx context.Context, templates []entity.Template) ([]entity.Template, error) {
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

func (s Service) Update(ctx context.Context, template entity.Template) (entity.Template, error) {
	err := s.storage.Store(ctx, template)
	if err != nil {
		return entity.Template{}, err
	}

	return template, nil
}

func (s Service) FindById(ctx context.Context, templateUUID uuid.UUID) (entity.Template, error) {
	return s.storage.GetById(ctx, templateUUID)
}

func (s Service) FindByEventId(ctx context.Context, eventUUID uuid.UUID) (entity.Template, error) {
	return s.storage.GetByEventId(ctx, eventUUID)
}

func (s Service) DeleteById(ctx context.Context, templateUUID uuid.UUID) error {
	return s.storage.DeleteById(ctx, templateUUID)
}
