package stat

import (
	"context"

	"github.com/atrian/go-notify-customer/internal/services/stat/entity"
)

type Service struct {
}

func (s Service) Start() {
	//TODO implement me
	panic("implement me")
}

func (s Service) Stop() {
	//TODO implement me
	panic("implement me")
}

func (s Service) All(ctx context.Context) []entity.Stat {
	//TODO implement me
	panic("implement me")
}

func (s Service) Store(statChan chan entity.Stat) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) FindByPersonUUID(ctx context.Context, orderId string) ([]entity.Stat, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) FindByEventUUID(ctx context.Context, orderId string) ([]entity.Stat, error) {
	//TODO implement me
	panic("implement me")
}
