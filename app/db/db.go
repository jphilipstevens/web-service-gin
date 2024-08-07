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

type Database struct {
	Client *sql.DB
}

func ConnectToDB(dbConfig config.DatabaseConfig) (*Database, error) {
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

	return &Database{db}, nil
}

func (db *Database) Close() {
	db.Client.Close()
}

func (db *Database) GetClient() *sql.DB {
	return db.Client
}

// ExecWithSpan executes a SQL query with tracing and returns the result.
func (db *Database) ExecWithSpan(ctx context.Context, serviceName string, query string, args ...any) (sql.Result, error) {
	startTime := time.Now()
	spanCtx, span := appTracer.CreateDownstreamSpan(ctx, serviceName)
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

	return result, nil
}

// QueryWithSpan executes a query with a new span and saves results to ClientContext
func (db *Database) QueryWithSpan(ctx context.Context, serviceName string, query string, args ...any) (*sql.Rows, error) {
	startTime := time.Now()
	spanCtx, span := appTracer.CreateDownstreamSpan(ctx, serviceName)
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
