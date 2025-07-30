package logutils

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

// WithTrace returns a logrus Entry with traceId and spanId fields attached.
// Always use this for logging inside a request context.
func WithTrace(ctx context.Context) *logrus.Entry {
	sc := trace.SpanContextFromContext(ctx)
	traceId := sc.TraceID().String()
	spanId := sc.SpanID().String()
	return logrus.WithFields(logrus.Fields{
		"traceId": traceId,
		"spanId":  spanId,
	})
}
