/*
This file contains utility functions for testing. Mostly used to make mock and stub calls easier.
*/
package testUtils

import (
	"context"
	"database/sql"
	"example/web-service-gin/app/appTracer"
	"example/web-service-gin/app/clientContext"
	"example/web-service-gin/app/db"

	"go.opentelemetry.io/otel/trace"
)

func CreateTestContext() context.Context {
	ctx := context.Background()
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

type dummyAppTracer struct {
}

func (d *dummyAppTracer) CreateSpan(ctx context.Context, serviceName string) (context.Context, trace.Span) {
	return ctx, nil
}

func NewAppTracer() appTracer.AppTracer {
	return &dummyAppTracer{}
}

func NewDatabase(mockedDB *sql.DB) db.Database {
	testDatabase := db.DatabaseImpl{
		Client:    mockedDB,
		AppTracer: NewAppTracer(),
	}
	return &testDatabase
}
