package notify

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/atrian/go-notify-customer/internal/dto"
	"github.com/atrian/go-notify-customer/pkg/logger"
)

const bufferSize = 3

// TestService_ProcessNotification проверяем правильность приоритезации уведомлений
func TestService_ProcessNotification(t *testing.T) {
	log := logger.NewZapLogger()

	resultChan := make(chan dto.Notification, bufferSize)
	notifications := []dto.Notification{
		{
			EventUUID: uuid.New(),
			Priority:  1,
		}, {
			EventUUID: uuid.New(),
			Priority:  999,
		}, {
			EventUUID: uuid.New(),
			Priority:  500,
		},
	}

	s := New(resultChan, log)

	_ = s.ProcessNotification(context.TODO(), notifications)

	res1 := <-resultChan
	if res1.Priority != 999 {
		t.Errorf("got %v, wanted %v", res1.Priority, 999)
	}

	res2 := <-resultChan
	if res2.Priority != 500 {
		t.Errorf("got %v, wanted %v", res2.Priority, 500)
	}

	res3 := <-resultChan
	if res3.Priority != 1 {
		t.Errorf("got %v, wanted %v", res3.Priority, 1)
	}
}
