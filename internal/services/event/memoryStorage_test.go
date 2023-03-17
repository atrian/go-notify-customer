package event

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/atrian/go-notify-customer/internal/dto"

	"github.com/google/uuid"
)

func TestNewMemoryStorage(t *testing.T) {
	ms := NewMemoryStorage()
	ctx := context.TODO()

	event := dto.Event{
		EventUUID:            uuid.New(),
		Title:                "Test event",
		Description:          "Test description",
		DefaultPriority:      100,
		NotificationChannels: []string{"sms", "email"},
	}

	event2 := dto.Event{
		EventUUID:            uuid.New(),
		Title:                "Test event 2",
		Description:          "Test description 2",
		DefaultPriority:      1,
		NotificationChannels: []string{"email"},
	}

	// Сохраняем в хранилище события
	_ = ms.Store(ctx, event)
	_ = ms.Store(ctx, event2)

	// Запрос несуществующего события
	_, err := ms.GetById(ctx, uuid.New())
	if !errors.Is(NotFound, err) {
		t.Errorf("Expected \"%v\" error, got \"%v\"", NotFound, err)
	}

	// Получение всех событий
	allEvents, _ := ms.All(ctx)
	if len(allEvents) != 2 {
		t.Errorf("Event expected \"%v\", got \"%v\"", 2, len(allEvents))
	}

	// Получение события
	res1, err := ms.GetById(ctx, event.EventUUID)

	if errors.Is(NotFound, err) {
		t.Errorf("Expected no error, got \"%v\"", err)
	}
	if !reflect.DeepEqual(res1, event) {
		t.Errorf("Event expected \"%v\", got \"%v\"", event, res1)
	}

	// Удаление события
	err = ms.DeleteById(ctx, event.EventUUID)
	if errors.Is(NotFound, err) {
		t.Errorf("Expected no error, got \"%v\"", err)
	}

	// Запрос удаленного события
	_, err = ms.GetById(ctx, event.EventUUID)
	if !errors.Is(NotFound, err) {
		t.Errorf("Expected \"%v\" error, got \"%v\"", NotFound, err)
	}

	// При повторном удалении ожидаем ошибку
	err = ms.DeleteById(ctx, event.EventUUID)
	if !errors.Is(NotFound, err) {
		t.Errorf("Expected \"%v\" error, got \"%v\"", NotFound, err)
	}

	// Обновление события
	event2.Title = "Updated"
	err = ms.Update(ctx, event2)
	if errors.Is(NotFound, err) {
		t.Errorf("Expected no error, got \"%v\"", err)
	}

	// Получение обновленного события
	updatedEvent, err := ms.GetById(ctx, event2.EventUUID)
	if errors.Is(NotFound, err) {
		t.Errorf("Expected no error, got \"%v\"", err)
	}
	if !reflect.DeepEqual(updatedEvent, event2) {
		t.Errorf("Event expected \"%v\", got \"%v\"", event, res1)
	}
}
