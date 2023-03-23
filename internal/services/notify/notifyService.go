// Package notify фронт сервис для приема уведомлений на отправку
// Выполняет приоритезацию уведомлений и передает далее в NotificationDispatcher
package notify

import (
	"container/heap"
	"context"
	"errors"

	"github.com/atrian/go-notify-customer/internal/dto"
	"github.com/atrian/go-notify-customer/internal/interfaces"
)

var (
	NotificationLimitExceeded = errors.New("notification limit exceed")
)

type Service struct {
	queue      PriorityQueue
	resultChan chan<- dto.Notification
	logger     interfaces.Logger
}

// New Конфигурация зависимостей сервиса
func New(resultChan chan dto.Notification, logger interfaces.Logger) *Service {
	s := Service{
		queue:      nil,        // очередь с приоритетом
		resultChan: resultChan, // выходной канал после приоритезации сообщений
		logger:     logger,
	}

	heap.Init(&s.queue)

	return &s
}

// Start стартовые процедуры - логгер?
func (s Service) Start(ctx context.Context) {
	s.logger.Info("Notification service started")
}

func (s Service) Stop() {
	close(s.resultChan)
	s.logger.Info("Notification service stopped")
}

func (s Service) ProcessNotification(ctx context.Context, notifications []dto.Notification) error {
	// приоритизация очереди уведомлений
	for i := 0; i < len(notifications); i++ {
		heap.Push(&s.queue, &notifications[i])
	}

	// обрабатываем очередь в порядке приоритета и отдаем в результирующий канал
	// TODO добавить RATE LIMITER на количество отправленных
	for s.queue.Len() > 0 {
		item := heap.Pop(&s.queue).(*dto.Notification)
		s.resultChan <- *item
	}

	return nil
}
