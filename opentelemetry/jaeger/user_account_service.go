package jaeger

import (
	"context"

	banking "github.com/morozovcookie/agat-banking"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var _ banking.UserAccountService = (*UserAccountService)(nil)

// UserAccountService represents a service for managing UserAccount data.
type UserAccountService struct {
	wrapped banking.UserAccountService

	tracer trace.Tracer
	attrs  []attribute.KeyValue
}

// NewUserAccountService returns a new UserAccountService instance.
func NewUserAccountService(
	svc banking.UserAccountService,
	tracer trace.Tracer,
	attrs ...attribute.KeyValue,
) *UserAccountService {
	return &UserAccountService{
		wrapped: svc,

		tracer: tracer,
		attrs:  attrs,
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
	attrs := append(svc.attrs, attribute.String("email_address", emailAddress))

	ctx, span := svc.tracer.Start(ctx, "UserAccountService.FindUserAccountByEmailAddress",
		trace.WithAttributes(attrs...))
	defer span.End()

	account, err := svc.wrapped.FindUserAccountByEmailAddress(ctx, emailAddress)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())

		return nil, err // nolint:wrapcheck
	}

	span.SetStatus(codes.Ok, "")
	span.SetAttributes(attribute.Stringer("account_id", account.ID),
		attribute.Stringer("user_id", account.User.ID))

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
	attrs := append(svc.attrs, attribute.String("username", userName))

	ctx, span := svc.tracer.Start(ctx, "UserAccountService.FindUserAccountByUserName",
		trace.WithAttributes(attrs...))
	defer span.End()

	account, err := svc.wrapped.FindUserAccountByUserName(ctx, userName)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())

		return nil, err // nolint:wrapcheck
	}

	span.SetStatus(codes.Ok, "")
	span.SetAttributes(attribute.Stringer("account_id", account.ID),
		attribute.Stringer("user_id", account.User.ID))

	return account, nil
}
