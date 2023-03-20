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
	"github.com/atrian/go-notify-customer/internal/services/template"
	"github.com/atrian/go-notify-customer/pkg/logger"
)

func ExampleHandler_UpdateTemplate() {
	templateUUID, _ := uuid.Parse("c10a7fa8-f162-46ce-a97e-ff8718d7eb7d")

	updateEndpoint := fmt.Sprintf("/api/v1/templates/%v", templateUUID)

	testTemplate := dto.Template{
		TemplateUUID: templateUUID,
		EventUUID:    uuid.UUID{},
		Title:        "Test",
		Description:  "Description",
		Body:         "Body",
		ChannelType:  "ChannelType",
	}

	// Подготавливаем все зависимости, логгер, конфигурацию приложения, хранилище (In Memory) и роутер
	appLogger := logger.NewZapLogger()
	appConf := mockHandlerConfig{}

	tService := template.New()
	_, _ = tService.Store(context.Background(), testTemplate)

	h := handlers.New(&appConf, nil, nil, nil, tService, appLogger)

	r := router.New(h, &appConf)

	// Запускаем тестовый сервер
	testServer := httptest.NewServer(r)
	defer testServer.Close()

	// Обновление данных
	testTemplate.Title = "Updated title"
	jData, _ := json.Marshal(testTemplate)
	jReader := bytes.NewReader(jData)

	request, _ := http.NewRequest(http.MethodPut, testServer.URL+updateEndpoint, jReader)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		appLogger.Error("http.DefaultClient.Do err", err)
	}

	responseBody, _ := io.ReadAll(response.Body)
	_ = response.Body.Close()

	// В случае успеха сервис отвечает кодом 200 и JSON содержащим текущее значение шаблона
	fmt.Println(response.StatusCode, string(responseBody))

	// Output:
	// 200 {"template_uuid":"c10a7fa8-f162-46ce-a97e-ff8718d7eb7d","event_uuid":"00000000-0000-0000-0000-000000000000","title":"Updated title","description":"Description","body":"Body","channel_type":"ChannelType"}
}

func ExampleHandler_DeleteTemplate() {
	templateUUID, _ := uuid.Parse("c10a7fa8-f162-46ce-a97e-ff8718d7eb7d")

	deleteEndpoint := fmt.Sprintf("/api/v1/templates/%v", templateUUID)

	testTemplate := dto.Template{
		TemplateUUID: templateUUID,
		EventUUID:    uuid.UUID{},
		Title:        "Test",
		Description:  "Description",
		Body:         "Body",
		ChannelType:  "ChannelType",
	}

	// Подготавливаем все зависимости, логгер, конфигурацию приложения, хранилище (In Memory) и роутер
	appLogger := logger.NewZapLogger()
	appConf := mockHandlerConfig{}

	tService := template.New()
	_, _ = tService.Store(context.Background(), testTemplate)

	h := handlers.New(&appConf, nil, nil, nil, tService, appLogger)

	r := router.New(h, &appConf)

	// Запускаем тестовый сервер
	testServer := httptest.NewServer(r)
	defer testServer.Close()

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

func ExampleHandler_GetTemplate() {
	testTemplate := dto.Template{
		EventUUID:   uuid.UUID{},
		Title:       "Test",
		Description: "Description",
		Body:        "Body",
		ChannelType: "ChannelType",
	}

	// Подготавливаем все зависимости, логгер, конфигурацию приложения, хранилище (In Memory) и роутер
	appLogger := logger.NewZapLogger()
	appConf := mockHandlerConfig{}

	tService := template.New()
	stored, _ := tService.Store(context.Background(), testTemplate)

	getEndpoint := fmt.Sprintf("/api/v1/templates/%v", stored.TemplateUUID)

	h := handlers.New(&appConf, nil, nil, nil, tService, appLogger)

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
	var getResult dto.Template
	_ = json.Unmarshal(responseBody, &getResult)
	getResult.TemplateUUID = uuid.UUID{}
	fmt.Println(response.StatusCode, getResult)

	// Output:
	// 200 {00000000-0000-0000-0000-000000000000 00000000-0000-0000-0000-000000000000 Test Description Body ChannelType}
}

func ExampleHandler_GetTemplates() {
	testTemplate := dto.Template{
		EventUUID:   uuid.UUID{},
		Title:       "GetAllTest",
		Description: "Description",
		Body:        "Body",
		ChannelType: "ChannelType",
	}

	testTemplate2 := dto.Template{
		EventUUID:   uuid.UUID{},
		Title:       "Test_2",
		Description: "Description",
		Body:        "Body",
		ChannelType: "ChannelType",
	}

	// Подготавливаем все зависимости, логгер, конфигурацию приложения, хранилище (In Memory) и роутер
	appLogger := logger.NewZapLogger()
	appConf := mockHandlerConfig{}

	tService := template.New()
	_, _ = tService.Store(context.Background(), testTemplate)
	_, _ = tService.Store(context.Background(), testTemplate2)

	getEndpoint := fmt.Sprintf("/api/v1/templates")

	h := handlers.New(&appConf, nil, nil, nil, tService, appLogger)

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
	var getResult []dto.Template
	_ = json.Unmarshal(responseBody, &getResult)

	stableSlice := make([]string, 0, 2)
	for _, t := range getResult {
		checkThis := strings.Builder{}
		checkThis.WriteString(t.Title)
		checkThis.WriteString(t.Body)
		checkThis.WriteString(t.Description)
		checkThis.WriteString(t.ChannelType)
		stableSlice = append(stableSlice, checkThis.String())
	}

	sort.Strings(stableSlice)

	fmt.Println(response.StatusCode, stableSlice)

	// Output:
	// 200 [GetAllTestBodyDescriptionChannelType Test_2BodyDescriptionChannelType]
}

func ExampleHandler_StoreTemplate() {
	storeEndpoint := fmt.Sprintf("/api/v1/templates")

	testTemplate := dto.IncomingTemplate{
		EventUUID:   uuid.UUID{},
		Title:       "TestStoreTemplate",
		Description: "Description",
		Body:        "Body",
		ChannelType: "ChannelType",
	}

	// Подготавливаем все зависимости, логгер, конфигурацию приложения, хранилище (In Memory) и роутер
	appLogger := logger.NewZapLogger()
	appConf := mockHandlerConfig{}

	tService := template.New()

	h := handlers.New(&appConf, nil, nil, nil, tService, appLogger)

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
	var getResult dto.Template
	_ = json.Unmarshal(responseBody, &getResult)

	checkThis := strings.Builder{}
	checkThis.WriteString(getResult.Title)
	checkThis.WriteString(getResult.Description)
	checkThis.WriteString(getResult.Body)
	checkThis.WriteString(getResult.ChannelType)

	// В случае успеха сервис отвечает кодом 200 и JSON содержащим текущее значение шаблона
	fmt.Println(response.StatusCode, checkThis.String())

	// Output:
	// 200 TestStoreTemplateDescriptionBodyChannelType
}

type mockHandlerConfig struct {
}

func (m *mockHandlerConfig) GetTrustedSubnetAddress() string {
	return ""
}

func (m *mockHandlerConfig) GetDefaultResponseContentType() string {
	return "application/json"
}
