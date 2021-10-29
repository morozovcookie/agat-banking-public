package aes

import (
	banking "github.com/morozovcookie/agat-banking"
)

var _ banking.SecretString = (*SecretString)(nil)

// SecretString represents an object that stores sensitive information.
type SecretString struct {
	encryptedString string
	decryptedString string
}

// EncryptedString returns sensitive information with encrypted string.
func (ss *SecretString) EncryptedString() string {
	return ss.encryptedString
}

// DecryptedString returns sensitive information with decrypted string.
func (ss *SecretString) DecryptedString() string {
	return ss.decryptedString
}

func (ss *SecretString) String() string {
	return ss.encryptedString
}
