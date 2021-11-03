package jwx

import (
	"context"
	"io"

	"github.com/lestrrat-go/jwx/jwt"
)

// TokenSigner represents a service for signing JWT.
type TokenSigner interface {
	// SignToken signs token.
	SignToken(ctx context.Context, dst io.Writer, src jwt.Token) error
}
