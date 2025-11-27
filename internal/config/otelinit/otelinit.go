package otelinit

import (
	"context"
	"crypto/tls"
	"log/slog"
	"os"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Config struct {
	OtelEndpoint string
	Headers      map[string]string
	Version      string
	Env          string
	SampleRatio  float64
}

// Setup initializes the OpenTelemetry SDK
func Setup(ctx context.Context, cfg Config) (func(context.Context) error, error) {
	if cfg.SampleRatio <= 0.0 {
		cfg.SampleRatio = 1.0
	}

	creds := credentials.NewTLS(&tls.Config{
		MinVersion: tls.VersionTLS12,
	})
	options := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(cfg.OtelEndpoint),
		otlptracegrpc.WithHeaders(cfg.Headers),
		otlptracegrpc.WithDialOption(grpc.WithTransportCredentials(creds)),
	}
	dialOpts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}

	exporter, err := otlptracegrpc.New(ctx, options...)
	if err != nil {
		return nil, err
	}

	hostname, _ := os.Hostname()
	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithOS(),
		resource.WithContainer(),
		resource.WithAttributes(
			semconv.ServiceName("Critiquefi"),
			semconv.ServiceVersion(cfg.Version),
			semconv.DeploymentEnvironmentName(cfg.Env),
			attribute.String("host.hostname", hostname),
		),
	)
	if err != nil {
		return nil, err
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(
			exporter,
			trace.WithMaxExportBatchSize(512),
			trace.WithBatchTimeout(5*time.Second),
		),
		trace.WithResource(res),
		trace.WithSampler(
			trace.ParentBased(trace.TraceIDRatioBased(cfg.SampleRatio)),
		),
	)

	otel.SetTracerProvider(tracerProvider)

	logExp, err := otlploggrpc.New(ctx,
		otlploggrpc.WithEndpoint(cfg.OtelEndpoint),
		otlploggrpc.WithHeaders(cfg.Headers),
		otlploggrpc.WithDialOption(dialOpts...),
	)
	if err != nil {
		_ = tracerProvider.Shutdown(ctx)
		return nil, err
	}

	loggerProvider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(logExp)),
		log.WithResource(res),
	)

	global.SetLoggerProvider(loggerProvider)

	slog.SetDefault(slog.New(otelslog.NewHandler("Critiquefi")))

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return func(ctx context.Context) error {
		return tracerProvider.Shutdown(ctx)
	}, nil
}
