package banking

import (
	"time"
)

// UserGroup represents a group of users.
type UserGroup struct {
	// ID is the user group unique identifier.
	ID ID

	// Name is the user group name.
	Name string

	// Users is the collection of users which assigned to the current group.
	Users []*User

	// CreatedAt is the time when user group was created.
	CreatedAt time.Time

	// UpdatedAt is the time when user group was updated.
	UpdatedAt time.Time
}

// UserGroupService represents a service for managing UserGroup data.
type UserGroupService interface{}
