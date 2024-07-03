package helper

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	tr "go.opentelemetry.io/otel/trace"
)

type IOpenTelemetry interface {
	StartTracer(ctx context.Context, name string) tr.Span
}

type opentelemetry struct {
	otl tr.Tracer
}

func InitTracer() (*trace.TracerProvider, error) {
	// initialize response for headers
	headers := map[string]string {
		"Content-Type": "application/json",
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(":4318"),
			otlptracehttp.WithHeaders(headers),
			otlptracehttp.WithInsecure(),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("something error while create new exporter: %v", err)
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(
			exporter,
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
			trace.WithBatchTimeout(trace.DefaultScheduleDelay * time.Millisecond),
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
		),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("cressida-app"),
			),
		),
	)

	return tracerProvider, nil
}

func NewOpenTelemetry() IOpenTelemetry {
	return &opentelemetry{otl: otel.Tracer("cresside-app")}
}

func (o *opentelemetry) StartTracer(ctx context.Context, name string) tr.Span {
	_, span := o.otl.Start(ctx, name)
	return span
}
