package rand

import (
	"context"
	"crypto/rand"
	"io"

	"github.com/morozovcookie/agat-banking/aes"
	"github.com/pkg/errors"
)

var _ aes.NonceGenerator = (*NonceGenerator)(nil)

// NonceGenerator represents a service for generating nonce by using rand.Reader.
type NonceGenerator struct{}

// NewNonceGenerator returns a new RandomNonceGenerator instance.
func NewNonceGenerator() *NonceGenerator {
	return &NonceGenerator{}
}

// GenerateNonce returns a nonce.
func (gen *NonceGenerator) GenerateNonce(_ context.Context, size int) ([]byte, error) {
	nonce := make([]byte, size)

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, errors.Wrap(err, "generate nonce")
	}

	return nonce, nil
}
