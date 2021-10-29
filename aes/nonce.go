package aes

import (
	"context"
)

// NonceGenerator represents a service for generating nonce.
type NonceGenerator interface {
	// GenerateNonce returns a nonce.
	GenerateNonce(ctx context.Context, size int) ([]byte, error)
}
