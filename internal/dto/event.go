package dto

import "github.com/google/uuid"

// Event структура бизнес события для передачи между слоями приложения
type Event struct {
	EventUUID            uuid.UUID `json:"event_uuid"`                      // EventUUID связь с UUID бизнес события
	Title                string    `json:"title"`                           // Title название бизнес события
	Description          string    `json:"description,omitempty"`           // Description описание бизнес события
	DefaultPriority      uint      `json:"default_priority,omitempty"`      // DefaultPriority приоритет уведомления с таким событием по умолчанию
	NotificationChannels []string  `json:"notification_channels,omitempty"` // NotificationChannels каналы отправки для данного события
}

// IncomingEvent структура входящего бизнес события для анмаршаллинга json
type IncomingEvent struct {
	Title                string   `json:"title"`                           // Title название бизнес события
	Description          string   `json:"description,omitempty"`           // Description описание бизнес события
	DefaultPriority      uint     `json:"default_priority,omitempty"`      // DefaultPriority приоритет уведомления с таким событием по умолчанию
	NotificationChannels []string `json:"notification_channels,omitempty"` // NotificationChannels каналы отправки для данного события
}
