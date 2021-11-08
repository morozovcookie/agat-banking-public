package banking

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/pkg/errors"
)

// ErrTokenDoesNotExist will be raised when token could not be found.
var ErrTokenDoesNotExist = errors.New("token does not exist")

// TokenType represents an enum which describes a possible types for token.
type TokenType int

const (
	// TokenTypeUnknown is the default TokenType value which represents that type of token could not be defined or was
	// not set.
	TokenTypeUnknown TokenType = iota

	// TokenTypeAccess is the type of token which must be used for resource access.
	TokenTypeAccess

	// TokenTypeRefresh is the type of token which must be used for retrieve new access token.
	TokenTypeRefresh
)

func (tt TokenType) String() string {
	if tt == TokenTypeAccess {
		return "access"
	}

	if tt == TokenTypeRefresh {
		return "refresh"
	}

	return ""
}

// Token represents an JWT.
type Token interface {
	fmt.Stringer

	// ID returns the token unique identifier.
	ID() ID

	// Type returns the token type.
	Type() TokenType

	// Account returns the subject of token.
	Account() *UserAccount

	// IssuedAt returns the UTC time when token was issued.
	IssuedAt() time.Time

	// Expiration returns the UTC time which after token will be expired.
	Expiration() time.Time

	// NotBefore returns the time before which the token is not valid.
	NotBefore() time.Time

	// Until returns the time which after token will be invalid.
	// NOTE: Expiration is the constant parameter, but ValidUntil could be changed (e.g. when user signing out).
	Until() time.Time

	// SecretString returns token as SecretString.
	SecretString() SecretString
}

// TokenBuilder represents a service for parametrized token building.
type TokenBuilder interface {
	// WithID sets up the token unique identifier.
	WithID(jti ID) TokenBuilder

	// WithAccount sets up the subject of token.
	WithAccount(sub *UserAccount) TokenBuilder

	// WithIssuedAt sets up the UTC time when token was issued.
	WithIssuedAt(iat time.Time) TokenBuilder

	// WithExpiration sets up the UTC time which after token will be expired.
	WithExpiration(exp time.Time) TokenBuilder

	// WithNotBefore sets up the time before which the token is not valid.
	WithNotBefore(nbf time.Time) TokenBuilder

	// WithValidUntil sets up the time which after token will be invalid.
	// NOTE: Expiration is the constant parameter, but ValidUntil could be changed (e.g. when user signing out).
	WithValidUntil(until time.Time) TokenBuilder

	// Build creates and returns the Token.
	Build(ctx context.Context) (Token, error)
}

// TokenBuilderCreator represents a service for producing TokenBuilder instance.
type TokenBuilderCreator interface {
	// CreateTokenBuilder returns a new TokeBuilder instance.
	CreateTokenBuilder(ctx context.Context) TokenBuilder
}

// TokenParser represents a service for parsing token.
type TokenParser interface {
	// ParseToken parse and returns a Token.
	ParseToken(ctx context.Context, r io.Reader) (Token, error)
}

// TokenService represents a service for managing token data.
type TokenService interface {
	// StoreToken stores a single Token.
	StoreToken(ctx context.Context, token Token) error

	// ExpireToken expires single Token.
	// Return the new Token state after update.
	ExpireToken(ctx context.Context, id ID) (Token, error)

	// FindTokenByID returns a single Token.
	FindTokenByID(ctx context.Context, id ID) (Token, error)

	// RemoveExpiredTokens removes expired tokens.
	// Return tokens list after remove.
	RemoveExpiredTokens(ctx context.Context, opts FindOptions) ([]Token, error)
}
