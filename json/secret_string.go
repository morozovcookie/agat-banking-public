package json

import (
	"bytes"
	"context"
	"encoding/json"

	banking "github.com/morozovcookie/agat-banking"
	"github.com/pkg/errors"
)

var (
	_ banking.SecretString = (*SecretString)(nil)
	_ json.Marshaler       = (*SecretString)(nil)
	_ json.Unmarshaler     = (*SecretString)(nil)
)

// SecretString represents an object that stores sensitive information.
type SecretString struct {
	wrapped       banking.SecretString
	secretFactory banking.SecretFactory
}

// NewSecretString returns a new SecretString instance.
func NewSecretString(wrapped banking.SecretString, secretFactory banking.SecretFactory) *SecretString {
	return &SecretString{
		wrapped:       wrapped,
		secretFactory: secretFactory,
	}
}

func (ss *SecretString) UnmarshalJSON(bb []byte) (err error) {
	var str string

	if err = json.NewDecoder(bytes.NewBuffer(bb)).Decode(&str); err != nil {
		return errors.Wrap(err, "unmarshal SecretString")
	}

	ss.wrapped, err = ss.secretFactory.CreateFromDecryptedData(context.Background(), bytes.NewBufferString(str))
	if err != nil {
		return errors.Wrap(err, "unmarshal SecretString")
	}

	return nil
}

func (ss *SecretString) MarshalJSON() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := json.NewEncoder(buf).Encode(ss.wrapped.DecryptedString()); err != nil {
		return nil, errors.Wrap(err, "marshal SecretString")
	}

	return buf.Bytes(), nil
}

func (ss *SecretString) String() string {
	return ss.wrapped.String()
}

// EncryptedString returns sensitive information with encrypted string.
func (ss *SecretString) EncryptedString() string {
	return ss.wrapped.EncryptedString()
}

// DecryptedString returns sensitive information with decrypted string.
func (ss *SecretString) DecryptedString() string {
	return ss.wrapped.DecryptedString()
}
