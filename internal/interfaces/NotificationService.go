package interfaces

import (
	"context"

	"github.com/atrian/go-notify-customer/internal/dto"
)

type NotificationService interface {
	// BaseService Общий сервисный интерфейс с методами Start и Stop
	BaseService

	// ProcessNotification приоритезация, лимитер уведомлений
	ProcessNotification(ctx context.Context, notification []dto.Notification) error
}
