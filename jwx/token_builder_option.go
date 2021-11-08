package jwx

import (
	"time"

	banking "github.com/morozovcookie/agat-banking"
)

// TokenBuilderOption represents an option for configure TokenBuilder object.
type TokenBuilderOption interface {
	apply(builder *TokenBuilder)
}

type tokenBuilderOptionFunc func(builder *TokenBuilder)

func (fn tokenBuilderOptionFunc) apply(builder *TokenBuilder) {
	fn(builder)
}

const (
	// DefaultTokenExpiresIn is the default lifetime for any type of token.
	DefaultTokenExpiresIn = time.Minute * 5

	// DefaultAccessTokenExpiresIn is the default lifetime for access token.
	DefaultAccessTokenExpiresIn = time.Minute * 15

	// DefaultRefreshTokenExpiresIn is the default lifetime for refresh token.
	DefaultRefreshTokenExpiresIn = time.Hour * 24 * 30
)

// WithExpiresIn sets up the time duration which token will be valid since it did issue.
func WithExpiresIn(expiresIn time.Duration) TokenBuilderOption {
	return tokenBuilderOptionFunc(func(builder *TokenBuilder) {
		builder.expiresIn = expiresIn
	})
}

// WithSigner sets up the service for signing token.
func WithSigner(signer TokenSigner) TokenBuilderOption {
	return tokenBuilderOptionFunc(func(builder *TokenBuilder) {
		builder.tokenSigner = signer
	})
}

// WithSecretFactory sets up the service for creating a banking.SecretString from token.
func WithSecretFactory(factory banking.SecretFactory) TokenBuilderOption {
	return tokenBuilderOptionFunc(func(builder *TokenBuilder) {
		builder.secretFactory = factory
	})
}

// WithIdentifierGenerator sets up the service for generate unique identifiers.
func WithIdentifierGenerator(generator banking.IdentifierGenerator) TokenBuilderOption {
	return tokenBuilderOptionFunc(func(builder *TokenBuilder) {
		builder.identifierGenerator = generator
	})
}

// WithTimer sets up the service for retrieving time value.
func WithTimer(timer banking.Timer) TokenBuilderOption {
	return tokenBuilderOptionFunc(func(builder *TokenBuilder) {
		builder.timer = timer
	})
}
