package zap

import (
	"context"

	banking "github.com/morozovcookie/agat-banking"
	"go.uber.org/zap"
)

var _ banking.AuthenticationService = (*AuthenticationService)(nil)

// AuthenticationService represents a service for managing user authentication process.
type AuthenticationService struct {
	wrapped banking.AuthenticationService
	logger  *zap.Logger
}

// NewAuthenticationService returns a new AuthenticationService instance.
func NewAuthenticationService(svc banking.AuthenticationService, logger *zap.Logger) *AuthenticationService {
	return &AuthenticationService{
		wrapped: svc,
		logger:  logger.With(zap.String("component", "AuthenticationService")),
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
	accessToken, refreshToken, err := svc.wrapped.AuthenticateUserByEmail(ctx, email, password)

	svc.logger.Debug("authenticate user by email", zap.String("email", email),
		zap.Stringer("password", password), zap.Stringer("access_token", accessToken),
		zap.Stringer("refresh_token", refreshToken), zap.Error(err))

	if err != nil {
		svc.logger.Error("authenticate user by email", zap.String("email", email),
			zap.Stringer("password", password), zap.Error(err))

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
	accessToken, refreshToken, err := svc.wrapped.AuthenticateUserByUsername(ctx, username, password)

	svc.logger.Debug("authenticate user by username", zap.String("username", username),
		zap.Stringer("password", password), zap.Stringer("access_token", accessToken),
		zap.Stringer("refresh_token", refreshToken), zap.Error(err))

	if err != nil {
		svc.logger.Error("authenticate user by username", zap.String("username", username),
			zap.Stringer("password", password), zap.Error(err))

		return nil, nil, err // nolint:wrapcheck
	}

	return accessToken, refreshToken, nil
}
