package main

import (
	banking "github.com/morozovcookie/agat-banking"
	"github.com/morozovcookie/agat-banking/jwx"
)

var _ jwx.TokenFactoryConfiguration = (*TokenFactoryConfiguration)(nil)

type TokenFactoryConfiguration struct {
	accessTokenSigner   jwx.TokenSigner
	refreshTokenSigner  jwx.TokenSigner
	secretFactory       banking.SecretFactory
	identifierGenerator banking.IdentifierGenerator
	timer               banking.Timer
}

func (t *TokenFactoryConfiguration) AccessTokenSigner() jwx.TokenSigner {
	return t.accessTokenSigner
}

func (t *TokenFactoryConfiguration) RefreshTokenSigner() jwx.TokenSigner {
	return t.refreshTokenSigner
}

func (t *TokenFactoryConfiguration) SecretFactory() banking.SecretFactory {
	return t.secretFactory
}

func (t *TokenFactoryConfiguration) IdentifierGenerator() banking.IdentifierGenerator {
	return t.identifierGenerator
}

func (t *TokenFactoryConfiguration) Timer() banking.Timer {
	return t.timer
}
