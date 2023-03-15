package notify

import (
	"container/heap"
	"errors"

	"github.com/atrian/go-notify-customer/internal/services/notify/entity"
)

var (
	NotificationLimitExceeded = errors.New("notification limit exceed")
)

type Service struct {
	//limiter           RateLimiter
	//notificationLimit int
	queue      PriorityQueue
	resultChan chan entity.Notification
}

// New Конфигурация зависимостей сервиса
func New(resultChan chan entity.Notification) *Service {
	s := Service{
		queue:      nil,
		resultChan: resultChan,
	}

	heap.Init(&s.queue)

	return &s
}

// Start стартовые процедуры - логгер?
func (s Service) Start() {

}

func (s Service) Stop() {
	close(s.resultChan)
}

func (s Service) ProcessNotification(notifications []entity.Notification) error {
	// приоритизация очереди уведомлений
	for i := 0; i < len(notifications); i++ {
		heap.Push(&s.queue, &notifications[i])
	}

	// обрабатываем очередь в порядке приоритета и отдаем в результирующий канал
	// TODO добавить RATE LIMITER на количество отправленных
	for s.queue.Len() > 0 {
		item := heap.Pop(&s.queue).(*entity.Notification)
		s.resultChan <- *item
	}

	return nil
}
