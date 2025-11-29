package telemetry

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/freekieb7/gopenehr/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	metricNoop "go.opentelemetry.io/otel/metric/noop"
	sdkLog "go.opentelemetry.io/otel/sdk/log"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
	sdkResource "go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	traceNoop "go.opentelemetry.io/otel/trace/noop"
)

type Telemetry struct {
	Metrics  *Metrics
	Tracing  *Tracing
	Logger   *Logger
	Shutdown func(context.Context) error
}

func Init(ctx context.Context, settings config.Settings) (*Telemetry, error) {
	res := NewResource(settings.Name, settings.Version)

	if settings.OtelEndpoint == "" {
		return setupNoop(), nil
	}

	return setupEnabled(ctx, res, settings)
}

func setupNoop() *Telemetry {
	mp := metricNoop.NewMeterProvider()
	tp := traceNoop.NewTracerProvider()

	otel.SetMeterProvider(mp)
	otel.SetTracerProvider(tp)

	m := NewMetrics(mp.Meter("noop"))
	t := NewTracing(tp.Tracer("noop"))

	return &Telemetry{
		Metrics:  m,
		Tracing:  t,
		Logger:   &Logger{slog.New(slog.NewTextHandler(os.Stdout, nil))},
		Shutdown: func(context.Context) error { return nil },
	}
}

func setupEnabled(ctx context.Context, res *sdkResource.Resource, settings config.Settings) (*Telemetry, error) {
	tExporterOpts := []otlptracegrpc.Option{otlptracegrpc.WithEndpoint(settings.OtelEndpoint)}
	if settings.OtelInsecure {
		tExporterOpts = append(tExporterOpts, otlptracegrpc.WithInsecure())
	}
	tExporter, err := otlptracegrpc.New(ctx, tExporterOpts...)
	if err != nil {
		return nil, err
	}

	mExporterOpts := []otlpmetricgrpc.Option{otlpmetricgrpc.WithEndpoint(settings.OtelEndpoint)}
	if settings.OtelInsecure {
		mExporterOpts = append(mExporterOpts, otlpmetricgrpc.WithInsecure())
	}
	mExporter, err := otlpmetricgrpc.New(ctx, mExporterOpts...)
	if err != nil {
		return nil, err
	}

	lExporterOpts := []otlploggrpc.Option{otlploggrpc.WithEndpoint(settings.OtelEndpoint)}
	if settings.OtelInsecure {
		lExporterOpts = append(lExporterOpts, otlploggrpc.WithInsecure())
	}
	lExporter, err := otlploggrpc.New(ctx, lExporterOpts...)
	if err != nil {
		return nil, err
	}

	tracerProvider := sdkTrace.NewTracerProvider(
		sdkTrace.WithBatcher(tExporter),
		sdkTrace.WithResource(res),
	)

	reader := sdkMetric.NewPeriodicReader(mExporter, sdkMetric.WithInterval(10*time.Second))
	meterProvider := sdkMetric.NewMeterProvider(
		sdkMetric.WithReader(reader),
		sdkMetric.WithResource(res),
	)

	logProvider := sdkLog.NewLoggerProvider(
		sdkLog.WithProcessor(sdkLog.NewBatchProcessor(lExporter)),
		sdkLog.WithResource(res),
	)

	return &Telemetry{
		Metrics:  NewMetrics(meterProvider.Meter("service")),
		Tracing:  NewTracing(tracerProvider.Tracer("service")),
		Logger:   NewLogger("service", logProvider),
		Shutdown: multiShutdown(tracerProvider.Shutdown, meterProvider.Shutdown, logProvider.Shutdown),
	}, nil
}

func multiShutdown(funcs ...func(context.Context) error) func(context.Context) error {
	return func(ctx context.Context) error {
		var err error
		for _, f := range funcs {
			if e := f(ctx); e != nil {
				err = errors.Join(err, e)
			}
		}
		return err
	}
}
