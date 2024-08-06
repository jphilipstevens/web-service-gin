package appTracer

import (
	"context"
	"example/web-service-gin/app/clientContext"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func CreateDownstreamSpan(ctx context.Context, serviceName string) (context.Context, trace.Span) {
	serverName := ctx.Value(clientContext.ServiceNameKey).(string)
	tracer := otel.Tracer(serverName)
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
