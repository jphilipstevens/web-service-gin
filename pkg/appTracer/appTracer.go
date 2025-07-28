/*
AppTracer is a wrapper for the OpenTelemetry tracer. It provides a way to create a span.
*/
package appTracer

import (
	"context"

	"github.com/jphilipstevens/web-service-gin/v2/pkg/config"
	"github.com/jphilipstevens/web-service-gin/v2/pkg/version"

	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// AppTracer is an interface for creating spans based on the current context.
type AppTracer interface {
	CreateSpan(ctx context.Context, serviceName string) (context.Context, trace.Span)
	Shutdown(ctx context.Context) error
}

type appTracerImpl struct {
	serverName string
	tracer     trace.Tracer
}

func initTracer(cfg config.Config) (trace.Tracer, error) {
	uptrace.ConfigureOpentelemetry(
		uptrace.WithDSN(cfg.Uptrace.DSN),
		uptrace.WithServiceName(cfg.AppName),
		uptrace.WithServiceVersion(version.Version),
	)

	return otel.Tracer(cfg.AppName), nil
}

// NewAppTracer creates a new AppTracer.
func NewAppTracer(cfg config.Config) AppTracer {
	tracer, err := initTracer(cfg)
	if err != nil {
		panic(err)
	}
	return &appTracerImpl{
		serverName: cfg.AppName,
		tracer:     tracer,
	}
}

// CreateSpan creates a span based on the parent span in the context.
// The span is created with the service name as the span name.
// The context is returned with the span added. This limits the span as a child of the current context without modifying the current context.
func (d *appTracerImpl) CreateSpan(ctx context.Context, serviceName string) (context.Context, trace.Span) {
	_, span := d.tracer.Start(ctx, serviceName)

	return ctx, span
}

func (d *appTracerImpl) Shutdown(ctx context.Context) error {
	return uptrace.Shutdown(ctx)
}
