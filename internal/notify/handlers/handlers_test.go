package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/atrian/go-notify-customer/internal/notify/handlers"
	"github.com/atrian/go-notify-customer/internal/notify/router"
	"github.com/atrian/go-notify-customer/internal/services/event"
	"github.com/atrian/go-notify-customer/pkg/logger"
)

type mockHandlerConfig struct {
}

func (m *mockHandlerConfig) GetTrustedSubnetAddress() string {
	return ""
}

func (m *mockHandlerConfig) GetDefaultResponseContentType() string {
	return "application/json"
}

type subnetConf struct {
	mockHandlerConfig
}

func (s *subnetConf) GetTrustedSubnetAddress() string {
	return "62.217.188.0/24" // random mask
}

func TestMiddlewareSubnetRestriction(t *testing.T) {
	// Подготавливаем все зависимости, логгер, конфигурацию приложения, хранилище (In Memory) и роутер
	appLogger := logger.NewZapLogger()

	appConf := subnetConf{}

	service := event.New()

	getEndpoint := fmt.Sprintf("/api/v1/events")

	h := handlers.New(&appConf, service, nil, nil, nil, appLogger)

	r := router.New(h, &appConf)

	// Запускаем тестовый сервер
	testServer := httptest.NewServer(r)
	defer testServer.Close()

	// Делаем тестовый запрос
	request, _ := http.NewRequest(http.MethodGet, testServer.URL+getEndpoint, nil)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		appLogger.Error("http.DefaultClient.Do err", err)
	}

	// Ожидаем 403 статус т.к. не прошли проверку по подсети
	assert.Equal(t, http.StatusForbidden, response.StatusCode)
}
