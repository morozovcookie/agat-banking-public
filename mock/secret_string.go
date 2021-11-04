package mock

import (
	banking "github.com/morozovcookie/agat-banking"
	"github.com/stretchr/testify/mock"
)

var _ banking.SecretString = (*SecretString)(nil)

// SecretString represents an object that stores sensitive information.
type SecretString struct {
	mock.Mock
}

// NewSecretString returns a new SecretString instance.
func NewSecretString() *SecretString {
	return &SecretString{}
}

func (ss *SecretString) String() string {
	return ss.Called().String(0)
}

// EncryptedString returns sensitive information with encrypted string.
func (ss *SecretString) EncryptedString() string {
	return ss.Called().String(0)
}

// DecryptedString returns sensitive information with decrypted string.
func (ss *SecretString) DecryptedString() string {
	return ss.Called().String(0)
}
