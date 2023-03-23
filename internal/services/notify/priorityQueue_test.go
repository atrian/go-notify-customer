package notify

import (
	"container/heap"
	"testing"

	"github.com/atrian/go-notify-customer/internal/dto"

	"github.com/google/uuid"
)

func TestPriorityQueue_Push(t *testing.T) {
	firstUUID := uuid.New()
	secondUUID := uuid.New()
	maxUUID := uuid.New()
	errUUID := uuid.New()

	notifications := []dto.Notification{
		{
			EventUUID: errUUID,
			Priority:  1500,
		}, {
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
		heap.Push(&pq, &notifications[i])
	}

	//  Bug testcase - use queue with heap!
	errNotification := heap.Pop(&pq).(*dto.Notification)
	if errNotification.EventUUID != errUUID {
		t.Errorf("got %q, wanted %q", errNotification.EventUUID, errUUID)
	}

	// Get max priority
	topNotification := heap.Pop(&pq).(*dto.Notification)
	if topNotification.EventUUID != maxUUID {
		t.Errorf("got %q, wanted %q", topNotification.EventUUID, maxUUID)
	}

	// Get mid-priority
	topNotification = heap.Pop(&pq).(*dto.Notification)
	if topNotification.EventUUID != secondUUID {
		t.Errorf("got %q, wanted %q", topNotification.EventUUID, secondUUID)
	}

	// Get low priority
	topNotification = heap.Pop(&pq).(*dto.Notification)
	if topNotification.EventUUID != firstUUID {
		t.Errorf("got %q, wanted %q", topNotification.EventUUID, firstUUID)
	}
}
