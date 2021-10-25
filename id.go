package banking

import (
	"context"
)

// ID represents an unique identifier.
type ID string

func (id ID) String() string {
	return string(id)
}

// EmptyID represents an empty identifier.
const EmptyID = ID("")

// IdentifierGenerator represents a service for generating unique identifier.
type IdentifierGenerator interface {
	// GenerateIdentifier returns unique identifier.
	GenerateIdentifier(ctx context.Context) (ID, error)
}
