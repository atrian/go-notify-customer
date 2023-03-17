package dto

import "github.com/google/uuid"

type Template struct {
	TemplateUUID uuid.UUID `json:"template_uuid"`         // TemplateUUID - id шаблона
	EventUUID    uuid.UUID `json:"event_uuid"`            // EventUUID связь с UUID бизнес события
	Title        string    `json:"title"`                 // Title название шаблона
	Description  string    `json:"description,omitempty"` // Description описание шаблона
	Body         string    `json:"body"`                  // Body тело шаблона
	ChannelType  string    `json:"channel_type"`          // ChannelType связь с каналом отправки
}
