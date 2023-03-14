package entity

import "github.com/google/uuid"

// Notification структура уведомления
type Notification struct {
	EventUUID     uuid.UUID      `json:"event_uuid"`               // EventUUID связь с UUID бизнес события
	PersonUUIDs   []uuid.UUID    `json:"person_uuids"`             // PersonUUIDs связь с пользователями - получателями уведомления
	MessageParams []MessageParam `json:"message_params,omitempty"` // MessageParams key-value подстановки в шаблон уведомления
	Priority      uint           `json:"priority,omitempty"`       // Priority опциональный приоритет уведомления
}

// MessageParam key-value подстановки в шаблон уведомления
type MessageParam struct {
	Key   string `json:"key"`   // Key ключ по которому будет произведен поиск в теле уведомления
	Value string `json:"value"` // Value значение которое будет подставлено вместо ключа в шаблоне
}
