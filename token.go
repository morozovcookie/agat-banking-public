package banking

import (
	"context"
	"fmt"
	"time"
)

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

	// Expiration return the UTC time which after token will be expired.
	Expiration() time.Time

	// Value returns token as SecretString.
	Value() SecretString
}

// TokenFactory represents a service for creating access and refresh tokens.
type TokenFactory interface {
	// CreateAccessToken returns the access token.
	CreateAccessToken(ctx context.Context, account *UserAccount) (Token, error)

	// CreateRefreshToken returns the refresh token.
	CreateRefreshToken(ctx context.Context, account *UserAccount) (Token, error)
}
