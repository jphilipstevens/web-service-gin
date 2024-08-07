package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func TraceMiddleware(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		tracer := otel.Tracer(serviceName)

		_, span := tracer.Start(ctx, "http-server")
		defer span.End()

		ctx = trace.ContextWithSpan(ctx, span)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
