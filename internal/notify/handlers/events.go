package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/atrian/go-notify-customer/internal/dto"
	eventErrors "github.com/atrian/go-notify-customer/internal/services/event"
)

// UpdateEvent обновление бизнес события PUT /api/v1/events/{UUID-v4}
//
//	@Tags Event
//	@Summary обновление бизнес события
//	@Accept  json
//	@Produce json
//	@Param event_uuid path string true "ID бизнес события в формате UUID v4"
//	@Param event body dto.IncomingEvent true "Test"
//	@Success 200 {object} dto.Event
//	@Failure 400
//	@Failure 404
//	@Failure 500
//	@Router /api/v1/events/{event_uuid} [put]
func (h *Handler) UpdateEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		param := chi.URLParam(r, "eventUUID")

		event, err := h.unmarshallIncomingEvent(r)

		if err != nil {
			h.logger.Error("UpdateEvent cant unmarshallEvent", err)
			http.Error(w, "Bad JSON", http.StatusBadRequest)
			return
		}

		eventUUID, err := uuid.Parse(param)
		if err != nil {
			h.logger.Error("UpdateEvent cant Parse UUID", err)
			http.Error(w, "Bad UUID", http.StatusBadRequest)
			return
		}

		eventForUpdate := dto.Event{
			EventUUID:            eventUUID,
			Title:                event.Title,
			Description:          event.Description,
			DefaultPriority:      event.DefaultPriority,
			NotificationChannels: event.NotificationChannels,
		}

		result, err := h.services.event.Update(context.Background(), eventForUpdate)

		if err != nil {
			if errors.Is(err, eventErrors.NotFound) {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}

			http.Error(w, "Bad JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", h.conf.GetDefaultResponseContentType())
		w.WriteHeader(http.StatusOK)

		h.logger.Debug("Request OK")

		jsonEncErr := json.NewEncoder(w).Encode(result)
		if jsonEncErr != nil {
			h.logger.Error("json.NewEncoder err", jsonEncErr)
		}
	}
}

// StoreEvent сохранение бизнес события POST /api/v1/events
//
//	@Tags Event
//	@Summary сохранение бизнес события
//	@Accept  json
//	@Produce json
//	@Param event body dto.IncomingEvent true "Принимает dto события, отдает сохраненное событие с идентификатором"
//	@Success 200 {object} dto.Event
//	@Failure 400
//	@Failure 500
//	@Router /api/v1/events [post]
func (h *Handler) StoreEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		event, err := h.unmarshallEvent(r)
		if err != nil {
			h.logger.Error("StoreEvent cant unmarshallEvent", err)
			http.Error(w, "Bad JSON", http.StatusBadRequest)
			return
		}

		result, err := h.services.event.Store(context.Background(), event)

		if err != nil {
			http.Error(w, "Bad JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", h.conf.GetDefaultResponseContentType())
		w.WriteHeader(http.StatusOK)

		h.logger.Debug("Request OK")

		jsonEncErr := json.NewEncoder(w).Encode(result)
		if jsonEncErr != nil {
			h.logger.Error("json.NewEncoder err", jsonEncErr)
		}
	}
}

// DeleteEvent удаление бизнес события DELETE /api/v1/events/{UUID-v4}
//
//	@Tags Event
//	@Summary удаление бизнес события
//	@Produce json
//	@Param event_uuid path string true "ID бизнес события в формате UUID v4"
//	@Success 200
//	@Failure 400
//	@Failure 404
//	@Failure 500
//	@Router /api/v1/events/{event_uuid} [delete]
func (h *Handler) DeleteEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		param := chi.URLParam(r, "eventUUID")

		eventUUID, err := uuid.Parse(param)
		if err != nil {
			h.logger.Error("DeleteEvent Parse eventUUID", err)
			http.Error(w, "Bad eventUUID", http.StatusBadRequest)
			return
		}

		err = h.services.event.DeleteById(context.Background(), eventUUID)

		if err != nil {
			if errors.Is(err, eventErrors.NotFound) {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}

			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", h.conf.GetDefaultResponseContentType())
		w.WriteHeader(http.StatusOK)

		h.logger.Debug("Request OK")
	}
}

// GetEvent Запрос деталей бизнес события GET /api/v1/events/{UUID-v4}
//
//	@Tags Event
//	@Summary Запрос деталей бизнес события
//	@Produce json
//	@Param event_uuid path string true "ID события в формате UUID v4"
//	@Success 200 {object} dto.Event
//	@Failure 400
//	@Failure 404
//	@Failure 500
//	@Router /api/v1/events/{event_uuid} [get]
func (h *Handler) GetEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		param := chi.URLParam(r, "eventUUID")

		eventUUID, err := uuid.Parse(param)
		if err != nil {
			h.logger.Error("GetEvent Parse eventUUID", err)
			http.Error(w, "Bad eventUUID", http.StatusBadRequest)
			return
		}

		event, err := h.services.event.FindById(context.Background(), eventUUID)

		if err != nil {
			if errors.Is(err, eventErrors.NotFound) {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}

			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", h.conf.GetDefaultResponseContentType())
		w.WriteHeader(http.StatusOK)

		h.logger.Debug("Request OK")

		jsonEncErr := json.NewEncoder(w).Encode(event)
		if jsonEncErr != nil {
			h.logger.Error("json.NewEncoder err", jsonEncErr)
		}
	}
}

// GetEvents Запрос всех доступных шаблонов GET /api/v1/events
//
//	@Tags Event
//	@Summary Запрос всех доступных шаблонов
//	@Produce json
//	@Success 200 array dto.Event
//	@Failure 500
//	@Router /api/v1/events [get]
func (h *Handler) GetEvents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		events := h.services.event.All(context.Background())

		w.Header().Set("content-type", h.conf.GetDefaultResponseContentType())
		w.WriteHeader(http.StatusOK)

		h.logger.Debug("Request OK")

		jsonEncErr := json.NewEncoder(w).Encode(events)
		if jsonEncErr != nil {
			h.logger.Error("json.NewEncoder err", jsonEncErr)
		}
	}
}

// unmarshallEvent анмаршаллинг бизнес события
func (h *Handler) unmarshallEvent(r *http.Request) (dto.Event, error) {
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

	var event dto.Event
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&event)
	if err != nil {
		return dto.Event{}, err
	}

	return event, nil
}

// unmarshallEvent анмаршаллинг бизнес события
func (h *Handler) unmarshallIncomingEvent(r *http.Request) (dto.IncomingEvent, error) {
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

	var event dto.IncomingEvent
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&event)
	if err != nil {
		return dto.IncomingEvent{}, err
	}

	return event, nil
}
