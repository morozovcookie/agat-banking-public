package jwx

import (
	"context"
	"time"

	banking "github.com/morozovcookie/agat-banking"
	"github.com/pkg/errors"
)

var _ banking.TokenFactory = (*TokenFactory)(nil)

// TokenFactoryConfiguration represents a service for TokenFactory configuration.
type TokenFactoryConfiguration interface {
	// AccessTokenSigner returns a TokenSigner for signing access tokens.
	AccessTokenSigner() TokenSigner

	// RefreshTokenSigner returns a TokenSigner for signing refresh tokens.
	RefreshTokenSigner() TokenSigner

	// SecretFactory returns a banking.SecretFactory service instance.
	SecretFactory() banking.SecretFactory

	// IdentifierGenerator return a banking.IdentifierGenerator service instance.
	IdentifierGenerator() banking.IdentifierGenerator

	// Timer returns a banking.Timer service instance.
	Timer() banking.Timer
}

// TokenFactory represents a service for creating access and refresh tokens.
type TokenFactory struct {
	accessTokenSigner  TokenSigner
	refreshTokenSigner TokenSigner

	secretFactory       banking.SecretFactory
	identifierGenerator banking.IdentifierGenerator
	timer               banking.Timer

	accessTokenLifetime  time.Duration
	refreshTokenLifetime time.Duration
}

// NewTokenFactory returns a new TokenFactory instance.
func NewTokenFactory(config TokenFactoryConfiguration, opts ...TokenFactoryOption) *TokenFactory {
	factory := &TokenFactory{
		accessTokenSigner:  config.AccessTokenSigner(),
		refreshTokenSigner: config.RefreshTokenSigner(),

		secretFactory:       config.SecretFactory(),
		identifierGenerator: config.IdentifierGenerator(),
		timer:               config.Timer(),

		accessTokenLifetime:  DefaultAccessTokenLifetime,
		refreshTokenLifetime: DefaultRefreshTokenLifetime,
	}

	for _, opt := range opts {
		opt.apply(factory)
	}

	return factory
}

// CreateAccessToken returns the access token.
func (factory *TokenFactory) CreateAccessToken(
	ctx context.Context,
	account *banking.UserAccount,
) (
	banking.Token,
	error,
) {
	opts, err := factory.createTokenOptions(ctx, account, factory.accessTokenLifetime)
	if err != nil {
		return nil, errors.Wrap(err, "create access token")
	}

	token, err := NewToken(ctx, factory.secretFactory, factory.accessTokenSigner, banking.TokenTypeAccess, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "create access token")
	}

	return token, nil
}

// CreateRefreshToken returns the refresh token.
func (factory *TokenFactory) CreateRefreshToken(
	ctx context.Context,
	account *banking.UserAccount,
) (
	banking.Token,
	error,
) {
	opts, err := factory.createTokenOptions(ctx, account, factory.refreshTokenLifetime)
	if err != nil {
		return nil, errors.Wrap(err, "create refresh token")
	}

	token, err := NewToken(ctx, factory.secretFactory, factory.refreshTokenSigner, banking.TokenTypeRefresh, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "create refresh token")
	}

	return token, nil
}

func (factory *TokenFactory) createTokenOptions(
	ctx context.Context,
	account *banking.UserAccount,
	lifetime time.Duration,
) (
	[]TokenOption,
	error,
) {
	jti, err := factory.identifierGenerator.GenerateIdentifier(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "create token options")
	}

	iat, err := factory.timer.Time(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "create token options")
	}

	opts := []TokenOption{
		WithID(jti),
		WithIssuedAt(iat),
		WithNotBefore(iat),
		WithExpiration(iat.Add(lifetime)),
		WithSubject(account),
	}

	return opts, nil
}
