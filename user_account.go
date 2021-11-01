package banking

import (
	"context"
	"time"
)

// UserAccount represents a user account information.
type UserAccount struct {
	// ID is the user account unique identifier.
	ID ID

	// UserName is the user login.
	UserName string

	// EmailAddress is the user email address.
	EmailAddress string

	// PasswordHash is the user hashed password.
	PasswordHash string

	// User is the owner of account.
	User *User

	// CreatedAt is the time when user account was created.
	CreatedAt time.Time

	// UpdatedAt is the time when user account was updated.
	UpdateAt time.Time
}

// ComparePassword returns true if passed password value and stored password value are the same.
func (ua *UserAccount) ComparePassword(password string) bool {
	return false
}

// UserAccountService represents a service for managing UserAccount data.
type UserAccountService interface {
	// FindUserAccountByEmailAddress returns UserAccount by UserAccount.EmailAddress.
	FindUserAccountByEmailAddress(ctx context.Context, emailAddress string) (*UserAccount, error)

	// FindUserAccountByUserName returns UserAccount by UserAccount.UserName.
	FindUserAccountByUserName(ctx context.Context, userName string) (*UserAccount, error)
}
