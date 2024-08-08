package appTracer

import (
	"context"
	"example/web-service-gin/config"

	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type AppTracer interface {
	CreateDownstreamSpan(ctx context.Context, serviceName string) (context.Context, trace.Span)
}

type appTracerImpl struct {
	serverName string
	tracer     trace.Tracer
}

func initTracer(configFile config.ConfigFile) (trace.Tracer, error) {
	uptrace.ConfigureOpentelemetry(
		uptrace.WithDSN(configFile.Uptrace.DSN),
		uptrace.WithServiceName(configFile.AppName),
		uptrace.WithServiceVersion(configFile.AppVersion),
	)

	return otel.Tracer(configFile.AppName), nil
}

func NewDownstreamSpan(configFile config.ConfigFile) AppTracer {
	tracer, err := initTracer(configFile)
	if err != nil {
		panic(err)
	}
	return &appTracerImpl{
		serverName: configFile.AppName,
		tracer:     tracer,
	}
}

func (d *appTracerImpl) CreateDownstreamSpan(ctx context.Context, serviceName string) (context.Context, trace.Span) {
	_, span := d.tracer.Start(ctx, serviceName)

	return ctx, span
}
