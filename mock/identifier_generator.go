package mock

import (
	"context"

	banking "github.com/morozovcookie/agat-banking"
	"github.com/stretchr/testify/mock"
)

var _ banking.IdentifierGenerator = (*IdentifierGenerator)(nil)

// IdentifierGenerator represents a service for generating unique identifier.
type IdentifierGenerator struct {
	mock.Mock
}

// NewIdentifierGenerator returns a new IdentifierGenerator instance.
func NewIdentifierGenerator() *IdentifierGenerator {
	return &IdentifierGenerator{}
}

// GenerateIdentifier returns unique identifier.
func (gen *IdentifierGenerator) GenerateIdentifier(_ context.Context) (banking.ID, error) {
	args := gen.Called()

	return args.Get(0).(banking.ID), args.Error(1)
}
