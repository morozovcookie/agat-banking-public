package v1

import (
	"net/http"
)

const (
	//
	SignInPathPrefix = "/api/v1/signin"

	//
	SignOutPathPrefix = "/api/v1/signout"

	//
	RefreshTokenPathPrefix = "/api/v1/refresh"
)

//
type AuthenticationHandler struct{}

//
func NewAuthenticationHandler() *AuthenticationHandler {
	return &AuthenticationHandler{}
}

func (h *AuthenticationHandler) handleSignIn(w http.ResponseWriter, r *http.Request) {

}

func (h *AuthenticationHandler) handleRefreshToken(w http.ResponseWriter, r *http.Request) {

}

func (h *AuthenticationHandler) handleSignOut(w http.ResponseWriter, r *http.Request) {

}
