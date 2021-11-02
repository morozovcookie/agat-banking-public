package jaeger

import (
	"context"

	banking "github.com/morozovcookie/agat-banking"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var _ banking.AuthenticationService = (*AuthenticationService)(nil)

// AuthenticationService represents a service for managing user authentication process.
type AuthenticationService struct {
	wrapped banking.AuthenticationService

	tracer trace.Tracer
	attrs  []attribute.KeyValue
}

// NewAuthenticationService returns a new AuthenticationService instance.
func NewAuthenticationService(
	svc banking.AuthenticationService,
	tracer trace.Tracer,
	attrs ...attribute.KeyValue,
) *AuthenticationService {
	return &AuthenticationService{
		wrapped: svc,

		tracer: tracer,
		attrs:  attrs,
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
	attrs := append(svc.attrs, attribute.String("email", email))

	ctx, span := svc.tracer.Start(ctx, "AuthenticationService.AuthenticateUserByEmail",
		trace.WithAttributes(attrs...))
	defer span.End()

	accessToken, refreshToken, err := svc.wrapped.AuthenticateUserByEmail(ctx, email, password)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())

		return nil, nil, err // nolint:wrapcheck
	}

	span.SetStatus(codes.Ok, "")

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
	attrs := append(svc.attrs, attribute.String("username", username))

	ctx, span := svc.tracer.Start(ctx, "AuthenticationService.AuthenticateUserByUsername",
		trace.WithAttributes(attrs...))
	defer span.End()

	accessToken, refreshToken, err := svc.wrapped.AuthenticateUserByUsername(ctx, username, password)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())

		return nil, nil, err // nolint:wrapcheck
	}

	span.SetStatus(codes.Ok, "")

	return accessToken, refreshToken, nil
}
