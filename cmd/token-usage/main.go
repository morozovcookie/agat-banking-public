package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"log"

	banking "github.com/morozovcookie/agat-banking"
	"github.com/morozovcookie/agat-banking/aes"
	"github.com/morozovcookie/agat-banking/aes/nonce"
	"github.com/morozovcookie/agat-banking/jwx"
	"github.com/morozovcookie/agat-banking/nanoid"
	"github.com/morozovcookie/agat-banking/time"
)

func main() {
	var (
		identifierGenerator = nanoid.NewIdentifierGenerator()
		factoryConfig       = &TokenFactoryConfiguration{
			accessTokenSigner:   nil,
			refreshTokenSigner:  nil,
			secretFactory:       nil,
			identifierGenerator: identifierGenerator,
			timer:               time.NewUTCTimer(),
		}
		err error
	)

	accessTokenPrivateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatalln(err)
	}

	factoryConfig.accessTokenSigner, err = jwx.NewRS512TokenSigner(accessTokenPrivateKey)
	if err != nil {
		log.Fatalln(err)
	}

	refreshTokenPrivateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatalln(err)
	}

	factoryConfig.refreshTokenSigner, err = jwx.NewRS512TokenSigner(refreshTokenPrivateKey)
	if err != nil {
		log.Fatalln(err)
	}

	factoryConfig.secretFactory, err = aes.NewSecretFactory(nonce.NewRandomNonceGenerator(),
		bytes.NewBufferString("a7681ff138d941377c55aefb4ab667b833a823e582c91317f5b5e33c09e6891e"))
	if err != nil {
		log.Fatalln(err)
	}

	var (
		factory = jwx.NewTokenFactory(factoryConfig)

		ctx     = context.Background()
		account = new(banking.UserAccount)
	)

	if account.ID, err = identifierGenerator.GenerateIdentifier(ctx); err != nil {
		log.Fatalln(err)
	}

	createAccessToken(ctx, factory, account)

	createRefreshToken(ctx, factory, account)
}

func createAccessToken(ctx context.Context, factory banking.TokenFactory, account *banking.UserAccount) {
	accessToken, err := factory.CreateAccessToken(ctx, account)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(accessToken.Value().DecryptedString())
}

func createRefreshToken(ctx context.Context, factory banking.TokenFactory, account *banking.UserAccount) {
	refreshToken, err := factory.CreateRefreshToken(ctx, account)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(refreshToken.Value().DecryptedString())
}
