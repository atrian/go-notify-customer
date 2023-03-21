package stat

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/atrian/go-notify-customer/internal/dto"
)

type Service struct {
	statChan chan dto.Stat
	ctx      context.Context
	storage  Storager
}

func New(ctx context.Context, statChan chan dto.Stat) *Service {
	s := Service{
		statChan: statChan,
		storage:  NewMemoryStorage(),
		ctx:      ctx,
	}

	return &s
}

func (s Service) Start() {
	go func(ctx context.Context, statChan chan dto.Stat) {
		for {
			select {
			case stat := <-statChan:
				_ = s.Store(stat)

			case <-ctx.Done():
				// TODO shutdown
				return

			default:
				// do nothing
			}
		}
	}(s.ctx, s.statChan)
}

func (s Service) Stop() {
	//TODO implement me
	panic("implement me")
}

func (s Service) All(ctx context.Context) []dto.Stat {
	res, err := s.storage.All(ctx)
	if err != nil {
		// TODO handle err
	}

	if res == nil {
		return []dto.Stat{}
	}

	return res
}

func (s Service) Store(stat dto.Stat) error {
	stat.StatUUID = uuid.New()
	stat.CreatedAt = time.Now().Format(dateTimeFormat) // сохраняем время записи

	return s.storage.Store(s.ctx, stat)
}

func (s Service) FindByPersonUUID(ctx context.Context, personUUID uuid.UUID) ([]dto.Stat, error) {
	return s.storage.GetByPersonId(ctx, personUUID)
}

func (s Service) FindByNotificationId(ctx context.Context, notificationUUID uuid.UUID) ([]dto.Stat, error) {
	return s.storage.GetByNotificationId(ctx, notificationUUID)
}
