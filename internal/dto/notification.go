package dto

import "github.com/google/uuid"

// Notification структура уведомления для внутренних интерфейсов
type Notification struct {
	Index            int            // Index индекс уведомления в очереди
	NotificationUUID uuid.UUID      `json:"notification_uuid"`        // NotificationUUID id уведомления в системе
	EventUUID        uuid.UUID      `json:"event_uuid"`               // EventUUID связь с UUID бизнес события
	PersonUUIDs      []uuid.UUID    `json:"person_uuids"`             // PersonUUIDs связь с пользователями - получателями уведомления
	MessageParams    []MessageParam `json:"message_params,omitempty"` // MessageParams key-value подстановки в шаблон уведомления
	Priority         uint           `json:"priority,omitempty"`       // Priority опциональный приоритет уведомления
}

// IncomingNotification структура уведомления для внешних интерфейсов
type IncomingNotification struct {
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
