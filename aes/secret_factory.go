package aes

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"io"

	banking "github.com/morozovcookie/agat-banking"
	"github.com/pkg/errors"
)

// CipherKeyLength is the required length of cipher key.
const CipherKeyLength = 32

var (
	_ banking.SecretFactory = (*SecretFactory)(nil)

	// ErrWrongCipherTextLength is the error that will be raised when length of ciphertext is less than size of nonce.
	ErrWrongCipherTextLength = errors.New("wrong ciphertext length")

	// ErrWrongKeyLength is the error that will be raised when length of key is less than 32 bytes.
	ErrWrongKeyLength = errors.New("wrong key length")
)

// SecretFactory represents a service initialize SecretString object.
type SecretFactory struct {
	nonceGenerator NonceGenerator
	aead           cipher.AEAD
}

// NewSecretFactory returns a new SecretStringFactory instance.
func NewSecretFactory(nonceGenerator NonceGenerator, key io.Reader) (*SecretFactory, error) {
	buf := bytes.NewBuffer(nil)

	if _, err := buf.ReadFrom(key); err != nil {
		return nil, errors.Wrap(err, "init secret string factory")
	}

	bb, err := hex.DecodeString(buf.String())
	if err != nil {
		return nil, errors.Wrap(err, "init secret string factory")
	}

	if len(bb) < CipherKeyLength {
		return nil, errors.Wrap(ErrWrongKeyLength, "init secret string factory")
	}

	block, err := aes.NewCipher(bb)
	if err != nil {
		return nil, errors.Wrap(err, "init secret string factory")
	}

	f := &SecretFactory{
		nonceGenerator: nonceGenerator,
		aead:           nil,
	}

	if f.aead, err = cipher.NewGCM(block); err != nil {
		return nil, errors.Wrap(err, "init secret string factory")
	}

	return f, nil
}

// CreateFromEncryptedData creates SecretString object from encrypted data.
func (f *SecretFactory) CreateFromEncryptedData(_ context.Context, r io.Reader) (banking.SecretString, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(r); err != nil {
		return nil, errors.Wrap(err, "create from encrypted string")
	}

	bb, err := hex.DecodeString(buf.String())
	if err != nil {
		return nil, errors.Wrap(err, "create from encrypted string")
	}

	nonceSize := f.aead.NonceSize()
	if len(bb) < nonceSize {
		return nil, errors.Wrap(ErrWrongCipherTextLength, "create from encrypted string")
	}

	nonce, ciphertext := bb[:nonceSize], bb[nonceSize:]

	plaintext, err := f.aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.Wrap(err, "create from encrypted string")
	}

	return &SecretString{
		encryptedString: buf.String(),
		decryptedString: bytes.NewBuffer(plaintext).String(),
	}, nil
}

// CreateFromDecryptedData creates SecretString object from decrypted data.
func (f *SecretFactory) CreateFromDecryptedData(ctx context.Context, r io.Reader) (banking.SecretString, error) {
	nonce, err := f.nonceGenerator.GenerateNonce(ctx, f.aead.NonceSize())
	if err != nil {
		return nil, errors.Wrap(err, "create from decrypted data")
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(r); err != nil {
		return nil, errors.Wrap(err, "create from decrypted data")
	}

	encrypted := f.aead.Seal(nonce, nonce, buf.Bytes(), nil)

	return &SecretString{
		encryptedString: hex.EncodeToString(encrypted),
		decryptedString: buf.String(),
	}, nil
}
