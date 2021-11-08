package jaeger

import (
	"context"

	banking "github.com/morozovcookie/agat-banking"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TokenService represents a service for managing token data.
type TokenService struct {
	tracer  trace.Tracer
	wrapped banking.TokenService
	attrs   []attribute.KeyValue
}

// NewTokenService returns a new TokenService instance.
func NewTokenService(tracer trace.Tracer, svc banking.TokenService, attrs ...attribute.KeyValue) *TokenService {
	return &TokenService{
		tracer:  tracer,
		wrapped: svc,
		attrs:   attrs,
	}
}

// StoreToken stores a single Token.
func (svc *TokenService) StoreToken(ctx context.Context, token banking.Token) error {
	attrs := append(svc.attrs, attribute.Stringer("token", token))

	ctx, span := svc.tracer.Start(ctx, "TokenService.StoreToken", trace.WithAttributes(attrs...))
	defer span.End()

	if err := svc.wrapped.StoreToken(ctx, token); err != nil {
		span.SetStatus(codes.Error, err.Error())

		return err // nolint:wrapcheck
	}

	span.SetStatus(codes.Ok, "")

	return nil
}

// ExpireToken expires single Token.
// Return the new Token state after update.
func (svc *TokenService) ExpireToken(ctx context.Context, id banking.ID) (banking.Token, error) {
	attrs := append(svc.attrs, attribute.Stringer("id", id))

	ctx, span := svc.tracer.Start(ctx, "TokenService.ExpireToken", trace.WithAttributes(attrs...))
	defer span.End()

	token, err := svc.wrapped.ExpireToken(ctx, id)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())

		return nil, err // nolint:wrapcheck
	}

	span.SetStatus(codes.Ok, "")
	span.SetAttributes(attribute.Stringer("token", token))

	return token, nil
}

// FindTokenByID returns a single Token.
func (svc *TokenService) FindTokenByID(ctx context.Context, id banking.ID) (banking.Token, error) {
	attrs := append(svc.attrs, attribute.Stringer("id", id))

	ctx, span := svc.tracer.Start(ctx, "TokenService.FindTokenByID", trace.WithAttributes(attrs...))
	defer span.End()

	token, err := svc.wrapped.FindTokenByID(ctx, id)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())

		return nil, err // nolint:wrapcheck
	}

	span.SetStatus(codes.Ok, "")
	span.SetAttributes(attribute.Stringer("token", token))

	return token, nil
}

// RemoveExpiredTokens removes expired tokens.
// Return tokens list after remove.
func (svc *TokenService) RemoveExpiredTokens(ctx context.Context, opts banking.FindOptions) ([]banking.Token, error) {
	attrs := append(svc.attrs, attribute.Int64("offset", int64(opts.Offset())),
		attribute.Int64("limit", int64(opts.Limit())))

	ctx, span := svc.tracer.Start(ctx, "", trace.WithAttributes(attrs...))
	defer span.End()

	tt, err := svc.wrapped.RemoveExpiredTokens(ctx, opts)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())

		return nil, err // nolint:wrapxcheck
	}

	span.SetStatus(codes.Ok, "")

	return tt, nil
}
