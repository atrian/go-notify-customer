package entity

import "github.com/google/uuid"

type Message struct {
	PersonUUID         uuid.UUID `json:"person_uuid"`
	Text               string    `json:"text"`
	Channel            string    `json:"channel"`
	DestinationAddress string    `json:"destination_address"`
}
