package event

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"

	"github.com/atrian/go-notify-customer/internal/dto"
)

var NotFound = errors.New("not found")

// MemoryStorage in-memory хранилище для сервиса event
// ! потокобезопасно, работает на sync.Map
// ! is safe for concurrent use
type MemoryStorage struct {
	data sync.Map
}

func NewMemoryStorage() *MemoryStorage {
	ms := MemoryStorage{}
	return &ms
}

func (m *MemoryStorage) All(ctx context.Context) ([]dto.Event, error) {
	var events []dto.Event

	m.data.Range(func(key, value interface{}) bool {
		event := value.(dto.Event)
		events = append(events, event)
		return true
	})

	return events, nil
}

func (m *MemoryStorage) Update(ctx context.Context, event dto.Event) error {
	m.data.Store(event.EventUUID.String(), event)

	return nil
}

func (m *MemoryStorage) Store(ctx context.Context, event dto.Event) error {
	m.data.Store(event.EventUUID.String(), event)

	return nil
}

func (m *MemoryStorage) GetById(ctx context.Context, eventUUID uuid.UUID) (dto.Event, error) {
	event, ok := m.data.Load(eventUUID.String())

	if !ok {
		return dto.Event{}, NotFound
	}

	return event.(dto.Event), nil
}

func (m *MemoryStorage) DeleteById(ctx context.Context, eventUUID uuid.UUID) error {
	_, ok := m.data.Load(eventUUID.String())

	if !ok {
		return NotFound
	}

	m.data.Delete(eventUUID.String())

	return nil
}
