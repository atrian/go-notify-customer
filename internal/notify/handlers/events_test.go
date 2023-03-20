package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"

	"github.com/google/uuid"

	"github.com/atrian/go-notify-customer/internal/dto"
	"github.com/atrian/go-notify-customer/internal/notify/handlers"
	"github.com/atrian/go-notify-customer/internal/notify/router"
	"github.com/atrian/go-notify-customer/internal/services/event"
	"github.com/atrian/go-notify-customer/pkg/logger"
)

func ExampleHandler_UpdateEvent() {
	testEvent := dto.Event{
		Title:                "Test",
		Description:          "Description",
		DefaultPriority:      1,
		NotificationChannels: []string{"sms"},
	}

	// Подготавливаем все зависимости, логгер, конфигурацию приложения, хранилище (In Memory) и роутер
	appLogger := logger.NewZapLogger()
	appConf := mockHandlerConfig{}

	eService := event.New()
	stored, _ := eService.Store(context.Background(), testEvent)

	h := handlers.New(&appConf, eService, nil, nil, nil, appLogger)

	r := router.New(h, &appConf)

	// Запускаем тестовый сервер
	testServer := httptest.NewServer(r)
	defer testServer.Close()

	// Обновление данных
	testEvent.Title = "Updated title"
	testEvent.EventUUID = stored.EventUUID
	jData, _ := json.Marshal(testEvent)
	jReader := bytes.NewReader(jData)

	updateEndpoint := fmt.Sprintf("/api/v1/events/%v", stored.EventUUID)

	request, _ := http.NewRequest(http.MethodPut, testServer.URL+updateEndpoint, jReader)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		appLogger.Error("http.DefaultClient.Do err", err)
	}

	responseBody, _ := io.ReadAll(response.Body)
	_ = response.Body.Close()

	// Манипуляции для отсечения изменяющегося при сохранении UUID
	var getResult dto.Event
	_ = json.Unmarshal(responseBody, &getResult)

	// В случае успеха сервис отвечает кодом 200 и JSON содержащим текущее значение шаблона
	fmt.Println(response.StatusCode, getResult.Title)

	// Output:
	// 200 Updated title
}

func ExampleHandler_DeleteEvent() {
	testEvent := dto.Event{
		Title:                "Test",
		Description:          "Description",
		DefaultPriority:      1,
		NotificationChannels: []string{"sms"},
	}

	// Подготавливаем все зависимости, логгер, конфигурацию приложения, хранилище (In Memory) и роутер
	appLogger := logger.NewZapLogger()
	appConf := mockHandlerConfig{}

	eService := event.New()
	stored, _ := eService.Store(context.Background(), testEvent)

	h := handlers.New(&appConf, eService, nil, nil, nil, appLogger)

	r := router.New(h, &appConf)

	// Запускаем тестовый сервер
	testServer := httptest.NewServer(r)
	defer testServer.Close()

	deleteEndpoint := fmt.Sprintf("/api/v1/events/%v", stored.EventUUID)

	// Удаление данных
	request, _ := http.NewRequest(http.MethodDelete, testServer.URL+deleteEndpoint, nil)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		appLogger.Error("http.DefaultClient.Do err", err)
	}

	responseBody, _ := io.ReadAll(response.Body)
	_ = response.Body.Close()

	// В случае успеха сервис отвечает кодом 200
	appLogger.Debug(fmt.Sprintf("DELETED OK - status: %v", response.StatusCode))

	response, err = http.DefaultClient.Do(request)
	if err != nil {
		appLogger.Error("http.DefaultClient.Do err", err)
	}

	// проверяем что запись удалена и повторный вызов вернет 404 ошибку
	responseBody, _ = io.ReadAll(response.Body)
	_ = response.Body.Close()

	appLogger.Debug(fmt.Sprintf("SECOND ATTEMPT FAILED - status: %v", response.StatusCode))
	fmt.Println(response.StatusCode, string(responseBody))

	// Output:
	// 404 Not found
}

