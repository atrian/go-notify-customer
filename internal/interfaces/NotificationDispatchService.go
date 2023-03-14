package interfaces

import "github.com/atrian/go-notify-customer/internal/services/notificationDispatcher/entity"

// NotificationDispatchService интерфейс сервиса отправки уведомлений в шину
type NotificationDispatchService interface {
	// BaseService Общий сервисный интерфейс с методами Start и Stop
	BaseService

	Dispatch(message entity.Message) error
}
