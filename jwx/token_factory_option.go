package jwx

import (
	"time"
)

// TokenFactoryOption represents an option for configure TokenFactory instance.
type TokenFactoryOption interface {
	apply(factory *TokenFactory)
}

type tokenFactoryOptionFunc func(factory *TokenFactory)

func (fn tokenFactoryOptionFunc) apply(factory *TokenFactory) {
	fn(factory)
}

// DefaultAccessTokenLifetime is the default time during access token will be valid.
const DefaultAccessTokenLifetime = time.Minute * 15

// WithAccessTokenLifetime sets up the time during access token will be valid.
func WithAccessTokenLifetime(lifetime time.Duration) TokenFactoryOption {
	return tokenFactoryOptionFunc(func(factory *TokenFactory) {
		factory.accessTokenLifetime = lifetime
	})
}

// DefaultRefreshTokenLifetime is the default time during refresh token will be valid.
const DefaultRefreshTokenLifetime = time.Hour * 24 * 14

// WithRefreshTokenLifetime sets up the time during refresh token will be valid.
func WithRefreshTokenLifetime(lifetime time.Duration) TokenFactoryOption {
	return tokenFactoryOptionFunc(func(factory *TokenFactory) {
		factory.refreshTokenLifetime = lifetime
	})
}
