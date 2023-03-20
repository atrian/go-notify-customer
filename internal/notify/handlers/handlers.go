package handlers

import "github.com/atrian/go-notify-customer/internal/interfaces"

type Handler struct {
	services services
	conf     handlerConfig
	logger   interfaces.Logger
}

type handlerConfig interface {
	GetDefaultResponseContentType() string
}

type services struct {
	event    interfaces.EventService
	notify   interfaces.NotificationService
	stat     interfaces.StatService
	template interfaces.TemplateService
}

func New(
	conf handlerConfig,
	event interfaces.EventService,
	notify interfaces.NotificationService,
	stat interfaces.StatService,
	template interfaces.TemplateService,
	logger interfaces.Logger,
) *Handler {

	h := Handler{
		conf: conf,
		services: services{
			event:    event,
			notify:   notify,
			stat:     stat,
			template: template,
		},
		logger: logger,
	}

	return &h
}
