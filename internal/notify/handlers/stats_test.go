package handlers_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sort"

	"github.com/google/uuid"

	"github.com/atrian/go-notify-customer/internal/dto"
	"github.com/atrian/go-notify-customer/internal/notify/handlers"
	"github.com/atrian/go-notify-customer/internal/notify/router"
	"github.com/atrian/go-notify-customer/internal/services/stat"
	"github.com/atrian/go-notify-customer/pkg/logger"
)

func ExampleHandler_GetStats() {
	testStat := dto.Stat{
		PersonUUID:       uuid.New(),
		NotificationUUID: uuid.New(),
		CreatedAt:        "",
		Status:           dto.Sent,
	}

	testStat2 := dto.Stat{
		PersonUUID:       uuid.New(),
		NotificationUUID: uuid.New(),
		CreatedAt:        "",
		Status:           dto.Failed,
	}

	// Подготавливаем все зависимости, логгер, конфигурацию приложения, хранилище (In Memory) и роутер
	appLogger := logger.NewZapLogger()
	appConf := mockHandlerConfig{}

	service := stat.New(nil, appLogger)
	_ = service.Store(context.TODO(), testStat)
	_ = service.Store(context.TODO(), testStat2)

	getEndpoint := fmt.Sprintf("/api/v1/stats")

	h := handlers.New(&appConf, nil, nil, service, nil, appLogger)

	r := router.New(h, &appConf)

	// Запускаем тестовый сервер
	testServer := httptest.NewServer(r)
	defer testServer.Close()

	// Удаление данных
	request, _ := http.NewRequest(http.MethodGet, testServer.URL+getEndpoint, nil)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		appLogger.Error("http.DefaultClient.Do err", err)
	}

	responseBody, _ := io.ReadAll(response.Body)
	_ = response.Body.Close()

	// В случае успеха сервис отвечает кодом 200
	appLogger.Debug(fmt.Sprintf("GET OK - status: %v", response.StatusCode))

	// Манипуляции для отсечения изменяющегося при сохранении UUID
	var getResult []dto.Stat
	_ = json.Unmarshal(responseBody, &getResult)

	stableSlice := make([]string, 0, 2)
	for _, t := range getResult {
		stableSlice = append(stableSlice, t.Status.String())
	}

	sort.Strings(stableSlice)

	fmt.Println(response.StatusCode, stableSlice)

	// Output:
	// 200 [failed sent]
}

func ExampleHandler_GetStatByPersonUUID() {
	personUUID := uuid.New()
	testStat := dto.Stat{
		PersonUUID:       personUUID,
		NotificationUUID: uuid.New(),
		CreatedAt:        "",
		Status:           dto.Sent,
	}

	testStat2 := dto.Stat{
		PersonUUID:       personUUID,
		NotificationUUID: uuid.New(),
		CreatedAt:        "",
		Status:           dto.Failed,
	}

	testStat3 := dto.Stat{
		PersonUUID:       uuid.New(),
		NotificationUUID: uuid.New(),
		CreatedAt:        "",
		Status:           dto.Failed,
	}

	// Подготавливаем все зависимости, логгер, конфигурацию приложения, хранилище (In Memory) и роутер
	appLogger := logger.NewZapLogger()
	appConf := mockHandlerConfig{}

	service := stat.New(nil, nil)
	_ = service.Store(context.TODO(), testStat)
	_ = service.Store(context.TODO(), testStat2)
	_ = service.Store(context.TODO(), testStat3)

	getEndpoint := fmt.Sprintf("/api/v1/stats/person/%v", personUUID)

	h := handlers.New(&appConf, nil, nil, service, nil, appLogger)

	r := router.New(h, &appConf)

	// Запускаем тестовый сервер
	testServer := httptest.NewServer(r)
	defer testServer.Close()

	// Удаление данных
	request, _ := http.NewRequest(http.MethodGet, testServer.URL+getEndpoint, nil)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		appLogger.Error("http.DefaultClient.Do err", err)
	}

	responseBody, _ := io.ReadAll(response.Body)
	_ = response.Body.Close()

	// В случае успеха сервис отвечает кодом 200
	appLogger.Debug(fmt.Sprintf("GET OK - status: %v", response.StatusCode))

	// Манипуляции для отсечения изменяющегося при сохранении UUID
	var getResult []dto.Stat
	_ = json.Unmarshal(responseBody, &getResult)

	stableSlice := make([]string, 0, 2)
	for _, t := range getResult {
		stableSlice = append(stableSlice, t.Status.String())
	}

	sort.Strings(stableSlice)

	fmt.Println(response.StatusCode, stableSlice)

	// Output:
	// 200 [failed sent]
}

func ExampleHandler_GetStatByNotificationId() {
	notificationUUID := uuid.New()
	testStat := dto.Stat{
		PersonUUID:       uuid.New(),
		NotificationUUID: notificationUUID,
		CreatedAt:        "",
		Status:           dto.Failed,
	}

	testStat2 := dto.Stat{
		PersonUUID:       uuid.New(),
		NotificationUUID: notificationUUID,
		CreatedAt:        "",
		Status:           dto.Failed,
	}

	testStat3 := dto.Stat{
		PersonUUID:       uuid.New(),
		NotificationUUID: uuid.New(),
		CreatedAt:        "",
		Status:           dto.Failed,
	}

	// Подготавливаем все зависимости, логгер, конфигурацию приложения, хранилище (In Memory) и роутер
	appLogger := logger.NewZapLogger()
	appConf := mockHandlerConfig{}

	service := stat.New(nil, nil)
	_ = service.Store(context.TODO(), testStat)
	_ = service.Store(context.TODO(), testStat2)
	_ = service.Store(context.TODO(), testStat3)

	getEndpoint := fmt.Sprintf("/api/v1/stats/notification/%v", notificationUUID)

	h := handlers.New(&appConf, nil, nil, service, nil, appLogger)

	r := router.New(h, &appConf)

	// Запускаем тестовый сервер
	testServer := httptest.NewServer(r)
	defer testServer.Close()

	// Удаление данных
	request, _ := http.NewRequest(http.MethodGet, testServer.URL+getEndpoint, nil)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		appLogger.Error("http.DefaultClient.Do err", err)
	}

	responseBody, _ := io.ReadAll(response.Body)
	_ = response.Body.Close()

	// В случае успеха сервис отвечает кодом 200
	appLogger.Debug(fmt.Sprintf("GET OK - status: %v", response.StatusCode))

	// Манипуляции для отсечения изменяющегося при сохранении UUID
	var getResult []dto.Stat
	_ = json.Unmarshal(responseBody, &getResult)

	stableSlice := make([]string, 0, 2)
	for _, t := range getResult {
		stableSlice = append(stableSlice, t.Status.String())
	}

	sort.Strings(stableSlice)

	fmt.Println(response.StatusCode, stableSlice)

	// Output:
	// 200 [failed failed]
}
