package mock

import (
	"context"

	"github.com/morozovcookie/agat-banking/aes"
	"github.com/stretchr/testify/mock"
)

var _ aes.NonceGenerator = (*NonceGenerator)(nil)

// NonceGenerator represents a service for generating nonce.
type NonceGenerator struct {
	mock.Mock
}

// NewNonceGenerator returns a new NonceGenerator instance.
func NewNonceGenerator() *NonceGenerator {
	return &NonceGenerator{}
}

// GenerateNonce returns a nonce.
func (gen *NonceGenerator) GenerateNonce(_ context.Context, size int) ([]byte, error) {
	args := gen.Called(size)

	return args.Get(0).([]byte), args.Error(1)
}
