/*
AppTracer is a wrapper for the OpenTelemetry tracer. It provides a way to create a span.
*/
package appTracer

import (
	"context"
	"example/web-service-gin/app/version"
	"example/web-service-gin/config"

	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// AppTracer is an interface for creating spans based on the current context.
type AppTracer interface {
	CreateSpan(ctx context.Context, serviceName string) (context.Context, trace.Span)
}

type appTracerImpl struct {
	serverName string
	tracer     trace.Tracer
}

func initTracer(configFile config.ConfigFile) (trace.Tracer, error) {
	uptrace.ConfigureOpentelemetry(
		uptrace.WithDSN(configFile.Uptrace.DSN),
		uptrace.WithServiceName(configFile.AppName),
		uptrace.WithServiceVersion(version.Version),
	)

	return otel.Tracer(configFile.AppName), nil
}

// NewAppTracer creates a new AppTracer.
func NewAppTracer(configFile config.ConfigFile) AppTracer {
	tracer, err := initTracer(configFile)
	if err != nil {
		panic(err)
	}
	return &appTracerImpl{
		serverName: configFile.AppName,
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
