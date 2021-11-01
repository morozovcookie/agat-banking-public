package zap

import (
	"context"

	banking "github.com/morozovcookie/agat-banking"
	"go.uber.org/zap"
)

var _ banking.UserAccountService = (*UserAccountService)(nil)

// UserAccountService represents a service for managing UserAccount data.
type UserAccountService struct {
	wrapped banking.UserAccountService
	logger  *zap.Logger
}

// NewUserAccountService returns a new UserAccountService instance.
func NewUserAccountService(svc banking.UserAccountService, logger *zap.Logger) *UserAccountService {
	return &UserAccountService{
		wrapped: svc,
		logger:  logger.With(zap.String("component", "UserAccountService")),
	}
}

// FindUserAccountByEmailAddress returns UserAccount by UserAccount.EmailAddress.
func (svc *UserAccountService) FindUserAccountByEmailAddress(
	ctx context.Context,
	emailAddress string,
) (
	*banking.UserAccount,
	error,
) {
	account, err := svc.wrapped.FindUserAccountByEmailAddress(ctx, emailAddress)

	svc.logger.Debug("find user account by email address", zap.String("email", emailAddress),
		zap.Any("account", account), zap.Error(err))

	if err != nil {
		svc.logger.Error("find user account by email", zap.String("email", emailAddress), zap.Error(err))

		return nil, err // nolint:wrapcheck
	}

	return account, nil
}

// FindUserAccountByUserName returns UserAccount by UserAccount.UserName.
func (svc *UserAccountService) FindUserAccountByUserName(
	ctx context.Context,
	userName string,
) (
	*banking.UserAccount,
	error,
) {
	account, err := svc.wrapped.FindUserAccountByUserName(ctx, userName)

	svc.logger.Debug("find user account by username", zap.String("username", userName),
		zap.Any("account", account), zap.Error(err))

	if err != nil {
		svc.logger.Error("find user account by username", zap.String("username", userName), zap.Error(err))

		return nil, err // nolint:wrapcheck
	}

	return account, nil
}
