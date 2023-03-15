package notify

import (
	"testing"

	"github.com/google/uuid"

	"github.com/atrian/go-notify-customer/internal/services/notify/entity"
)

const bufferSize = 3

func TestService_ProcessNotification(t *testing.T) {
	resultChan := make(chan entity.Notification, bufferSize)
	notifications := []entity.Notification{
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

	s := New(resultChan)

	_ = s.ProcessNotification(notifications)

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
