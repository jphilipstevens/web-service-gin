/*
This file contains utility functions for testing. Mostly used to make mock and stub calls easier.
*/
package testUtils

import (
	"context"
	"database/sql"

	"github.com/jphilipstevens/web-service-gin/app/appTracer"
	"github.com/jphilipstevens/web-service-gin/app/clientContext"
	"github.com/jphilipstevens/web-service-gin/app/db"

	"go.opentelemetry.io/otel"
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

// dummyAppTracer now returns a NoOpSpan instead of nil
type dummyAppTracer struct{}

func (d *dummyAppTracer) CreateSpan(ctx context.Context, serviceName string) (context.Context, trace.Span) {
	tracer := otel.GetTracerProvider().Tracer("test")
	return tracer.Start(ctx, serviceName)
}

func (d *dummyAppTracer) Shutdown(ctx context.Context) error {
	return nil
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
