package middleware

import (
	"context"

	"example/web-service-gin/app/clientContext"
	"strings"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

func ClientContextMiddleware(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Extract trace ID and span ID
		spanContext := trace.SpanContextFromContext(ctx)
		traceId := spanContext.TraceID().String()
		spanId := spanContext.SpanID().String()
		// Get IP address
		ip := c.ClientIP()

		userAgent := c.Request.UserAgent()
		if forwardedFor := c.GetHeader("X-Forwarded-For"); forwardedFor != "" {
			ips := strings.Split(forwardedFor, ",")
			if len(ips) > 0 {
				ip = strings.TrimSpace(ips[0])
			}
		}

		currentContext := clientContext.ClientContext{
			TraceId: traceId,
			SpanId:  spanId,
			Client: clientContext.ClientInfo{
				IP:        ip,
				UserAgent: userAgent,
			},
		}

		ctx = context.WithValue(ctx, clientContext.ServiceNameKey, serviceName)
		ctx = context.WithValue(ctx, clientContext.ClientContextKey, &currentContext)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
