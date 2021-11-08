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
	tokenType banking.TokenType
	account   *banking.UserAccount
	until     time.Time

	value jwt.Token

	secret banking.SecretString
}

func NewToken(tokenType banking.TokenType, account *banking.UserAccount, value jwt.Token) *Token {
	return &Token{
		tokenType: tokenType,
		account:   account,

		value: value,

		secret: nil,
	}
}

func (t *Token) Sign(ctx context.Context, signer TokenSigner, factory banking.SecretFactory) (err error) {
	buf := new(bytes.Buffer)

	if err = sign(ctx, signer, buf, t.value); err != nil {
		return errors.Wrap(err, "sign token")
	}

	if t.secret, err = encrypt(ctx, factory, buf); err != nil {
		return errors.Wrap(err, "sign token")
	}

	return nil
}

func sign(ctx context.Context, signer TokenSigner, dst io.Writer, src jwt.Token) error {
	return signer.SignToken(ctx, dst, src)
}

func encrypt(ctx context.Context, factory banking.SecretFactory, src io.Reader) (banking.SecretString, error) {
	return factory.CreateFromDecryptedData(ctx, src)
}

func (t *Token) String() string {
	return t.secret.String()
}

// ID returns the token unique identifier.
func (t *Token) ID() banking.ID {
	return banking.ID(t.value.JwtID())
}

// Type returns the token type.
func (t *Token) Type() banking.TokenType {
	return t.tokenType
}

// Account returns the subject of token.
func (t *Token) Account() *banking.UserAccount {
	return t.account
}

// IssuedAt returns the UTC time when token was issued.
func (t *Token) IssuedAt() time.Time {
	return t.value.IssuedAt()
}

// Expiration returns the UTC time which after token will be expired.
func (t *Token) Expiration() time.Time {
	return t.value.Expiration()
}

// NotBefore returns the time before which the token is not valid.
func (t *Token) NotBefore() time.Time {
	return t.value.NotBefore()
}

// Until returns the time which after token will be invalid.
// NOTE: Expiration is the constant parameter, but ValidUntil could be changed (e.g. when user signing out).
func (t *Token) Until() time.Time {
	return t.until
}

// SecretString returns token as SecretString.
func (t *Token) SecretString() banking.SecretString {
	return t.secret
}
