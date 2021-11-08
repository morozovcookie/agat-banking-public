package zap

import (
	"context"

	banking "github.com/morozovcookie/agat-banking"
	"go.uber.org/zap"
)

var _ banking.AuthenticationService = (*AuthenticationService)(nil)

// AuthenticationService represents a service for managing user authentication process.
type AuthenticationService struct {
	loggerCreator LoggerCreator
	wrapped       banking.AuthenticationService
}

// NewAuthenticationService returns a new AuthenticationService instance.
func NewAuthenticationService(creator LoggerCreator, svc banking.AuthenticationService) *AuthenticationService {
	return &AuthenticationService{
		loggerCreator: creator,
		wrapped:       svc,
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
	logger := svc.loggerCreator.CreateLogger(ctx, "AuthenticationService",
		"AuthenticateUserByEmail")

	accessToken, refreshToken, err := svc.wrapped.AuthenticateUserByEmail(ctx, email, password)

	logger.Debug("authenticate user by email", zap.String("email", email), zap.Error(err),
		zap.Stringer("password", password), zap.Stringer("access_token", accessToken),
		zap.Stringer("refresh_token", refreshToken))

	if err != nil {
		logger.Error("authenticate user by email", zap.String("email", email), zap.Error(err),
			zap.Stringer("password", password))

		return nil, nil, err // nolint:wrapcheck
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
	logger := svc.loggerCreator.CreateLogger(ctx, "AuthenticationService",
		"AuthenticateUserByUsername")

	accessToken, refreshToken, err := svc.wrapped.AuthenticateUserByUsername(ctx, username, password)

	logger.Debug("authenticate user by username", zap.String("username", username), zap.Error(err),
		zap.Stringer("password", password), zap.Stringer("access_token", accessToken),
		zap.Stringer("refresh_token", refreshToken))

	if err != nil {
		logger.Error("authenticate user by username", zap.String("username", username), zap.Error(err),
			zap.Stringer("password", password))

		return nil, nil, err // nolint:wrapcheck
	}

	return accessToken, refreshToken, nil
}
