package jwx

import (
	"context"

	banking "github.com/morozovcookie/agat-banking"
)

var _ banking.TokenBuilderCreator = (*TokenBuilderCreator)(nil)

// TokenBuilderCreator represents a service for producing TokenBuilder instance.
type TokenBuilderCreator struct {
	tokenType banking.TokenType
	opts      []TokenBuilderOption
}

// NewAccessTokenBuilderCreator returns a new TokenBuilderCreator instance that produce a new TokeBuilder instance for
// creating access token.
func NewAccessTokenBuilderCreator(opts ...TokenBuilderOption) *TokenBuilderCreator {
	return NewTokenBuilderCreator(banking.TokenTypeAccess,
		append([]TokenBuilderOption{WithExpiresIn(DefaultAccessTokenExpiresIn)}, opts...)...)
}

// NewRefreshTokenBuilderCreator returns a new TokenBuilderCreator instance that produce a new TokeBuilder instance for
// creating refresh token.
func NewRefreshTokenBuilderCreator(opts ...TokenBuilderOption) *TokenBuilderCreator {
	return NewTokenBuilderCreator(banking.TokenTypeRefresh,
		append([]TokenBuilderOption{WithExpiresIn(DefaultRefreshTokenExpiresIn)}, opts...)...)
}

// NewTokenBuilderCreator returns a new TokenBuilderCreator instance.
func NewTokenBuilderCreator(tokenType banking.TokenType, opts ...TokenBuilderOption) *TokenBuilderCreator {
	return &TokenBuilderCreator{
		tokenType: tokenType,
		opts:      opts,
	}
}

// CreateTokenBuilder returns a new TokeBuilder instance.
func (creator *TokenBuilderCreator) CreateTokenBuilder(_ context.Context) banking.TokenBuilder {
	return NewTokenBuilder(creator.tokenType, creator.opts...)
}
