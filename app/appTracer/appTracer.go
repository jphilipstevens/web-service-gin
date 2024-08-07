package appTracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type AppTracer interface {
	CreateDownstreamSpan(ctx context.Context, serviceName string) (context.Context, trace.Span)
}

type appTracerImpl struct {
	serverName string
}

func NewDownstreamSpan(serverName string) AppTracer {
	return &appTracerImpl{
		serverName: serverName,
	}
}

func (d *appTracerImpl) CreateDownstreamSpan(ctx context.Context, serviceName string) (context.Context, trace.Span) {
	tracer := otel.Tracer(d.serverName)
	spanCtx := trace.SpanContextFromContext(ctx)

	// Create a child span with a new span ID
	newSpanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    spanCtx.TraceID(),
		SpanID:     trace.SpanID{},
		TraceFlags: trace.FlagsSampled,
		TraceState: trace.TraceState{},
	})

	ctx, newSpan := tracer.Start(ctx, serviceName,
		trace.WithLinks(trace.Link{SpanContext: newSpanCtx}),
	)
	return ctx, newSpan
}
