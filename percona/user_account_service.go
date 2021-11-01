package percona

import (
	"context"
	"database/sql"
	stdtime "time"

	"github.com/Masterminds/squirrel"
	banking "github.com/morozovcookie/agat-banking"
	"github.com/morozovcookie/agat-banking/time"
	"github.com/pkg/errors"
)

var _ banking.UserAccountService = (*UserAccountService)(nil)

// UserAccountService represents a service for managing UserAccount data.
type UserAccountService struct {
	preparer Preparer
}

// NewUserAccountService returns a new UserAccountService instance.
func NewUserAccountService(preparer Preparer) *UserAccountService {
	return &UserAccountService{
		preparer: preparer,
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
	pred := squirrel.Eq{
		"email_address": emailAddress,
	}

	account, err := svc.findUserAccount(ctx, pred)
	if err != nil {
		return nil, errors.Wrap(err, "find user account by email address")
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
	pred := squirrel.Eq{
		"username": userName,
	}

	account, err := svc.findUserAccount(ctx, pred)
	if err != nil {
		return nil, errors.Wrap(err, "find user account by username")
	}

	return account, nil
}

func (svc *UserAccountService) findUserAccount(ctx context.Context, pred interface{}) (*banking.UserAccount, error) {
	query, args, err := squirrel.Select("account_id", "username", "email_address", "password_hash",
		"user_id", "created_at", "updated_at").
		From("user_accounts").
		Where(pred).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "find user account")
	}

	stmt, err := svc.preparer.PrepareContext(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "find user account")
	}
	defer func(ctx context.Context, stmt Stmt) {
		if closeErr := stmt.Close(ctx); closeErr != nil {
			err = closeErr
		}
	}(ctx, stmt)

	var (
		account = &banking.UserAccount{
			User: &banking.User{},
		}
		createdAt int64
		updatedAt sql.NullInt64
	)

	err = stmt.QueryRowContext(ctx, args...).Scan(&account.ID, &account.UserName, &account.EmailAddress,
		&account.PasswordHash, &account.User.ID, &createdAt, &updatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "find user account")
	}

	account.CreatedAt = stdtime.Unix(0, time.MillisecondsToNanoseconds(createdAt))

	if updatedAt.Valid {
		account.UpdateAt = stdtime.Unix(0, time.MillisecondsToNanoseconds(updatedAt.Int64))
	}

	return account, nil
}
