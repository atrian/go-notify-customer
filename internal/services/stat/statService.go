package stat

import (
	"context"

	"github.com/google/uuid"

	"github.com/atrian/go-notify-customer/internal/services/stat/entity"
)

type Service struct {
	statChan chan entity.Stat
	ctx      context.Context
	storage  Storager
}

func New(ctx context.Context, statChan chan entity.Stat) *Service {
	s := Service{
		statChan: statChan,
		storage:  NewMemoryStorage(),
		ctx:      ctx,
	}

	return &s
}

func (s Service) Start() {
	go func(ctx context.Context, statChan chan entity.Stat) {
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

func (s Service) All(ctx context.Context) []entity.Stat {
	res, err := s.storage.All(ctx)
	if err != nil {
		// TODO handle err
	}

	return res
}

func (s Service) Store(stat entity.Stat) error {
	return s.storage.Store(s.ctx, stat)
}

func (s Service) FindByPersonUUID(ctx context.Context, personUUID uuid.UUID) ([]entity.Stat, error) {
	return s.storage.GetByPersonId(ctx, personUUID)
}

func (s Service) FindByEventUUID(ctx context.Context, personUUID uuid.UUID) (entity.Stat, error) {
	return s.storage.GetByNotificationId(ctx, personUUID)
}
