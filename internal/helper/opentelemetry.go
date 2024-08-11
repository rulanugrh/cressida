package helper

import (
	"context"
	"fmt"
	"time"

	"github.com/rulanugrh/cressida/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
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

func NewOpenTelemetry() IOpenTelemetry {
	return &opentelemetry{otl: otel.Tracer("github.com/rulanugrh/cressida")}
}

func InitTracer() (*trace.TracerProvider, error) {
	// initialize response for headers
	headers := map[string]string {
		"Content-Type": "application/json",
	}

	cfg := config.GetConfig()

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(cfg.Opentelemetry.Endpoint),
			otlptracehttp.WithHeaders(headers),
			otlptracehttp.WithInsecure(),
		),
	)

	println(cfg.Opentelemetry.Endpoint)
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
		trace.WithSampler(trace.ParentBased(trace.TraceIDRatioBased(0.5))),
	)

	// set trace provider and text map propagator
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tracerProvider, nil
}

func (o *opentelemetry) StartTracer(ctx context.Context, name string) tr.Span {
	_, span := o.otl.Start(ctx, name)
	return span
}
