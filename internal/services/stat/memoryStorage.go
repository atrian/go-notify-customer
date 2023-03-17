package stat

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/atrian/go-notify-customer/internal/dto"

	"github.com/google/uuid"
)

const dateTimeFormat = "2006-01-02 15:04:05"

var NotFound = errors.New("not found")

type MemoryStorage struct {
	data sync.Map
}

func NewMemoryStorage() *MemoryStorage {
	ms := MemoryStorage{}
	return &ms
}

func (m *MemoryStorage) All(ctx context.Context) ([]dto.Stat, error) {
	var stats []dto.Stat

	m.data.Range(func(key, value interface{}) bool {
		stat := value.(dto.Stat)
		stats = append(stats, stat)
		return true
	})

	return stats, nil
}

func (m *MemoryStorage) Store(ctx context.Context, stat dto.Stat) error {
	stat.CreatedAt = time.Now().Format(dateTimeFormat) // сохраняем время записи

	m.data.Store(stat.NotificationUUID, stat)

	return nil
}

func (m *MemoryStorage) GetByNotificationId(ctx context.Context, notificationUUID uuid.UUID) (dto.Stat, error) {
	var (
		stat  dto.Stat
		found bool
	)

	m.data.Range(func(key, value interface{}) bool {
		candidate := value.(dto.Stat)
		if candidate.NotificationUUID == notificationUUID {
			stat = candidate
			found = true
			return false
		}
		return true
	})

	if !found {
		return dto.Stat{}, NotFound
	}

	return stat, nil
}

func (m *MemoryStorage) GetByPersonId(ctx context.Context, personUUID uuid.UUID) ([]dto.Stat, error) {
	var (
		stats []dto.Stat
	)

	m.data.Range(func(key, value interface{}) bool {
		candidate := value.(dto.Stat)
		if candidate.PersonUUID == personUUID {
			stats = append(stats, candidate)
		}
		return true
	})

	if len(stats) == 0 {
		return nil, NotFound
	}

	return stats, nil
}
