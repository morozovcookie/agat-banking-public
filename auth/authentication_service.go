package auth

import (
	"context"

	banking "github.com/morozovcookie/agat-banking"
	"github.com/pkg/errors"
)

var _ banking.AuthenticationService = (*AuthenticationService)(nil)

// AuthenticationService represents a service for managing user authentication process.
type AuthenticationService struct {
	userAccountService         banking.UserAccountService
	accessTokenBuilderCreator  banking.TokenBuilderCreator
	refreshTokenBuilderCreator banking.TokenBuilderCreator
	tokenService               banking.TokenService
}

// NewAuthenticationService returns a new AuthenticationService instance.
func NewAuthenticationService(
	userAccountService banking.UserAccountService,
	accessTokenBuilderCreator banking.TokenBuilderCreator,
	refreshTokenBuilderCreator banking.TokenBuilderCreator,
	tokenService banking.TokenService,
) *AuthenticationService {
	return &AuthenticationService{
		userAccountService:         userAccountService,
		accessTokenBuilderCreator:  accessTokenBuilderCreator,
		refreshTokenBuilderCreator: refreshTokenBuilderCreator,
		tokenService:               tokenService,
	}
}

// AuthenticateUserByEmail authenticates user by email address and password.
func (svc *AuthenticationService) AuthenticateUserByEmail(
	ctx context.Context,
	email string,
	password banking.SecretString,
) (
	banking.Token,
	banking.Token,
	error,
) {
	account, err := svc.userAccountService.FindUserAccountByEmailAddress(ctx, email)
	if err != nil {
		return nil, nil, errors.Wrap(err, "authenticate user by email")
	}

	accessToken, refreshToken, err := svc.authenticateUser(ctx, account, password)
	if err != nil {
		return nil, nil, errors.Wrap(err, "authenticate user by email")
	}

	return accessToken, refreshToken, nil
}

// AuthenticateUserByUsername authenticates user by username and password.
func (svc *AuthenticationService) AuthenticateUserByUsername(
	ctx context.Context,
	username string,
	password banking.SecretString,
) (
	banking.Token,
	banking.Token,
	error,
) {
	account, err := svc.userAccountService.FindUserAccountByUserName(ctx, username)
	if err != nil {
		return nil, nil, errors.Wrap(err, "authenticate user by username")
	}

	accessToken, refreshToken, err := svc.authenticateUser(ctx, account, password)
	if err != nil {
		return nil, nil, errors.Wrap(err, "authenticate user by username")
	}

	return accessToken, refreshToken, nil
}

func (svc *AuthenticationService) authenticateUser(
	ctx context.Context,
	account *banking.UserAccount,
	password banking.SecretString,
) (
	banking.Token,
	banking.Token,
	error,
) {
	if isSamePassword := account.ComparePassword(password); !isSamePassword {
		return nil, nil, errors.Wrap(banking.ErrIncorrectPassword, "authenticate user")
	}

	accessToken, err := svc.createAccessToken(ctx, account)
	if err != nil {
		return nil, nil, errors.Wrap(err, "authenticate user")
	}

	refreshToken, err := svc.createRefreshToken(ctx, account)
	if err != nil {
		return nil, nil, errors.Wrap(err, "authenticate user")
	}

	if err = svc.tokenService.StoreToken(ctx, refreshToken); err != nil {
		return nil, nil, errors.Wrap(err, "authenticate user")
	}

	return accessToken, refreshToken, nil
}

func (svc *AuthenticationService) createAccessToken(
	ctx context.Context,
	account *banking.UserAccount,
) (
	banking.Token,
	error,
) {
	token, err := svc.accessTokenBuilderCreator.CreateTokenBuilder(ctx).
		WithAccount(account).
		Build(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "create access token")
	}

	return token, nil
}

func (svc *AuthenticationService) createRefreshToken(
	ctx context.Context,
	account *banking.UserAccount,
) (
	banking.Token,
	error,
) {
	token, err := svc.refreshTokenBuilderCreator.CreateTokenBuilder(ctx).
		WithAccount(account).
		Build(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "create access token")
	}

	return token, nil
}
