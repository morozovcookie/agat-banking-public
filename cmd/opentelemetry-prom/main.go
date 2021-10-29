package main

import (
	"context"
	"log"
	"net/http"

	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func main() {
	var (
		ctx = context.Background()

		config = prometheus.Config{}
		ctrl   = controller.New(
			processor.NewFactory(
				selector.NewWithHistogramDistribution(
					histogram.WithExplicitBoundaries(config.DefaultHistogramBoundaries),
				),
				metric.CumulativeExportKindSelector(),
				processor.WithMemory(true),
			),
			controller.WithResource(
				resource.NewWithAttributes(
					semconv.SchemaURL,
					semconv.ServiceNameKey.String("opentelemetry-prom"),
					semconv.ServiceVersionKey.String("1.0.0"),
					semconv.DeploymentEnvironmentKey.String("production"),
				),
			),
		)
	)

	exporter, err := prometheus.New(config, ctrl)
	if err != nil {
		log.Fatalln(err)
	}

	var (
		provider = exporter.MeterProvider()

		perconaMeter  = provider.Meter("com.banking.percona")
		postgresMeter = provider.Meter("com.banking.postgres")
	)

	perconaQueryErrorsCounter, err := perconaMeter.NewInt64Counter("sql.query.errors")
	if err != nil {
		log.Fatalln(err)
	}

	postgresQueryErrorsCounter, err := postgresMeter.NewInt64Counter("sql.query.errors")
	if err != nil {
		log.Fatalln(err)
	}

	perconaQueryErrorsCounter.Add(ctx, 1, semconv.DBSystemMySQL)
	postgresQueryErrorsCounter.Add(ctx, 1, semconv.DBSystemKey.String("clickhouse"))

	http.HandleFunc("/metric", exporter.ServeHTTP)

	log.Fatalln(http.ListenAndServe(":2222", nil))
}
