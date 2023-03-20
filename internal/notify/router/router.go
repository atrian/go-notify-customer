package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/atrian/go-notify-customer/internal/notify/handlers"
	customMiddleware "github.com/atrian/go-notify-customer/internal/notify/middleware"
)

type Router struct {
	*chi.Mux
	conf securityConfig
}

type securityConfig interface {
	GetTrustedSubnetAddress() string
}

// RegisterMiddlewares общие middlewares для всех маршрутов
// Вызывать ДО регистрации маршрутов
func (r *Router) RegisterMiddlewares() *Router {
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	return r
}

// RegisterRoutes регистрация всех маршрутов бизнес логики приложения
// Вызывать ПОСЛЕ регистрации всех middlewares
func (r *Router) RegisterRoutes(handler *handlers.Handler) *Router {
	// Конфигурируем MW ограничения соединений для доверенных сетей
	trustedMW := customMiddleware.TrustedSubnetMW(r.conf.GetTrustedSubnetAddress())

	r.Group(func(r chi.Router) {
		r.Use(trustedMW)

		r.Route("/api/v1", func(r chi.Router) {

			// Сервис шаблонов сообщений
			r.Route("/templates", func(r chi.Router) {
				r.Get("/", nil)  // GET /templates
				r.Post("/", nil) // POST /templates

				r.Route("/{templateUUID}", func(r chi.Router) {
					// GET /templates/93ebac94-cf39-4728-9bba-472ac93a4368
					r.Get("/", nil)
					// PUT /templates/93ebac94-cf39-4728-9bba-472ac93a4368
					r.Put("/", handler.UpdateTemplate())
					// DELETE /templates/93ebac94-cf39-4728-9bba-472ac93a4368
					r.Delete("/", handler.DeleteTemplate())
				})
			})

		})
	})

	return r
}

// New возвращает роутер со стандартной конфигурацией.
// Принимает слайс дополнительных кастомных middleware
func New(handler *handlers.Handler, config securityConfig) *Router {
	router := Router{
		Mux:  chi.NewMux(),
		conf: config,
	}

	// middlewares
	router.RegisterMiddlewares()

	// routes
	router.RegisterRoutes(handler)

	return &router
}
