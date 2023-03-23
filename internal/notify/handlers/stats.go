package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	statErrors "github.com/atrian/go-notify-customer/internal/services/stat"
)

// GetStats Запрос всех доступных шаблонов GET /api/v1/stats
//
//	@Tags Stat
//	@Summary Запрос всей статистики
//	@Produce json
//	@Success 200 array dto.Stat
//	@Failure 500
//	@Router /api/v1/stats [get]
func (h *Handler) GetStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats := h.services.stat.All(context.Background())

		w.Header().Set("content-type", h.conf.GetDefaultResponseContentType())
		w.WriteHeader(http.StatusOK)

		h.logger.Debug("Request OK")

		jsonEncErr := json.NewEncoder(w).Encode(stats)
		if jsonEncErr != nil {
			h.logger.Error("json.NewEncoder err", jsonEncErr)
		}
	}
}

// GetStatByPersonUUID запрос статистики отправок по пользователю GET /api/v1/stats/person/{UUID-v4}
//
//	@Tags Stat
//	@Summary Запрос статистики отправок по пользователю
//	@Produce json
//	@Param person_uuid path string true "ID пользователя в формате UUID v4"
//	@Success 200 array dto.Stat
//	@Failure 400
//	@Failure 404
//	@Failure 500
//	@Router /api/v1/stats/person/{person_uuid} [get]
func (h *Handler) GetStatByPersonUUID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		param := chi.URLParam(r, "personUUID")

		personUUID, err := uuid.Parse(param)
		if err != nil {
			h.logger.Error("GetStatByPersonUUID Parse personUUID err", err)
			http.Error(w, "Bad personUUID", http.StatusBadRequest)
			return
		}

		stats, err := h.services.stat.FindByPersonUUID(context.Background(), personUUID)

		if err != nil {
			if errors.Is(err, statErrors.NotFound) {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}

			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", h.conf.GetDefaultResponseContentType())
		w.WriteHeader(http.StatusOK)

		h.logger.Debug("Request OK")

		jsonEncErr := json.NewEncoder(w).Encode(stats)
		if jsonEncErr != nil {
			h.logger.Error("json.NewEncoder err", jsonEncErr)
		}
	}
}

// GetStatByNotificationId запрос статистики отправок по уведомлению GET /api/v1/stats/notification/{UUID-v4}
//
//	@Tags Stat
//	@Summary Запрос статистики отправок по уведомлению
//	@Produce json
//	@Param notification_uuid path string true "ID уведомления в формате UUID v4"
//	@Success 200 array dto.Stat
//	@Failure 400
//	@Failure 404
//	@Failure 500
//	@Router /api/v1/stats/notification/{notification_uuid} [get]
func (h *Handler) GetStatByNotificationId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		param := chi.URLParam(r, "notificationUUID")

		notificationUUID, err := uuid.Parse(param)
		if err != nil {
			h.logger.Error("GetStatByNotificationId Parse notificationUUID err", err)
			http.Error(w, "Bad personUUID", http.StatusBadRequest)
			return
		}

		stats, err := h.services.stat.FindByNotificationId(context.Background(), notificationUUID)

		if err != nil {
			if errors.Is(err, statErrors.NotFound) {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}

			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", h.conf.GetDefaultResponseContentType())
		w.WriteHeader(http.StatusOK)

		h.logger.Debug("Request OK")

		jsonEncErr := json.NewEncoder(w).Encode(stats)
		if jsonEncErr != nil {
			h.logger.Error("json.NewEncoder err", jsonEncErr)
		}
	}
}
