package jwx

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/lestrrat-go/jwx/jwt"
	banking "github.com/morozovcookie/agat-banking"
	"github.com/pkg/errors"
)

var _ banking.Token = (*Token)(nil)

// Token represents an JWT.
type Token struct {
	token     jwt.Token
	tokenType banking.TokenType

	subject *banking.UserAccount
	value   banking.SecretString
}

func NewToken(
	ctx context.Context,
	factory banking.SecretFactory,
	signer TokenSigner,
	tokenType banking.TokenType,
	opts ...TokenOption,
) (
	t *Token,
	err error,
) {
	t = &Token{
		token:     jwt.New(),
		tokenType: tokenType,
	}

	if err = t.token.Set(`type`, tokenType.String()); err != nil {
		return nil, errors.Wrap(err, "init token")
	}

	for _, opt := range opts {
		if err = opt.apply(t); err != nil {
			return nil, errors.Wrap(err, "init token")
		}
	}

	buf := new(bytes.Buffer)

	if err = t.sign(ctx, signer, buf); err != nil {
		return nil, errors.Wrap(err, "init token")
	}

	if t.value, err = t.encrypt(ctx, factory, buf); err != nil {
		return nil, errors.Wrap(err, "init token")
	}

	return t, nil
}

func (t *Token) sign(ctx context.Context, signer TokenSigner, w io.Writer) error {
	if err := signer.SignToken(ctx, w, t.token); err != nil {
		return errors.Wrap(err, "sign token")
	}

	return nil
}

func (t *Token) encrypt(ctx context.Context, factory banking.SecretFactory, r io.Reader) (banking.SecretString, error) {
	ss, err := factory.CreateFromDecryptedData(ctx, r)
	if err != nil {
		return nil, errors.Wrap(err, "encrypt token")
	}

	return ss, nil
}

func (t *Token) String() string {
	return t.value.String()
}

// ID returns the token unique identifier.
func (t *Token) ID() banking.ID {
	return banking.ID(t.token.JwtID())
}

// Type returns the token type.
func (t *Token) Type() banking.TokenType {
	return t.tokenType
}

// Account returns the subject of token.
func (t *Token) Account() *banking.UserAccount {
	return t.subject
}

// IssuedAt returns the UTC time when token was issued.
func (t *Token) IssuedAt() time.Time {
	return t.token.IssuedAt()
}

// Expiration return the duration which after token will be expired.
func (t *Token) Expiration() time.Time {
	return t.token.Expiration()
}

// Value returns token as SecretString.
func (t *Token) Value() banking.SecretString {
	return t.value
}
