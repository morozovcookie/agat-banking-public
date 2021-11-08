package main

import (
	"bytes"
	"context"
	stdrand "crypto/rand"
	stdrsa "crypto/rsa"
	"encoding/hex"
	"io"
	"log"

	banking "github.com/morozovcookie/agat-banking"
	"github.com/morozovcookie/agat-banking/aes"
	"github.com/morozovcookie/agat-banking/aes/rand"
	"github.com/morozovcookie/agat-banking/jwx"
	"github.com/morozovcookie/agat-banking/jwx/rsa"
	"github.com/morozovcookie/agat-banking/nanoid"
	"github.com/morozovcookie/agat-banking/time"
)

func main() {
	accessTokenPrivateKey, err := stdrsa.GenerateKey(stdrand.Reader, 4096)
	if err != nil {
		log.Fatalln(err)
	}

	accessTokenSigner, err := rsa.NewRS512TokenSigner(accessTokenPrivateKey)
	if err != nil {
		log.Fatalln(err)
	}

	refreshTokenPrivateKey, err := stdrsa.GenerateKey(stdrand.Reader, 4096)
	if err != nil {
		log.Fatalln(err)
	}

	refreshTokenSigner, err := rsa.NewRS512TokenSigner(refreshTokenPrivateKey)
	if err != nil {
		log.Fatalln(err)
	}

	aesKey := make([]byte, aes.CipherKeyLength)

	if _, err = io.ReadFull(stdrand.Reader, aesKey); err != nil {
		log.Fatalln(err)
	}

	key := new(bytes.Buffer)

	if _, err = hex.NewEncoder(key).Write(aesKey); err != nil {
		log.Fatalln(err)
	}

	secretFactory, err := aes.NewSecretFactory(rand.NewNonceGenerator(), key)
	if err != nil {
		log.Fatalln(err)
	}

	var (
		ctx     = context.Background()
		account = new(banking.UserAccount)
	)

	if account.ID, err = nanoid.NewIdentifierGenerator().GenerateIdentifier(ctx); err != nil {
		log.Fatalln(err)
	}

	createAccessToken(ctx, accessTokenSigner, secretFactory, account)

	createRefreshToken(ctx, refreshTokenSigner, secretFactory, account)
}

func createAccessToken(
	ctx context.Context,
	signer jwx.TokenSigner,
	factory banking.SecretFactory,
	account *banking.UserAccount,
) {
	var (
		generator = nanoid.NewIdentifierGenerator()
		timer     = time.NewUTCTimer()

		opts = []jwx.TokenBuilderOption{
			jwx.WithSigner(signer),
			jwx.WithSecretFactory(factory),
			jwx.WithTimer(timer),
			jwx.WithIdentifierGenerator(generator),
		}
	)

	accessToken, err := jwx.NewAccessTokenBuilderCreator(opts...).
		CreateTokenBuilder(ctx).
		WithAccount(account).
		Build(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(accessToken.SecretString().DecryptedString())
}

func createRefreshToken(
	ctx context.Context,
	signer jwx.TokenSigner,
	factory banking.SecretFactory,
	account *banking.UserAccount,
) {
	var (
		generator = nanoid.NewIdentifierGenerator()
		timer     = time.NewUTCTimer()

		opts = []jwx.TokenBuilderOption{
			jwx.WithSigner(signer),
			jwx.WithSecretFactory(factory),
			jwx.WithTimer(timer),
			jwx.WithIdentifierGenerator(generator),
		}
	)

	refreshToken, err := jwx.NewRefreshTokenBuilderCreator(opts...).
		CreateTokenBuilder(ctx).
		WithAccount(account).
		Build(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(refreshToken.SecretString().DecryptedString())
}
