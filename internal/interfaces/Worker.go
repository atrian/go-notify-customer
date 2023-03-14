package interfaces

import "github.com/atrian/go-notify-customer/internal/services/notificationDispatcher/entity"

// Worker интерфейс типового воркера для отправки сообщений
type Worker interface {
	// Start метод запускающий воркера,
	// consumeQueue - очередь чтения,
	// successQueue - очередь записи положительного результата,
	// failQueue - очередь записи отризательного результата
	Start(consumeQueue string, successQueue string, failQueue string)

	// Send метод отправки сообщения через сервис провайдер на внешний сервис
	Send(message entity.Message)

	// Stop корректная остановка воркера
	Stop()
}
