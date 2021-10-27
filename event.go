package banking

import (
	"time"
)

// Event represents a paid event or purchase.
type Event struct {
	// ID is the event unique identifier.
	ID ID

	// Name is the event name.
	Name string

	// TemplateID is the unique identifier of EventTemplate used for creating an event.
	// NOTE: TemplateID is optional field. User could create an event without any template.
	TemplateID *ID

	// Creator is the User that creates an event.
	Creator *User

	// Payer is the User who pays for event.
	Payer *User

	// Participants is the list of users among whom the payment is divided.
	Participants []*User

	// CreatedAt is the time when event was created.
	CreatedAt time.Time

	// UpdatedAt is the time when event was updated.
	UpdatedAt time.Time
}

// EventService represents a service for managing Event data.
type EventService interface{}
