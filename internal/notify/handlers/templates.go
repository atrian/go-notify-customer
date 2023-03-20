package handlers

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/atrian/go-notify-customer/internal/dto"
	templateErrors "github.com/atrian/go-notify-customer/internal/services/template"
)

// UpdateTemplate обновление шаблона сообщения PUT /api/v1/templates/{UUID-v4}
//
//	@Tags Template
//	@Summary обновление шаблона сообщения
//	@Accept  json
//	@Produce json
//	@Param metrics body dto.Template true
//	@Success 200 dto.Template
//	@Failure 400
//	@Failure 404
//	@Failure 500
//	@Router /api/v1/templates/{UUID-v4} [put]
func (h *Handler) UpdateTemplate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		template, err := h.unmarshallTemplate(r)
		if err != nil {
			h.logger.Error("UpdateJSONMetrics cant unmarshallMetric", err)
			http.Error(w, "Bad JSON", http.StatusBadRequest)
		}

		result, err := h.services.template.Update(context.Background(), template)

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

// decodeGzipBody распаковка GZIP тела запроса
func (h *Handler) decodeGzipBody(gzipR io.Reader) io.Reader {
	gz, err := gzip.NewReader(gzipR)
	if err != nil {
		h.logger.Error("decodeGzipBody cant set up gzip decoder", err)
	}
	return gz
}
