package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/atrian/go-notify-customer/internal/dto"
	"github.com/atrian/go-notify-customer/internal/notify/handlers"
	"github.com/atrian/go-notify-customer/internal/notify/router"
	"github.com/atrian/go-notify-customer/internal/services/notify"
	"github.com/atrian/go-notify-customer/pkg/logger"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
)

func ExampleHandler_ProcessNotifications() {
	notifications := []dto.Notification{
		{
			NotificationUUID: uuid.New(),
			EventUUID:        uuid.New(),
			PersonUUIDs: []uuid.UUID{
				uuid.New(),
			},
			MessageParams: nil,
			Priority:      100,
		}, {
			NotificationUUID: uuid.New(),
			EventUUID:        uuid.New(),
			PersonUUIDs: []uuid.UUID{
				uuid.New(),
			},
			MessageParams: nil,
			Priority:      1500,
		},
	}

	// Подготавливаем все зависимости, логгер, конфигурацию приложения, хранилище (In Memory) и роутер
	appLogger := logger.NewZapLogger()
	appConf := mockHandlerConfig{}

	resultChan := make(chan dto.Notification)
	done := make(chan struct{})

	go func() {
		first := <-resultChan
		appLogger.Info(fmt.Sprintf("first notification priority: %v", first.Priority))
		second := <-resultChan
		appLogger.Info(fmt.Sprintf("second notification priority: %v", second.Priority))

		close(done)
	}()

	service := notify.New(resultChan, appLogger)
	getEndpoint := fmt.Sprintf("/api/v1/notifications")

	h := handlers.New(&appConf, nil, service, nil, nil, appLogger)
	r := router.New(h, &appConf)

	// Запускаем тестовый сервер
	testServer := httptest.NewServer(r)
	defer testServer.Close()

	// Отправка данных
	jData, _ := json.Marshal(notifications)
	jReader := bytes.NewReader(jData)
	request, _ := http.NewRequest(http.MethodPost, testServer.URL+getEndpoint, jReader)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		appLogger.Error("http.DefaultClient.Do err", err)
	}

	// В случае успеха сервис отвечает кодом 200
	appLogger.Debug(fmt.Sprintf("GET OK - status: %v", response.StatusCode))

	// Манипуляции для отсечения изменяющегося при сохранении UUID
	<-done
	fmt.Println(response.StatusCode)

	// Output:
	// 200
}
