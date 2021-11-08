package jaeger

import (
	"context"
	"io"

	banking "github.com/morozovcookie/agat-banking"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TokenParser represents a service for parsing token.
type TokenParser struct {
	tracer  trace.Tracer
	wrapped banking.TokenParser
	attrs   []attribute.KeyValue
}

// NewTokenParser returns a new TokenParser instance.
func NewTokenParser(tracer trace.Tracer, parser banking.TokenParser, attrs ...attribute.KeyValue) *TokenParser {
	return &TokenParser{
		tracer:  tracer,
		wrapped: parser,
		attrs:   attrs,
	}
}

// ParseToken parse and returns a Token.
func (parser *TokenParser) ParseToken(ctx context.Context, r io.Reader) (banking.Token, error) {
	ctx, span := parser.tracer.Start(ctx, "TokenParser.ParseToken", trace.WithAttributes(parser.attrs...))
	defer span.End()

	token, err := parser.wrapped.ParseToken(ctx, r)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())

		return nil, err // nolint:wrapcheck
	}

	span.SetStatus(codes.Ok, "")
	span.SetAttributes(attribute.Stringer("token", token))

	return token, nil
}
