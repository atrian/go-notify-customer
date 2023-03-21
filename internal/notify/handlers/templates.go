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
	templateErrors "github.com/atrian/go-notify-customer/internal/services/template"
)

// UpdateTemplate обновление шаблона сообщения PUT /api/v1/templates/{UUID-v4}
//
//	@Tags Template
//	@Summary обновление шаблона сообщения
//	@Accept  json
//	@Produce json
//	@Param templates_uuid path string true "ID шаблона сообщения в формате UUID v4"
//	@Param template body dto.IncomingTemplate true "Принимает dto шаблона сообщения, возвращает JSON с обновленными данными"
//	@Success 200 {object} dto.Template
//	@Failure 400
//	@Failure 404
//	@Failure 500
//	@Router /api/v1/templates/{templates_uuid} [put]
func (h *Handler) UpdateTemplate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		param := chi.URLParam(r, "templateUUID")
		templateUUID, err := uuid.Parse(param)
		if err != nil {
			h.logger.Error("UpdateTemplate cant Parse UUID", err)
			http.Error(w, "Bad template UUID", http.StatusBadRequest)
			return
		}

		template, err := h.unmarshallIncomingTemplate(r)
		if err != nil {
			h.logger.Error("UpdateTemplate cant unmarshallTemplate", err)
			http.Error(w, "Bad JSON", http.StatusBadRequest)
			return
		}

		updateTemplate := dto.Template{
			TemplateUUID: templateUUID,
			EventUUID:    template.EventUUID,
			Title:        template.Title,
			Description:  template.Description,
			Body:         template.Body,
			ChannelType:  template.ChannelType,
		}

		result, err := h.services.template.Update(context.Background(), updateTemplate)

		if err != nil {
			if errors.Is(err, templateErrors.NotFound) {
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

// StoreTemplate сохранение шаблона сообщения POST /api/v1/templates
//
//	@Tags Template
//	@Summary сохранение шаблона сообщения
//	@Accept  json
//	@Produce json
//	@Param template body dto.IncomingTemplate true "Принимает dto нового шаблона сообщения, возвращает JSON сохраненными данными и идентификатором"
//	@Success 200 {object} dto.Template
//	@Failure 400
//	@Failure 500
//	@Router /api/v1/templates [post]
func (h *Handler) StoreTemplate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		template, err := h.unmarshallIncomingTemplate(r)
		if err != nil {
			h.logger.Error("StoreTemplate cant unmarshallTemplate", err)
			http.Error(w, "Bad JSON", http.StatusBadRequest)
			return
		}

		storeTemplate := dto.Template{
			EventUUID:   template.EventUUID,
			Title:       template.Title,
			Description: template.Description,
			Body:        template.Body,
			ChannelType: template.ChannelType,
		}

		result, err := h.services.template.Store(context.Background(), storeTemplate)

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

// DeleteTemplate удаление шаблона сообщения DELETE /api/v1/templates/{UUID-v4}
//
//	@Tags Template
//	@Summary удаление шаблона сообщения
//	@Produce json
//	@Param template_uuid path string true "ID шаблона в формате UUID v4"
//	@Success 200
//	@Failure 400
//	@Failure 404
//	@Failure 500
//	@Router /api/v1/templates/{template_uuid} [delete]
func (h *Handler) DeleteTemplate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		param := chi.URLParam(r, "templateUUID")

		templateUUID, err := uuid.Parse(param)
		if err != nil {
			h.logger.Error("DeleteTemplate Parse templateUUID", err)
			http.Error(w, "Bad templateUUID", http.StatusBadRequest)
			return
		}

		err = h.services.template.DeleteById(context.Background(), templateUUID)

		if err != nil {
			if errors.Is(err, templateErrors.NotFound) {
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

// GetTemplate Запрос деталей шаблона сообщения GET /api/v1/templates/{UUID-v4}
//
//	@Tags Template
//	@Summary Запрос деталей шаблона сообщения
//	@Produce json
//	@Param template_uuid path string true "ID шаблона в формате UUID v4"
//	@Success 200 {object} dto.Template
//	@Failure 400
//	@Failure 404
//	@Failure 500
//	@Router /api/v1/templates/{template_uuid} [get]
func (h *Handler) GetTemplate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		param := chi.URLParam(r, "templateUUID")

		templateUUID, err := uuid.Parse(param)
		if err != nil {
			h.logger.Error("GetTemplate Parse templateUUID", err)
			http.Error(w, "Bad templateUUID", http.StatusBadRequest)
			return
		}

		template, err := h.services.template.FindById(context.Background(), templateUUID)

		if err != nil {
			if errors.Is(err, templateErrors.NotFound) {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}

			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", h.conf.GetDefaultResponseContentType())
		w.WriteHeader(http.StatusOK)

		h.logger.Debug("Request OK")

		jsonEncErr := json.NewEncoder(w).Encode(template)
		if jsonEncErr != nil {
			h.logger.Error("json.NewEncoder err", jsonEncErr)
		}
	}
}

// GetTemplates Запрос всех доступных шаблонов GET /api/v1/templates
//
//	@Tags Template
//	@Summary Запрос деталей шаблона сообщения
//	@Produce json
//	@Success 200 array dto.Template
//	@Failure 500
//	@Router /api/v1/templates [get]
func (h *Handler) GetTemplates() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templates := h.services.template.All(context.Background())

		w.Header().Set("content-type", h.conf.GetDefaultResponseContentType())
		w.WriteHeader(http.StatusOK)

		h.logger.Debug("Request OK")

		jsonEncErr := json.NewEncoder(w).Encode(templates)
		if jsonEncErr != nil {
			h.logger.Error("json.NewEncoder err", jsonEncErr)
		}
	}
}

// unmarshallTemplate анмаршаллинг шаблона сообщения
func (h *Handler) unmarshallTemplate(r *http.Request) (dto.Template, error) {
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

	var template dto.Template
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&template)
	if err != nil {
		return dto.Template{}, err
	}

	return template, nil
}

// unmarshallIncomingTemplate анмаршаллинг шаблона нового сообщения
func (h *Handler) unmarshallIncomingTemplate(r *http.Request) (dto.IncomingTemplate, error) {
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

	var template dto.IncomingTemplate
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&template)
	if err != nil {
		return dto.IncomingTemplate{}, err
	}

	return template, nil
}
