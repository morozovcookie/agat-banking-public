package jwx

import (
	"time"

	"github.com/lestrrat-go/jwx/jwt"
	banking "github.com/morozovcookie/agat-banking"
)

// TokenOption represents an option for configure Token instance.
type TokenOption interface {
	apply(t *Token) error
}

type tokenOptionFunc func(t *Token) error

func (fn tokenOptionFunc) apply(t *Token) error {
	return fn(t)
}

// WithID sets up the token unique identifier.
func WithID(jti banking.ID) TokenOption {
	return tokenOptionFunc(func(t *Token) error {
		return t.token.Set(jwt.JwtIDKey, jti.String())
	})
}

// WithIssuedAt sets up the token issue time.
func WithIssuedAt(iat time.Time) TokenOption {
	return tokenOptionFunc(func(t *Token) error {
		return t.token.Set(jwt.IssuedAtKey, iat)
	})
}

// WithExpiration sets up the token expiration time.
func WithExpiration(exp time.Time) TokenOption {
	return tokenOptionFunc(func(t *Token) error {
		return t.token.Set(jwt.ExpirationKey, exp)
	})
}

// WithSubject sets up the token subject(user).
func WithSubject(sub *banking.UserAccount) TokenOption {
	return tokenOptionFunc(func(t *Token) error {
		t.subject = sub

		return t.token.Set(jwt.SubjectKey, sub.ID.String())
	})
}

// WithNotBefore sets up the time before which the token is not valid.
func WithNotBefore(nbf time.Time) TokenOption {
	return tokenOptionFunc(func(t *Token) error {
		return t.token.Set(jwt.NotBeforeKey, nbf)
	})
}
