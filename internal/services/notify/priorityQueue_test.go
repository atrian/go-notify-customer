package notify

import (
	"container/heap"
	"testing"

	"github.com/google/uuid"

	"github.com/atrian/go-notify-customer/internal/services/notify/entity"
)

func TestPriorityQueue_Push(t *testing.T) {
	firstUUID := uuid.New()
	secondUUID := uuid.New()
	maxUUID := uuid.New()

	notifications := []entity.Notification{
		{
			EventUUID: firstUUID,
			Priority:  1,
		}, {
			EventUUID: secondUUID,
			Priority:  2,
		}, {
			EventUUID: maxUUID,
			Priority:  999,
		},
	}

	var pq PriorityQueue
	heap.Init(&pq)

	for i := 0; i < len(notifications); i++ {
		pq.Push(&notifications[i])
	}

	// Get max priority
	topNotification := pq.Pop().(*entity.Notification)
	if topNotification.EventUUID != maxUUID {
		t.Errorf("got %q, wanted %q", topNotification.EventUUID, maxUUID)
	}

	// Get mid-priority
	topNotification = pq.Pop().(*entity.Notification)
	if topNotification.EventUUID != secondUUID {
		t.Errorf("got %q, wanted %q", topNotification.EventUUID, secondUUID)
	}

	// Get low priority
	topNotification = pq.Pop().(*entity.Notification)
	if topNotification.EventUUID != firstUUID {
		t.Errorf("got %q, wanted %q", topNotification.EventUUID, firstUUID)
	}
}
