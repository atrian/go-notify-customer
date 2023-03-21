package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/atrian/go-notify-customer/internal/dto"
	"github.com/atrian/go-notify-customer/internal/services/notify"
	"github.com/google/uuid"
	"io"
	"net/http"
	"strings"
)

// ProcessNotifications отправка уведомлений POST /api/v1/notifications
//
//	@Tags Notifications
//	@Summary отправка уведомлений
//	@Accept  json
//	@Produce json
//	@Param notification body []dto.IncomingNotification true "Принимает JSON dto уведомлений, возвращает код 200 при успешной постановке, 429 при привышении лимита"
//	@Success 200
//	@Failure 400
//	@Failure 429
//	@Failure 500
//	@Router /api/v1/notifications [post]
func (h *Handler) ProcessNotifications() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		incomingNotifications, err := h.unmarshallIncomingNotifications(r)
		if err != nil {
			h.logger.Error("ProcessNotifications cant unmarshallIncomingNotifications", err)
			http.Error(w, "Bad JSON", http.StatusBadRequest)
			return
		}

		notifications := make([]dto.Notification, 0, len(incomingNotifications))

		for _, n := range incomingNotifications {
			notifications = append(notifications, dto.Notification{
				EventUUID:     n.EventUUID,
				PersonUUIDs:   n.PersonUUIDs,
				MessageParams: n.MessageParams,
				Priority:      n.Priority,
			})
		}

		err = h.services.notify.ProcessNotification(notifications)
		if err != nil {
			if errors.Is(err, notify.NotificationLimitExceeded) {
				http.Error(w, "limit exceed", http.StatusTooManyRequests)
				return
			}

			http.Error(w, "Server side error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", h.conf.GetDefaultResponseContentType())
		w.WriteHeader(http.StatusOK)

		h.logger.Debug("Request OK")
	}
}

// SeedDemoData Создает бизнес событие и шаблон к нему, возвращает подготовленный JSON для запроса
// через ProcessNotifications в POST /api/v1/notifications
// GET /api/v1/notifications/seed
//
//	@Tags Notifications
//	@Summary Создает бизнес событие и шаблон к нему, возвращает подготовленный JSON для запроса POST /api/v1/notifications
//	@Accept  json
//	@Produce json
//	@Success 200 {object} dto.IncomingNotification
//	@Failure 500
//	@Router /api/v1/notifications/seed [get]
func (h *Handler) SeedDemoData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		channel := "mail"

		// создаем бизнес событие
		event := dto.Event{
			Title:           "[demo] New appointment",
			Description:     "[demo] Event description",
			DefaultPriority: 10,
			NotificationChannels: []string{
				channel,
			},
		}

		event, err := h.services.event.Store(context.Background(), event)
		if err != nil {
			h.logger.Error("SeedDemoData h.services.event.Store err", err)
			http.Error(w, "Server side error", http.StatusInternalServerError)
			return
		}

		// создаем шаблон к нему
		template := dto.Template{
			EventUUID:   event.EventUUID,
			Title:       "[demo] Template title",
			Description: "[demo] Template description",
			Body:        "Demo message with two placeholders for date:[date] and service:[service]",
			ChannelType: channel,
		}

		template, err = h.services.template.Store(context.Background(), template)
		if err != nil {
			h.logger.Error("SeedDemoData h.services.template.Store err", err)
			http.Error(w, "Server side error", http.StatusInternalServerError)
			return
		}

		// создаем структуру тела ответа
		responseJSON := dto.IncomingNotification{
			EventUUID: event.EventUUID,
			PersonUUIDs: []uuid.UUID{
				uuid.New(),
			},
			MessageParams: []dto.MessageParam{
				{
					Key:   "date",
					Value: "21.03.2023",
				}, {
					Key:   "service",
					Value: "Demo service",
				},
			},
			Priority: event.DefaultPriority,
		}

		w.Header().Set("content-type", h.conf.GetDefaultResponseContentType())
		w.WriteHeader(http.StatusOK)

		jsonEncErr := json.NewEncoder(w).Encode([]dto.IncomingNotification{responseJSON})
		if jsonEncErr != nil {
			h.logger.Error("json.NewEncoder err", jsonEncErr)
		}

		h.logger.Debug("Seed OK")
	}
}

// unmarshallIncomingNotifications анмаршаллинг уведомлений
func (h *Handler) unmarshallIncomingNotifications(r *http.Request) ([]dto.IncomingNotification, error) {
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

	var notifications []dto.IncomingNotification
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&notifications)
	if err != nil {
		return nil, err
	}

	return notifications, nil
}
