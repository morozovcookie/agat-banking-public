package zap

import (
	"context"

	banking "github.com/morozovcookie/agat-banking"
	"go.uber.org/zap"
)

var _ banking.TokenService = (*TokenService)(nil)

// TokenService represents a service for managing token data.
type TokenService struct {
	loggerCreator LoggerCreator
	wrapped       banking.TokenService
}

// NewTokenService returns a new TokenService instance.
func NewTokenService(creator LoggerCreator, svc banking.TokenService) *TokenService {
	return &TokenService{
		loggerCreator: creator,
		wrapped:       svc,
	}
}

// StoreToken stores a single Token.
func (svc *TokenService) StoreToken(ctx context.Context, token banking.Token) error {
	logger := svc.loggerCreator.CreateLogger(ctx, "TokenService", "StoreToken")

	err := svc.wrapped.StoreToken(ctx, token)

	logger.Debug("store token", zap.Stringer("token", token), zap.Error(err))

	if err != nil {
		logger.Error("store token", zap.Stringer("token", token), zap.Error(err))

		return err // nolint:wrapcheck
	}

	return nil
}

// ExpireToken expires single Token.
// Return the new Token state after update.
func (svc *TokenService) ExpireToken(ctx context.Context, id banking.ID) (banking.Token, error) {
	logger := svc.loggerCreator.CreateLogger(ctx, "TokenService", "ExpireToken")

	token, err := svc.wrapped.ExpireToken(ctx, id)

	logger.Debug("expire token", zap.Stringer("id", id), zap.Stringer("token", token), zap.Error(err))

	if err != nil {
		logger.Error("expire token", zap.Stringer("id", id), zap.Stringer("token", token),
			zap.Error(err))

		return nil, err // nolint:wrapcheck
	}

	return token, nil
}

// FindTokenByID returns a single Token.
func (svc *TokenService) FindTokenByID(ctx context.Context, id banking.ID) (banking.Token, error) {
	logger := svc.loggerCreator.CreateLogger(ctx, "TokenService", "FindTokenByID")

	token, err := svc.wrapped.FindTokenByID(ctx, id)

	logger.Debug("find token by id", zap.Stringer("id", id), zap.Stringer("token", token),
		zap.Error(err))

	if err != nil {
		logger.Error("find token by id", zap.Stringer("id", id), zap.Stringer("token", token),
			zap.Error(err))

		return nil, err // nolint:wrapcheck
	}

	return token, nil
}

// RemoveExpiredTokens removes expired tokens.
// Return tokens list after remove.
func (svc *TokenService) RemoveExpiredTokens(ctx context.Context, opts banking.FindOptions) ([]banking.Token, error) {
	logger := svc.loggerCreator.CreateLogger(ctx, "TokenService", "RemoveExpiredTokens")

	tt, err := svc.wrapped.RemoveExpiredTokens(ctx, opts)

	logger.Debug("remove expired tokens", zap.Any("opts", opts), zap.Any("tokens", tt), zap.Error(err))

	if err != nil {
		logger.Error("remove expired tokens", zap.Any("opts", opts), zap.Any("tokens", tt),
			zap.Error(err))

		return nil, err // nolint:wrapcheck
	}

	return tt, nil
}
