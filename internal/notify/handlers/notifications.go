package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/atrian/go-notify-customer/internal/dto"
	"github.com/atrian/go-notify-customer/internal/services/notify"
)

// ProcessNotifications отправка уведомлений POST /api/v1/notifications
//
//	@Tags Notifications
//	@Summary отправка уведомлений
//	@Accept  json
//	@Produce json
//	@Param metrics body array dto.Notification true
//	@Success 200
//	@Failure 400
//	@Failure 429
//	@Failure 500
//	@Router /api/v1/notifications [post]
func (h *Handler) ProcessNotifications() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		notifications, err := h.unmarshallNotifications(r)
		if err != nil {
			h.logger.Error("ProcessNotifications cant unmarshallNotifications", err)
			http.Error(w, "Bad JSON", http.StatusBadRequest)
		}

		err = h.services.notify.ProcessNotification(notifications)
		if err != nil {
			if errors.Is(err, notify.NotificationLimitExceeded) {
				http.Error(w, "limit exceed", http.StatusTooManyRequests)
				return
			}

			http.Error(w, "Bad JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", h.conf.GetDefaultResponseContentType())
		w.WriteHeader(http.StatusOK)

		h.logger.Debug("Request OK")
	}
}

// unmarshallNotifications анмаршаллинг уведомлений
func (h *Handler) unmarshallNotifications(r *http.Request) ([]dto.Notification, error) {
	var body io.Reader

	// если в заголовках установлен Content-Encoding gzip, распаковываем тело
	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		body = h.decodeGzipBody(r.Body)
	} else {
		body = r.Body
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.logger.Error("Body io.ReadCloser error", err)
		}
	}(r.Body)

	var notifications []dto.Notification
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&notifications)
	if err != nil {
		return nil, err
	}

	return notifications, nil
}