func ExampleHandler_GetEvent() {
	testEvent := dto.Event{
		Title:                "Title",
		Description:          "Description",
		DefaultPriority:      0,
		NotificationChannels: nil,
	}

	// Подготавливаем все зависимости, логгер, конфигурацию приложения, хранилище (In Memory) и роутер
	appLogger := logger.NewZapLogger()
	appConf := mockHandlerConfig{}

	eService := event.New()
	stored, _ := eService.Store(context.Background(), testEvent)

	getEndpoint := fmt.Sprintf("/api/v1/events/%v", stored.EventUUID)

	h := handlers.New(&appConf, eService, nil, nil, nil, appLogger)

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
	var getResult dto.Event
	_ = json.Unmarshal(responseBody, &getResult)
	getResult.EventUUID = uuid.UUID{}
	fmt.Println(response.StatusCode, getResult)

	// Output:
	// 200 {00000000-0000-0000-0000-000000000000 Title Description 0 []}
}

func ExampleHandler_GetEvents() {
	testEvent := dto.Event{
		Title:                "Title 1",
		Description:          "Description 1",
		DefaultPriority:      0,
		NotificationChannels: nil,
	}

	testEvent2 := dto.Event{
		Title:                "Title 2",
		Description:          "Description 2",
		DefaultPriority:      0,
		NotificationChannels: nil,
	}

	// Подготавливаем все зависимости, логгер, конфигурацию приложения, хранилище (In Memory) и роутер
	appLogger := logger.NewZapLogger()
	appConf := mockHandlerConfig{}

	service := event.New()
	_, _ = service.Store(context.Background(), testEvent)
	_, _ = service.Store(context.Background(), testEvent2)

	getEndpoint := fmt.Sprintf("/api/v1/events")

	h := handlers.New(&appConf, service, nil, nil, nil, appLogger)

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
	var getResult []dto.Event
	_ = json.Unmarshal(responseBody, &getResult)

	stableSlice := make([]string, 0, 2)
	for _, t := range getResult {
		checkThis := strings.Builder{}
		checkThis.WriteString(t.Title)
		checkThis.WriteString(t.Description)
		stableSlice = append(stableSlice, checkThis.String())
	}

	sort.Strings(stableSlice)

	fmt.Println(response.StatusCode, stableSlice)

	// Output:
	// 200 [Title 1Description 1 Title 2Description 2]
}

func ExampleHandler_StoreEvent() {
	storeEndpoint := fmt.Sprintf("/api/v1/events")

	testTemplate := dto.IncomingEvent{
		Title:                "EventTitle",
		Description:          "Description",
		DefaultPriority:      1,
		NotificationChannels: []string{"sms"},
	}

	// Подготавливаем все зависимости, логгер, конфигурацию приложения, хранилище (In Memory) и роутер
	appLogger := logger.NewZapLogger()
	appConf := mockHandlerConfig{}

	service := event.New()

	h := handlers.New(&appConf, service, nil, nil, nil, appLogger)

	r := router.New(h, &appConf)

	// Запускаем тестовый сервер
	testServer := httptest.NewServer(r)
	defer testServer.Close()

	// Сохранение данных
	jData, _ := json.Marshal(testTemplate)
	jReader := bytes.NewReader(jData)

	request, _ := http.NewRequest(http.MethodPost, testServer.URL+storeEndpoint, jReader)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		appLogger.Error("http.DefaultClient.Do err", err)
	}

	responseBody, _ := io.ReadAll(response.Body)
	_ = response.Body.Close()

	// Манипуляции для отсечения изменяющегося при сохранении UUID
	var getResult dto.Event
	_ = json.Unmarshal(responseBody, &getResult)

	checkThis := strings.Builder{}
	checkThis.WriteString(getResult.Title)
	checkThis.WriteString(getResult.Description)

	// В случае успеха сервис отвечает кодом 200 и JSON содержащим текущее значение шаблона
	fmt.Println(response.StatusCode, checkThis.String())

	// Output:
	// 200 EventTitleDescription
}
