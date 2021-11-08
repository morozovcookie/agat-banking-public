package jaeger

import (
	"context"
	"time"

	banking "github.com/morozovcookie/agat-banking"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var _ banking.TokenBuilder = (*TokenBuilder)(nil)

// TokenBuilder represents a service for parametrized token building.
type TokenBuilder struct {
	tracer  trace.Tracer
	wrapped banking.TokenBuilder
	attrs   []attribute.KeyValue
}

// NewTokenBuilder returns a new TokenBuilder instance.
func NewTokenBuilder(tracer trace.Tracer, builder banking.TokenBuilder, attrs ...attribute.KeyValue) *TokenBuilder {
	return &TokenBuilder{
		tracer:  tracer,
		wrapped: builder,
		attrs:   attrs,
	}
}

// WithID sets up the token unique identifier.
func (builder *TokenBuilder) WithID(jti banking.ID) banking.TokenBuilder {
	return builder.wrapped.WithID(jti)
}

// WithAccount sets up the subject of token.
func (builder *TokenBuilder) WithAccount(sub *banking.UserAccount) banking.TokenBuilder {
	return builder.wrapped.WithAccount(sub)
}

// WithIssuedAt sets up the UTC time when token was issued.
func (builder *TokenBuilder) WithIssuedAt(iat time.Time) banking.TokenBuilder {
	return builder.wrapped.WithIssuedAt(iat)
}

// WithExpiration sets up the UTC time which after token will be expired.
func (builder *TokenBuilder) WithExpiration(exp time.Time) banking.TokenBuilder {
	return builder.wrapped.WithExpiration(exp)
}

// WithNotBefore sets up the time before which the token is not valid.
func (builder *TokenBuilder) WithNotBefore(nbf time.Time) banking.TokenBuilder {
	return builder.wrapped.WithNotBefore(nbf)
}

// WithValidUntil sets up the time which after token will be invalid.
// NOTE: Expiration is the constant parameter, but ValidUntil could be changed (e.g. when user signing out).
func (builder *TokenBuilder) WithValidUntil(until time.Time) banking.TokenBuilder {
	return builder.wrapped.WithValidUntil(until)
}

// Build creates and returns the Token.
func (builder *TokenBuilder) Build(ctx context.Context) (banking.Token, error) {
	ctx, span := builder.tracer.Start(ctx, "TokenBuilder.Build", trace.WithAttributes(builder.attrs...))
	defer span.End()

	token, err := builder.wrapped.Build(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())

		return nil, err // nolint:wrapcheck
	}

	span.SetStatus(codes.Ok, "")
	span.SetAttributes(attribute.Stringer("token", token))

	return token, nil
}
