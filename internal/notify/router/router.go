package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/atrian/go-notify-customer/docs"
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
		// Swagger
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
		))

		r.Route("/api/v1", func(r chi.Router) {

			// Сервис шаблонов сообщений
			r.Route("/templates", func(r chi.Router) {
				r.Get("/", handler.GetTemplates())   // GET /templates
				r.Post("/", handler.StoreTemplate()) // POST /templates

				r.Route("/{templateUUID}", func(r chi.Router) {
					// GET /templates/93ebac94-cf39-4728-9bba-472ac93a4368
					r.Get("/", handler.GetTemplate())
					// PUT /templates/93ebac94-cf39-4728-9bba-472ac93a4368
					r.Put("/", handler.UpdateTemplate())
					// DELETE /templates/93ebac94-cf39-4728-9bba-472ac93a4368
					r.Delete("/", handler.DeleteTemplate())
				})
			})

			// Сервис бизнес событий
			r.Route("/events", func(r chi.Router) {
				r.Get("/", handler.GetEvents())   // GET /events
				r.Post("/", handler.StoreEvent()) // POST /events

				r.Route("/{eventUUID}", func(r chi.Router) {
					// GET /events/93ebac94-cf39-4728-9bba-472ac93a4368
					r.Get("/", handler.GetEvent())
					// PUT /events/93ebac94-cf39-4728-9bba-472ac93a4368
					r.Put("/", handler.UpdateEvent())
					// DELETE /events/93ebac94-cf39-4728-9bba-472ac93a4368
					r.Delete("/", handler.DeleteEvent())
				})
			})

			// Сервис статистики
			r.Route("/stats", func(r chi.Router) {
				// GET /stats
				r.Get("/", handler.GetStats())
				// GET /stats/person/{personUUID}
				r.Get("/person/{personUUID}", handler.GetStatByPersonUUID())
				// GET /stats/person/{notificationUUID}
				r.Get("/notification/{notificationUUID}", handler.GetStatByNotificationId())
			})

			// Сервис уведомлений
			r.Route("/notifications", func(r chi.Router) {
				r.Post("/", handler.ProcessNotifications())
				r.Get("/seed", handler.SeedDemoData())
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
