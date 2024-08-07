package testUtils

import (
	"context"
	"example/web-service-gin/app/clientContext"
)

func CreateTestContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, clientContext.ServiceNameKey, "test")
	currentContext := clientContext.ClientContext{
		TraceId: "test-trace-id",
		SpanId:  "test-span-id",
		Client: clientContext.ClientInfo{
			IP:        "127.0.0.1",
			UserAgent: "test-user-agent",
		},
	}
	ctx = context.WithValue(ctx, clientContext.ClientContextKey, &currentContext)
	return ctx
}
