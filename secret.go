package banking

import (
	"context"
	"fmt"
	"io"
)

// SecretString represents an object that stores sensitive information.
type SecretString interface {
	fmt.Stringer

	// EncryptedString returns sensitive information with encrypted string.
	EncryptedString() string

	// DecryptedString returns sensitive information with decrypted string.
	DecryptedString() string
}

// SecretFactory represents a service initialize SecretString object.
type SecretFactory interface {
	// CreateFromEncryptedData creates SecretString object from encrypted data.
	CreateFromEncryptedData(ctx context.Context, r io.Reader) (SecretString, error)

	// CreateFromDecryptedData creates SecretString object from decrypted data.
	CreateFromDecryptedData(ctx context.Context, r io.Reader) (SecretString, error)
}
