package jaeger

import (
	"context"

	banking "github.com/morozovcookie/agat-banking"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var _ banking.IdentifierGenerator = (*IdentifierGenerator)(nil)

// IdentifierGenerator represents a service for generating unique identifier.
type IdentifierGenerator struct {
	tracer  trace.Tracer
	wrapped banking.IdentifierGenerator
	attrs   []attribute.KeyValue
}

// NewIdentifierGenerator returns a new IdentifierGenerator instance.
func NewIdentifierGenerator(
	tracer trace.Tracer,
	svc banking.IdentifierGenerator,
	attrs ...attribute.KeyValue,
) *IdentifierGenerator {
	return &IdentifierGenerator{
		tracer:  tracer,
		wrapped: svc,
		attrs:   attrs,
	}
}

// GenerateIdentifier returns unique identifier.
func (gen *IdentifierGenerator) GenerateIdentifier(ctx context.Context) (banking.ID, error) {
	ctx, span := gen.tracer.Start(ctx, "", trace.WithAttributes(gen.attrs...))
	defer span.End()

	id, err := gen.wrapped.GenerateIdentifier(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())

		return banking.EmptyID, err // nolint:wrapcheck
	}

	span.SetStatus(codes.Ok, "")
	span.SetAttributes(attribute.Stringer("id", id))

	return id, nil
}
