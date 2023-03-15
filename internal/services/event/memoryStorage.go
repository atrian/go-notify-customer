package event

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"

	"github.com/atrian/go-notify-customer/internal/services/event/entity"
)

var NotFound = errors.New("not found")

type MemoryStorage struct {
	data sync.Map
}

func NewMemoryStorage() *MemoryStorage {
	ms := MemoryStorage{}
	return &ms
}

func (m *MemoryStorage) All(ctx context.Context) ([]entity.Event, error) {
	var events []entity.Event

	m.data.Range(func(key, value interface{}) bool {
		event := value.(entity.Event)
		events = append(events, event)
		return true
	})

	return events, nil
}

func (m *MemoryStorage) Update(ctx context.Context, event entity.Event) error {
	m.data.Store(event.EventUUID.String(), event)

	return nil
}

func (m *MemoryStorage) Store(ctx context.Context, event entity.Event) error {
	m.data.Store(event.EventUUID.String(), event)

	return nil
}

func (m *MemoryStorage) GetById(ctx context.Context, eventUUID uuid.UUID) (entity.Event, error) {
	event, ok := m.data.Load(eventUUID.String())

	if !ok {
		return entity.Event{}, NotFound
	}

	return event.(entity.Event), nil
}

func (m *MemoryStorage) DeleteById(ctx context.Context, eventUUID uuid.UUID) error {
	_, ok := m.data.Load(eventUUID.String())

	if !ok {
		return NotFound
	}

	m.data.Delete(eventUUID.String())

	return nil
}
