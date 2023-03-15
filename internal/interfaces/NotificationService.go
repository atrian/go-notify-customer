package interfaces

import "github.com/atrian/go-notify-customer/internal/services/notify/entity"

type NotificationService interface {
	// BaseService Общий сервисный интерфейс с методами Start и Stop
	BaseService

	// ProcessNotification приоритезация, лимитер уведомлений
	ProcessNotification(notification []entity.Notification) error
}
