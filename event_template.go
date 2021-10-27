package banking

import (
	"time"
)

// EventTemplate represents a template for creating Event.
type EventTemplate struct {
	// ID is the event template unique identifier.
	ID ID

	// Name is the event template name.
	Name string

	// CreatedAt is the time when event template was created.
	CreatedAt time.Time

	// UpdatedAt is the time when event template was updated.
	UpdatedAt time.Time
}

// EventTemplateService represents a service fot managing EventTemplate data.
type EventTemplateService interface{}
