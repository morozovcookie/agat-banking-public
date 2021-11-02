package banking

import (
	"context"

	"github.com/pkg/errors"
)

// ErrIncorrectPassword is the error that will be raised during password comparison process when two comparable
// password are not the same.
var ErrIncorrectPassword = errors.New("incorrect password")

// AuthenticationService represents a service for managing user authentication process.
type AuthenticationService interface {
	// AuthenticateUserByEmail authenticates user by email address and password.
	AuthenticateUserByEmail(ctx context.Context, email string, password SecretString) (Token, Token, error)

	// AuthenticateUserByUsername authenticates user by username and password.
	AuthenticateUserByUsername(ctx context.Context, username string, password SecretString) (Token, Token, error)
}
