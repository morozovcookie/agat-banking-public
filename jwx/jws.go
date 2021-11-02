package jwx

import (
	"bytes"
	"context"
	"crypto/rsa"
	"io"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/pkg/errors"
)

// TokenSigner represents a service for signing JWT.
type TokenSigner interface {
	// SignToken signs token.
	SignToken(ctx context.Context, dst io.Writer, src jwt.Token) error
}

// RS512TokenSigner represents a service for signing JWT.
type RS512TokenSigner struct {
	alg jwa.SignatureAlgorithm
	key jwk.Key
}

// NewRS512TokenSigner returns a new RS512TokenSigner instance.
func NewRS512TokenSigner(rsaKey *rsa.PrivateKey) (signer *RS512TokenSigner, err error) {
	signer = &RS512TokenSigner{
		alg: jwa.RS512,
		key: nil,
	}

	if signer.key, err = jwk.New(rsaKey); err != nil {
		return nil, errors.Wrap(err, "init RS512TokenSigner")
	}

	return signer, nil
}

// SignToken signs token.
func (signer *RS512TokenSigner) SignToken(_ context.Context, dst io.Writer, src jwt.Token) error {
	signed, err := jwt.Sign(src, signer.alg, signer.key)
	if err != nil {
		return errors.Wrap(err, "sign token")
	}

	if _, err = io.Copy(dst, bytes.NewBuffer(signed)); err != nil {
		return errors.Wrap(err, "sign token")
	}

	return nil
}
