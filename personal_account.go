package banking

import (
	"time"
)

// PersonalAccount represents a user personal account.
type PersonalAccount struct {
	// ID is the unique identifier of user personal account.
	ID ID

	// Owner is the user that owned of personal account.
	Owner *User

	// CreatedAt is the time when personal account was created.
	CreatedAt time.Time
}

// PersonalAccountService represents a service for managing PersonalAccount data.
type PersonalAccountService interface{}
