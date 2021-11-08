package jaeger

import (
	"context"
	"time"

	banking "github.com/morozovcookie/agat-banking"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var _ banking.Timer = (*Timer)(nil)

// Timer represents a service for retrieving time data.
type Timer struct {
	tracer  trace.Tracer
	wrapped banking.Timer
	attrs   []attribute.KeyValue
}

// NewTimer returns a new Timer instance.
func NewTimer(tracer trace.Tracer, timer banking.Timer, attrs ...attribute.KeyValue) *Timer {
	return &Timer{
		tracer:  tracer,
		wrapped: timer,
		attrs:   attrs,
	}
}

// Time returns time data.
func (timer *Timer) Time(ctx context.Context) (time.Time, error) {
	ctx, span := timer.tracer.Start(ctx, "", trace.WithAttributes(timer.attrs...))
	defer span.End()

	t, err := timer.wrapped.Time(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())

		return time.Time{}, err // nolint:wrapcheck
	}

	span.SetStatus(codes.Ok, "")
	span.SetAttributes(attribute.Stringer("time", t))

	return t, nil
}
