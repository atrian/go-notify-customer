package template

import (
	"context"
	"errors"
	"sync"

	"github.com/atrian/go-notify-customer/internal/dto"

	"github.com/google/uuid"
)

var NotFound = errors.New("not found")

type MemoryStorage struct {
	data sync.Map
}

func NewMemoryStorage() *MemoryStorage {
	ms := MemoryStorage{}
	return &ms
}

func (m *MemoryStorage) All(ctx context.Context) ([]dto.Template, error) {
	var templates []dto.Template

	m.data.Range(func(key, value interface{}) bool {
		template := value.(dto.Template)
		templates = append(templates, template)
		return true
	})

	return templates, nil
}

func (m *MemoryStorage) Update(ctx context.Context, template dto.Template) error {
	m.data.Store(template.TemplateUUID.String(), template)

	return nil
}

func (m *MemoryStorage) Store(ctx context.Context, template dto.Template) error {
	m.data.Store(template.TemplateUUID.String(), template)

	return nil
}

func (m *MemoryStorage) GetById(ctx context.Context, templateUUID uuid.UUID) (dto.Template, error) {
	template, ok := m.data.Load(templateUUID.String())

	if !ok {
		return dto.Template{}, NotFound
	}

	return template.(dto.Template), nil
}

func (m *MemoryStorage) GetByEventId(ctx context.Context, eventUUID uuid.UUID) (dto.Template, error) {
	var (
		template dto.Template
		exist    bool
	)

	m.data.Range(func(key, value interface{}) bool {
		candidate := value.(dto.Template)
		if candidate.EventUUID == eventUUID {
			template = candidate
			exist = true
			return false
		}
		return true
	})

	if !exist {
		return dto.Template{}, NotFound
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
