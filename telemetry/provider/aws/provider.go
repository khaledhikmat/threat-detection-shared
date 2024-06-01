package awsotelprovider

import (
	"context"
	"time"

	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// Processor is a AWS OpenTelemetry provider.
// Code is from aws-otel-community/sample-apps/go-sample-app/collection/otelprovider.go
func Processor(ctx context.Context, service string) (shutdown func(context.Context) error, err error) {
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(service),
	)

	// Setup trace related
	tp, err := setupTraceProvider(ctx, res)
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(xray.Propagator{}) // Set AWS X-Ray propagator

	exp, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	meterProvider := metric.NewMeterProvider(metric.WithResource(res), metric.WithReader(metric.NewPeriodicReader(exp)), metric.WithView(metric.NewView(
		metric.Instrument{Name: "mp_histogram"},
		metric.Stream{Aggregation: metric.AggregationExplicitBucketHistogram{
			Boundaries: []float64{100, 300, 500},
		}},
	)))

	otel.SetMeterProvider(meterProvider)

	return func(context.Context) (err error) {
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()

		// pushes any last exports to the receiver
		err = meterProvider.Shutdown(ctx)
		if err != nil {
			return err
		}
		err = tp.Shutdown(ctx)
		if err != nil {
			return err
		}
		return nil
	}, nil
}

// setupTraceProvider configures a trace exporter and an AWS X-Ray ID Generator.
func setupTraceProvider(ctx context.Context, res *resource.Resource) (*sdktrace.TracerProvider, error) {
	// INSECURE !! NOT TO BE USED FOR ANYTHING IN PRODUCTION
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())

	if err != nil {
		return nil, err
	}

	idg := xray.NewIDGenerator()

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
		sdktrace.WithIDGenerator(idg),
	)
	return tp, nil
}
