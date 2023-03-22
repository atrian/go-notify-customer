package interfaces

import "context"

// BaseService интерфейс типового сервиса внутри приложения
type BaseService interface {
	// Start метод выполняется при запуске сервиса: регистрация в gateway, загрузка зависимостей
	Start(ctx context.Context)

	// Stop метод выполняется при завершении работы сервиса: закрытие ресурсов, отмена регистрации в gateway
	Stop()
}
