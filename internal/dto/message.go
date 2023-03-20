package dto

import "github.com/google/uuid"

type Message struct {
	NotificationUUID   uuid.UUID `json:"notification_uuid"`
	PersonUUID         uuid.UUID `json:"person_uuid"`
	Text               string    `json:"text"`
	Channel            string    `json:"channel"`
	DestinationAddress string    `json:"destination_address"`
}
