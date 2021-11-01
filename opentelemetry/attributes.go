package opentelemetry

import (
	"strings"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func SQLAttributesFromQuery(query string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.DBStatementKey.String(query),
		semconv.DBOperationKey.String(strings.SplitAfterN(strings.TrimSpace(query), " ", 1)[0]),
	}
}
