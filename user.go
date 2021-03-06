package banking

import (
	"time"
)

// User represents an information set about user.
type User struct {
	// ID is the user unique identifier.
	ID ID

	// FirstName is the user first name.
	FirstName string

	// LastName is the user last name.
	LastName string

	// CreatedAt is the time when user was created.
	CreatedAt time.Time

	// UpdatedAt is the time when user was updated.
	UpdatedAt time.Time
}
