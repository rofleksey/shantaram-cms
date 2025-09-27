package telemetry

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"shantaram/pkg/build"
	"shantaram/pkg/config"
	"time"

	"github.com/getsentry/sentry-go"
	sentryotel "github.com/getsentry/sentry-go/otel"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

var ServiceName = "shantaram-api"
var MetricsPort = 8081

type Telemetry struct {
	TracerProvider oteltrace.TracerProvider
	MeterProvider  otelmetric.MeterProvider
	Tracer         oteltrace.Tracer
	Meter          otelmetric.Meter
	Shutdown       func(context.Context) error
}

func Init(cfg *config.Config) (*Telemetry, error) {
	if cfg.Telemetry.OTLPEndpoint == "" {
		noopTracer := noop.NewTracerProvider()
		noopMeter := sdkmetric.NewMeterProvider()

		return &Telemetry{
			TracerProvider: noopTracer,
			MeterProvider:  noopMeter,
			Tracer:         noopTracer.Tracer("noop"),
			Meter:          noopMeter.Meter("noop"),
			Shutdown:       func(context.Context) error { return nil },
		}, nil
	}

	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(ServiceName),
			semconv.ServiceVersionKey.String(build.Tag),
			semconv.DeploymentEnvironmentKey.String("production"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	var shutdownFuncs []func(context.Context) error

	tracerProvider, err := initTracerProvider(ctx, res, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize tracer provider: %w", err)
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)

	meterProvider, err := initMeterProvider(ctx, res, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize meter provider: %w", err)
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)

	otel.SetTracerProvider(tracerProvider)
	otel.SetMeterProvider(meterProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		sentryotel.NewSentryPropagator(),
	))

	tracer := tracerProvider.Tracer(ServiceName)
	meter := meterProvider.Meter(ServiceName)

	shutdown := func(ctx context.Context) error {
		var errs []error
		for _, fn := range shutdownFuncs {
			if err := fn(ctx); err != nil {
				errs = append(errs, err)
			}
		}
		if len(errs) > 0 {
			return fmt.Errorf("shutdown errors: %v", errs)
		}
		return nil
	}

	slog.InfoContext(ctx, "Telemetry initialized")

	return &Telemetry{
		TracerProvider: tracerProvider,
		MeterProvider:  meterProvider,
		Tracer:         tracer,
		Meter:          meter,
		Shutdown:       shutdown,
	}, nil
}

func initTracerProvider(ctx context.Context, res *resource.Resource, cfg *config.Config) (*sdktrace.TracerProvider, error) {
	traceExporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(cfg.Telemetry.OTLPEndpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(sentryotel.NewSentrySpanProcessor()),
	), nil
}

func initMeterProvider(ctx context.Context, res *resource.Resource, cfg *config.Config) (*sdkmetric.MeterProvider, error) {
	promExporter, err := prometheus.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create Prometheus exporter: %w", err)
	}

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		server := &http.Server{
			Addr:         fmt.Sprintf(":%d", MetricsPort),
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		}
		slog.InfoContext(ctx, "Prometheus metrics server started",
			slog.String("address", server.Addr),
		)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.ErrorContext(ctx, "Prometheus metrics server error", slog.Any("error", err))
		}
	}()

	return sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(promExporter),
		sdkmetric.WithResource(res),
	), nil
}

func InitSentry(cfg *config.Config) error {
	if cfg.Sentry.DSN == "" {
		return nil
	}

	return sentry.Init(sentry.ClientOptions{
		Dsn:              cfg.Sentry.DSN,
		Environment:      "production",
		TracesSampleRate: 1.0,
		EnableTracing:    true,
		TracesSampler: func(ctx sentry.SamplingContext) float64 {
			if ctx.Span.Name == "GET /v1/healthz" {
				return 0.0
			}
			return 1.0
		},
		AttachStacktrace: true,
	})
}
