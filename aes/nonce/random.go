package nonce

import (
	"context"
	"crypto/rand"
	"io"

	"github.com/morozovcookie/agat-banking/aes"
	"github.com/pkg/errors"
)

var _ aes.NonceGenerator = (*RandomNonceGenerator)(nil)

// RandomNonceGenerator represents a service for generating nonce by using rand.Reader.
type RandomNonceGenerator struct{}

// NewRandomNonceGenerator returns a new RandomNonceGenerator instance.
func NewRandomNonceGenerator() *RandomNonceGenerator {
	return &RandomNonceGenerator{}
}

// GenerateNonce returns a nonce.
func (gen *RandomNonceGenerator) GenerateNonce(_ context.Context, size int) ([]byte, error) {
	nonce := make([]byte, size)

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, errors.Wrap(err, "generate nonce")
	}

	return nonce, nil
}
