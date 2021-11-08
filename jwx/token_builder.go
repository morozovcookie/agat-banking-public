package jwx

import (
	"context"
	"time"

	"github.com/lestrrat-go/jwx/jwt"
	banking "github.com/morozovcookie/agat-banking"
	"github.com/pkg/errors"
)

var _ banking.TokenBuilder = (*TokenBuilder)(nil)

// TokenBuilder represents a service for parametrized token building.
type TokenBuilder struct {
	tokenType banking.TokenType
	expiresIn time.Duration
	subject   *banking.UserAccount
	until     time.Time

	token jwt.Token

	identifierGenerator banking.IdentifierGenerator
	secretFactory       banking.SecretFactory
	tokenSigner         TokenSigner
	timer               banking.Timer
}

// NewTokenBuilder returns a new TokenBuilder instance.
func NewTokenBuilder(tokenType banking.TokenType, opts ...TokenBuilderOption) *TokenBuilder {
	builder := &TokenBuilder{
		tokenType: tokenType,
		expiresIn: DefaultTokenExpiresIn,
		subject:   nil,

		token: jwt.New(),

		identifierGenerator: nil,
		secretFactory:       nil,
		tokenSigner:         nil,
		timer:               nil,
	}

	_ = builder.token.Set(`typ`, tokenType.String())

	for _, opt := range opts {
		opt.apply(builder)
	}

	return builder
}

// WithID sets up the token unique identifier.
func (builder *TokenBuilder) WithID(jti banking.ID) banking.TokenBuilder {
	_ = builder.token.Set(jwt.JwtIDKey, jti.String())

	return builder
}

// WithAccount sets up the subject of token.
func (builder *TokenBuilder) WithAccount(sub *banking.UserAccount) banking.TokenBuilder {
	_ = builder.token.Set(jwt.SubjectKey, sub.ID.String())

	builder.subject = sub

	return builder
}

// WithIssuedAt sets up the UTC time when token was issued.
func (builder *TokenBuilder) WithIssuedAt(iat time.Time) banking.TokenBuilder {
	_ = builder.token.Set(jwt.IssuerKey, iat)

	return builder
}

// WithExpiration sets up the UTC time which after token will be expired.
func (builder *TokenBuilder) WithExpiration(exp time.Time) banking.TokenBuilder {
	_ = builder.token.Set(jwt.ExpirationKey, exp)

	builder.until = exp

	return builder
}

// WithNotBefore sets up the time before which the token is not valid.
func (builder *TokenBuilder) WithNotBefore(nbf time.Time) banking.TokenBuilder {
	_ = builder.token.Set(jwt.NotBeforeKey, nbf)

	return builder
}

// WithValidUntil sets up the time which after token will be invalid.
// NOTE: Expiration is the constant parameter, but ValidUntil could be changed (e.g. when user signing out).
func (builder *TokenBuilder) WithValidUntil(until time.Time) banking.TokenBuilder {
	builder.until = until

	return builder
}

// Build creates and returns the Token.
func (builder *TokenBuilder) Build(ctx context.Context) (banking.Token, error) {
	for _, fn := range []func(ctx context.Context) error{
		builder.setUpID,
		builder.setupIssuedAt,
		builder.setUpExpiration,
		builder.setupUpNotBefore,
	} {
		if err := fn(ctx); err != nil {
			return nil, errors.Wrap(err, "build token")
		}
	}

	token := NewToken(builder.tokenType, builder.subject, builder.token)

	if err := token.Sign(ctx, builder.tokenSigner, builder.secretFactory); err != nil {
		return nil, errors.Wrap(err, "build token")
	}

	return token, nil
}

func (builder *TokenBuilder) setUpID(ctx context.Context) error {
	if _, ok := builder.token.Get(jwt.JwtIDKey); ok {
		return nil
	}

	jti, err := builder.identifierGenerator.GenerateIdentifier(ctx)
	if err != nil {
		return errors.Wrap(err, "set up jti")
	}

	return builder.token.Set(jwt.JwtIDKey, jti.String())
}

func (builder *TokenBuilder) setupIssuedAt(ctx context.Context) error {
	if _, ok := builder.token.Get(jwt.IssuedAtKey); ok {
		return nil
	}

	iat, err := builder.timer.Time(ctx)
	if err != nil {
		return errors.Wrap(err, "set up iat")
	}

	if err = builder.token.Set(jwt.IssuedAtKey, iat); err != nil {
		return errors.Wrap(err, "set up iat")
	}

	if err = builder.setUpExpiration(ctx); err != nil {
		return errors.Wrap(err, "set up iat")
	}

	return builder.setupUpNotBefore(ctx)
}

func (builder *TokenBuilder) setUpExpiration(ctx context.Context) error {
	if _, ok := builder.token.Get(jwt.ExpirationKey); ok {
		return nil
	}

	if err := builder.setupIssuedAt(ctx); err != nil {
		return errors.Wrap(err, "set up exp")
	}

	exp := builder.token.IssuedAt().Add(builder.expiresIn)

	if err := builder.token.Set(jwt.ExpirationKey, exp); err != nil {
		return errors.Wrap(err, "set up exp")
	}

	builder.until = exp

	return nil
}

func (builder *TokenBuilder) setupUpNotBefore(ctx context.Context) error {
	if _, ok := builder.token.Get(jwt.NotBeforeKey); ok {
		return nil
	}

	if err := builder.setupIssuedAt(ctx); err != nil {
		return errors.Wrap(err, "set up nbf")
	}

	return builder.token.Set(jwt.NotBeforeKey, builder.token.IssuedAt())
}
