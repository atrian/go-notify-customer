package dto

import "github.com/google/uuid"

type Stat struct {
	PersonUUID       uuid.UUID  `json:"person_uuid"`       // PersonUUID связь отправленного уведомления с клиентом
	NotificationUUID uuid.UUID  `json:"notification_uuid"` // NotificationUUID связь с уведомлением
	CreatedAt        string     `json:"created_at"`        // CreatedAt дата и время отправки
	Status           StatStatus `json:"status"`            // Status статус отправки
}

// StatStatus Статусы обработки заказа
type StatStatus int64

const (
	Sent   StatStatus = iota + 1 // Уведомление отправлено
	Failed                       // Ошибка отправки
)

func (s StatStatus) String() string {
	switch s {
	case Sent:
		return "sent"
	case Failed:
		return "failed"
	}
	return "unknown"
}
