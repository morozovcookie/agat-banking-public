package v1

import (
	"context"
	stdjson "encoding/json"
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"
	banking "github.com/morozovcookie/agat-banking"
	"github.com/morozovcookie/agat-banking/json"
	"github.com/pkg/errors"
)

const (
	// SignInPathPrefix is the path prefix for handling user sign in request.
	SignInPathPrefix = "/signin"

	// SignOutPathPrefix is the path prefix for handling user sign out request.
	SignOutPathPrefix = "/signout"

	// RefreshTokenPathPrefix is the path prefix for handling refresh user token request.
	RefreshTokenPathPrefix = "/refresh"
)

var _ http.Handler = (*AuthenticationHandler)(nil)

// AuthenticationHandler represents ah HTTP handler for processing authentication requests.
type AuthenticationHandler struct {
	*Handler

	authenticationService banking.AuthenticationService
	secretFactory         banking.SecretFactory
}

// NewAuthenticationHandler returns a new AuthenticationHandler instance.
func NewAuthenticationHandler(authenticationService banking.AuthenticationService) *AuthenticationHandler {
	h := &AuthenticationHandler{
		Handler: NewHandler(),

		authenticationService: authenticationService,
	}

	h.router.Route(BasePathPrefix, func(r chi.Router) {
		r.Post(SignInPathPrefix, h.handleSignIn)
	})

	return h
}

// SignInRequest represents a set of data that should be passed by user for starting authentication process.
type SignInRequest struct {
	// Username is the user account name.
	Username string `json:"username"`

	// Password is the password from user account.
	Password *json.SecretString `json:"password"`

	useEmailAddress bool
}

// EmailAddressRegex is the regular expression that will be used for validating email address.
const EmailAddressRegex = `(?:[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_` + "`" +
	`{|}~-]+)*|"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*")@` +
	`(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\[(?:(?:(2(5[0-5]|[0-4][0-9])|` +
	`1[0-9][0-9]|[1-9]?[0-9]))\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:` +
	`[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\])`

func decodeSignInRequest(_ context.Context, factory banking.SecretFactory, r *http.Request) (*SignInRequest, error) {
	req := &SignInRequest{
		Username:        "",
		Password:        json.NewSecretString(nil, factory),
		useEmailAddress: false,
	}

	if err := stdjson.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(err, "decode SignInRequest")
	}

	reg := regexp.MustCompile(EmailAddressRegex)

	req.useEmailAddress = reg.MatchString(req.Username)

	return req, nil
}

// SignInResponse represents a set of data that will be returned after successfully finished authentication process.
type SignInResponse struct {
	// AccessToken is the short-live JWT for accessing to data API.
	AccessToken *json.SecretString `json:"access_token"`

	// ExpiresIn is the time which after AccessToken will be invalid.
	ExpiresIn int64 `json:"expires_in"`

	// TokenType is the type of token.
	TokenType string `json:"token_type"`

	// RefreshToken is the long-live JWT for refreshing tokens pair.
	RefreshToken *json.SecretString `json:"refresh_token"`
}

func (h *AuthenticationHandler) handleSignIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req, err := decodeSignInRequest(ctx, h.secretFactory, r)
	if err != nil {
		badRequestError(ctx, w)

		return
	}

	accessToken, refreshToken, err := authorize(ctx, h.authenticationService, req)
	if err != nil {
		return
	}

	resp := &SignInResponse{
		AccessToken:  json.NewSecretString(accessToken.SecretString(), h.secretFactory),
		ExpiresIn:    accessToken.Expiration().Sub(accessToken.IssuedAt()).Milliseconds(),
		TokenType:    "Bearer",
		RefreshToken: json.NewSecretString(refreshToken.SecretString(), h.secretFactory),
	}

	encodeResponse(ctx, w, http.StatusOK, resp)

	putRefreshTokenIntoCookie(ctx, w, r, refreshToken)
}

func authorize(
	ctx context.Context,
	svc banking.AuthenticationService,
	req *SignInRequest,
) (
	banking.Token,
	banking.Token,
	error,
) {
	if req.useEmailAddress {
		// return svc.AuthenticateUserByEmail(ctx, req.Username, req.Password)
		return svc.AuthenticateUserByEmail(ctx, req.Username, nil)
	}

	// return svc.AuthenticateUserByUsername(ctx, req.Username, req.Password)
	return svc.AuthenticateUserByUsername(ctx, req.Username, nil)
}

func putRefreshTokenIntoCookie(_ context.Context, w http.ResponseWriter, r *http.Request, refreshToken banking.Token) {
	cookie := &http.Cookie{
		Name:       "refresh_token",
		Value:      refreshToken.SecretString().DecryptedString(),
		Path:       "/",
		Domain:     r.URL.Host,
		Expires:    refreshToken.Expiration(),
		RawExpires: "",
		MaxAge:     int(refreshToken.Expiration().Sub(refreshToken.IssuedAt()).Seconds()),
		Secure:     true,
		HttpOnly:   true,
		SameSite:   http.SameSiteStrictMode,
		Raw:        "",
		Unparsed:   nil,
	}

	http.SetCookie(w, cookie)
}
