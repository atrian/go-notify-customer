package notify

import (
	"github.com/atrian/go-notify-customer/internal/services/notify/entity"
)

// PriorityQueue приоритетная очередь, реализует heap.Interface
type PriorityQueue []*entity.Notification

// Len текущая длина очереди
func (pq PriorityQueue) Len() int { return len(pq) }

// Less сравнение элементов по приоритету
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority > pq[j].Priority
}

// Swap обмен мест элементов
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

// Push добавление уведомления
func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*entity.Notification)
	item.Index = n
	*pq = append(*pq, item)
}

// Pop получение "верхнего" - самого приоритетного элемента
func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}
