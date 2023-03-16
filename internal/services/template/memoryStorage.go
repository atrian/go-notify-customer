package template

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"

	"github.com/atrian/go-notify-customer/internal/services/template/entity"
)

var NotFound = errors.New("not found")

type MemoryStorage struct {
	data sync.Map
}

func NewMemoryStorage() *MemoryStorage {
	ms := MemoryStorage{}
	return &ms
}

func (m *MemoryStorage) All(ctx context.Context) ([]entity.Template, error) {
	var templates []entity.Template

	m.data.Range(func(key, value interface{}) bool {
		template := value.(entity.Template)
		templates = append(templates, template)
		return true
	})

	return templates, nil
}

func (m *MemoryStorage) Update(ctx context.Context, template entity.Template) error {
	m.data.Store(template.TemplateUUID.String(), template)

	return nil
}

func (m *MemoryStorage) Store(ctx context.Context, template entity.Template) error {
	m.data.Store(template.TemplateUUID.String(), template)

	return nil
}

func (m *MemoryStorage) GetById(ctx context.Context, templateUUID uuid.UUID) (entity.Template, error) {
	template, ok := m.data.Load(templateUUID.String())

	if !ok {
		return entity.Template{}, NotFound
	}

	return template.(entity.Template), nil
}

func (m *MemoryStorage) GetByEventId(ctx context.Context, eventUUID uuid.UUID) (entity.Template, error) {
	var (
		template entity.Template
		exist    bool
	)

	m.data.Range(func(key, value interface{}) bool {
		candidate := value.(entity.Template)
		if candidate.EventUUID == eventUUID {
			template = candidate
			exist = true
			return false
		}
		return true
	})

	if !exist {
		return entity.Template{}, NotFound
	}

	return template, nil
}

func (m *MemoryStorage) DeleteById(ctx context.Context, templateUUID uuid.UUID) error {
	_, ok := m.data.Load(templateUUID.String())

	if !ok {
		return NotFound
	}

	m.data.Delete(templateUUID.String())

	return nil
}
