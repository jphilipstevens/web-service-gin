package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"example/web-service-gin/app/appTracer"
	"example/web-service-gin/app/clientContext"
)

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`

	DBName  string `mapstructure:"dbname"`
	SSLMode string `mapstructure:"sslmode"`
}

type Database struct {
	client *sql.DB
}

func ConnectToDB(config DatabaseConfig) (*Database, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)
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
	db.client.Close()
}

func (db *Database) GetClient() *sql.DB {
	return db.client
}

// ExecWithSpan executes a SQL query with tracing and returns the result.
func (db *Database) ExecWithSpan(ctx context.Context, serviceName string, query string, args ...any) (sql.Result, error) {
	startTime := time.Now()
	spanCtx, span := appTracer.CreateDownstreamSpan(ctx, serviceName)
	defer span.End()

	stmt, err := db.client.PrepareContext(spanCtx, query)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(spanCtx, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}

	currentContext := ctx.Value(clientContext.ClientContextKey).(*clientContext.ClientContext)
	newDatabaseCall := clientContext.DatabaseCall{
		ServiceTransaction: clientContext.ServiceTransaction{
			ServiceName: serviceName,
			SpanId:      span.SpanContext().TraceID().String(),
		},
		Query:        query,
		ResponseTime: time.Since(startTime),
		Error:        err,
	}
	currentContext.Database = append(currentContext.Database, newDatabaseCall)
	_ = context.WithValue(ctx, clientContext.ClientContextKey, currentContext)

	return result, nil
}

// QueryWithSpan executes a query with a new span and saves results to ClientContext
func (db *Database) QueryWithSpan(ctx context.Context, serviceName string, query string, args ...any) (*sql.Rows, error) {
	startTime := time.Now()
	spanCtx, span := appTracer.CreateDownstreamSpan(ctx, serviceName)
	defer span.End()

	rows, err := db.client.QueryContext(spanCtx, query, args...)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	currentContext := ctx.Value(clientContext.ClientContextKey).(*clientContext.ClientContext)
	newDatabaseCall := clientContext.DatabaseCall{
		ServiceTransaction: clientContext.ServiceTransaction{
			ServiceName: serviceName,
			SpanId:      span.SpanContext().TraceID().String(),
		},
		Query:        query,
		ResponseTime: time.Since(startTime),
		Error:        err,
	}
	currentContext.Database = append(currentContext.Database, newDatabaseCall)
	_ = context.WithValue(ctx, clientContext.ClientContextKey, currentContext)

	return rows, nil
}
