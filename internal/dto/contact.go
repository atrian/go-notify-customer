package dto

import "github.com/google/uuid"

type PersonContacts struct {
	PersonUUID uuid.UUID `json:"person_uuid"`
	Contacts   []Contact `json:"contacts,omitempty"`
}

type Contact struct {
	Channel     string `json:"channel"`
	Destination string `json:"destination"`
}
