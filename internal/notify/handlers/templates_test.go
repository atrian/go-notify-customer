package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

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

type mockHandlerConfig struct {
}

func (m *mockHandlerConfig) GetTrustedSubnetAddress() string {
	return ""
}

func (m *mockHandlerConfig) GetDefaultResponseContentType() string {
	return "application/json"
}
