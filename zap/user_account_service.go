package zap

import (
	"context"

	banking "github.com/morozovcookie/agat-banking"
	"go.uber.org/zap"
)

var _ banking.UserAccountService = (*UserAccountService)(nil)

// UserAccountService represents a service for managing UserAccount data.
type UserAccountService struct {
	loggerCreator LoggerCreator
	wrapped       banking.UserAccountService
}

// NewUserAccountService returns a new UserAccountService instance.
func NewUserAccountService(creator LoggerCreator, svc banking.UserAccountService) *UserAccountService {
	return &UserAccountService{
		loggerCreator: creator,
		wrapped:       svc,
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
	logger := svc.loggerCreator.CreateLogger(ctx, "UserAccountService",
		"FindUserAccountByEmailAddress")

	account, err := svc.wrapped.FindUserAccountByEmailAddress(ctx, emailAddress)

	logger.Debug("find user account by email address", zap.String("email", emailAddress), zap.Error(err),
		zap.Any("account", account))

	if err != nil {
		logger.Error("find user account by email", zap.String("email", emailAddress), zap.Error(err))

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
	logger := svc.loggerCreator.CreateLogger(ctx, "UserAccountService",
		"FindUserAccountByUserName")

	account, err := svc.wrapped.FindUserAccountByUserName(ctx, userName)

	logger.Debug("find user account by username", zap.String("username", userName), zap.Error(err),
		zap.Any("account", account))

	if err != nil {
		logger.Error("find user account by username", zap.String("username", userName), zap.Error(err))

		return nil, err // nolint:wrapcheck
	}

	return account, nil
}
