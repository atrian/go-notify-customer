package stat

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/atrian/go-notify-customer/internal/services/stat/entity"
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

func (m *MemoryStorage) All(ctx context.Context) ([]entity.Stat, error) {
	var stats []entity.Stat

	m.data.Range(func(key, value interface{}) bool {
		stat := value.(entity.Stat)
		stats = append(stats, stat)
		return true
	})

	return stats, nil
}

func (m *MemoryStorage) Store(ctx context.Context, stat entity.Stat) error {
	stat.CreatedAt = time.Now().Format(dateTimeFormat) // сохраняем время записи

	m.data.Store(stat.NotificationUUID, stat)

	return nil
}

func (m *MemoryStorage) GetByNotificationId(ctx context.Context, notificationUUID uuid.UUID) (entity.Stat, error) {
	var (
		stat  entity.Stat
		found bool
	)

	m.data.Range(func(key, value interface{}) bool {
		candidate := value.(entity.Stat)
		if candidate.NotificationUUID == notificationUUID {
			stat = candidate
			found = true
			return false
		}
		return true
	})

	if !found {
		return entity.Stat{}, NotFound
	}

	return stat, nil
}

func (m *MemoryStorage) GetByPersonId(ctx context.Context, personUUID uuid.UUID) ([]entity.Stat, error) {
	var (
		stats []entity.Stat
	)

	m.data.Range(func(key, value interface{}) bool {
		candidate := value.(entity.Stat)
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
