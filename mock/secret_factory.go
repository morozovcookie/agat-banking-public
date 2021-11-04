package mock

import (
	"context"
	"io"

	banking "github.com/morozovcookie/agat-banking"
	"github.com/stretchr/testify/mock"
)

var _ banking.SecretFactory = (*SecretFactory)(nil)

// SecretFactory represents a service initialize SecretString object.
type SecretFactory struct {
	mock.Mock
}

// NewSecretFactory returns a new SecretFactory instance.
func NewSecretFactory() *SecretFactory {
	return &SecretFactory{}
}

// CreateFromEncryptedData creates SecretString object from encrypted data.
func (factory *SecretFactory) CreateFromEncryptedData(_ context.Context, r io.Reader) (banking.SecretString, error) {
	args := factory.Called(r)

	return args.Get(0).(banking.SecretString), args.Error(1)
}

// CreateFromDecryptedData creates SecretString object from decrypted data.
func (factory *SecretFactory) CreateFromDecryptedData(_ context.Context, r io.Reader) (banking.SecretString, error) {
	args := factory.Called(r)

	return args.Get(0).(banking.SecretString), args.Error(1)
}
