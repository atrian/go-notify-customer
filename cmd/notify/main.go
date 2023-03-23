package main

import (
	"context"

	"github.com/atrian/go-notify-customer/internal/notify"
)

// @Title Go-notify-client
// @Description Сервис отправки уведомлений.
// @Version 1.0

// @Host localhost:8080
// @BasePath /

// @Tag.name Event
// @Tag.description "Группа запросов бизнес событий"

// @Tag.name Notifications
// @Tag.description "Группа запросов уведомлений"

// @Tag.name Stat
// @Tag.description "Группа запросов для работы со статистикой отправки"

// @Tag.name Template
// @Tag.description "Группа запросов для работы с шаблонами сообщений. Для создания сообщений требуется предварительное создание бизнес событий"

func main() {
	ctx := context.Background()

	application := notify.New()
	application.Run(ctx)
}
