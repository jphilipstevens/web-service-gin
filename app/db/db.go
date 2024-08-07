package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"example/web-service-gin/app/appTracer"
	"example/web-service-gin/app/clientContext"
	"example/web-service-gin/config"
)

type Database interface {
	ExecContext(serviceName string, ctx context.Context, query string, args ...any) (*sql.Result, error)
	QueryContext(serviceName string, ctx context.Context, query string, args ...any) (*sql.Rows, error)
	GetClient() *sql.DB
	Close()
}

type DatabaseImpl struct {
	Client    *sql.DB
	AppTracer appTracer.AppTracer
}

func NewDatabase(dbConfig config.DatabaseConfig, appTracer appTracer.AppTracer) (Database, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.SSLMode)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	dbImpl := &DatabaseImpl{
		Client:    db,
		AppTracer: appTracer,
	}

	return dbImpl, nil
}

func (db *DatabaseImpl) Close() {
	db.Client.Close()
}

func (db *DatabaseImpl) GetClient() *sql.DB {
	return db.Client
}

// ExecContext executes a SQL query with tracing and returns the result.
func (db *DatabaseImpl) ExecContext(serviceName string, ctx context.Context, query string, args ...any) (*sql.Result, error) {
	startTime := time.Now()
	spanCtx, span := db.AppTracer.CreateDownstreamSpan(ctx, serviceName)
	defer span.End()

	result, err := db.Client.ExecContext(spanCtx, query, args...)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("error executing query: %w", err)
	}

	newDatabaseCall := clientContext.DatabaseCall{
		ServiceTransaction: clientContext.ServiceTransaction{
			ServiceName: serviceName,
			SpanId:      span.SpanContext().TraceID().String(),
		},
		Query:        query,
		ResponseTime: time.Since(startTime),
		Error:        err,
	}
	clientContext.AddDatabaseCall(ctx, newDatabaseCall)

	return &result, nil
}

// QueryContext executes a query with a new span and saves results to ClientContext
func (db *DatabaseImpl) QueryContext(serviceName string, ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	startTime := time.Now()
	spanCtx, span := db.AppTracer.CreateDownstreamSpan(ctx, serviceName)
	defer span.End()

	rows, err := db.Client.QueryContext(spanCtx, query, args...)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	newDatabaseCall := clientContext.DatabaseCall{
		ServiceTransaction: clientContext.ServiceTransaction{
			ServiceName: serviceName,
			SpanId:      span.SpanContext().TraceID().String(),
		},
		Query:        query,
		ResponseTime: time.Since(startTime),
		Error:        err,
	}
	clientContext.AddDatabaseCall(ctx, newDatabaseCall)

	return rows, nil
}
