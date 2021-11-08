package jaeger

import (
	"context"
	"io"

	banking "github.com/morozovcookie/agat-banking"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var _ banking.SecretFactory = (*SecretFactory)(nil)

// SecretFactory represents a service initialize SecretString object.
type SecretFactory struct {
	tracer  trace.Tracer
	wrapped banking.SecretFactory
	attrs   []attribute.KeyValue
}

// NewSecretFactory returns a new SecretFactory instance.
func NewSecretFactory(tracer trace.Tracer, factory banking.SecretFactory, attrs ...attribute.KeyValue) *SecretFactory {
	return &SecretFactory{
		tracer:  tracer,
		wrapped: factory,
		attrs:   attrs,
	}
}

// CreateFromEncryptedData creates SecretString object from encrypted data.
func (factory *SecretFactory) CreateFromEncryptedData(ctx context.Context, r io.Reader) (banking.SecretString, error) {
	ctx, span := factory.tracer.Start(ctx, "SecretFactory.CreateFromEncryptedData",
		trace.WithAttributes(factory.attrs...))
	defer span.End()

	secret, err := factory.wrapped.CreateFromEncryptedData(ctx, r)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())

		return nil, err // nolint:wrapcheck
	}

	span.SetStatus(codes.Ok, "")
	span.SetAttributes(attribute.Stringer("secret", secret))

	return secret, nil
}

// CreateFromDecryptedData creates SecretString object from decrypted data.
func (factory *SecretFactory) CreateFromDecryptedData(ctx context.Context, r io.Reader) (banking.SecretString, error) {
	ctx, span := factory.tracer.Start(ctx, "SecretFactory.CreateFromDecryptedData",
		trace.WithAttributes(factory.attrs...))
	defer span.End()

	secret, err := factory.wrapped.CreateFromDecryptedData(ctx, r)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())

		return nil, err // nolint:wrapcheck
	}

	span.SetStatus(codes.Ok, "")
	span.SetAttributes(attribute.Stringer("secret", secret))

	return secret, nil
}
