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

// Database interface defines methods for interacting with the database.
// It provides an abstraction layer for database operations, allowing for
// easier testing and potential swapping of database implementations.

type Database interface {
	ExecContext(serviceName string, ctx context.Context, query string, args ...any) (*sql.Result, error)
	QueryContext(serviceName string, ctx context.Context, query string, args ...any) (*sql.Rows, error)
	GetClient() *sql.DB
	Close()
}

// DatabaseImpl implements the Database interface and provides methods for
// database operations with tracing capabilities.
//
// Fields:
//   - Client: A pointer to the underlying sql.DB instance for direct database access.
//   - AppTracer: An instance of AppTracer for tracing database operations.
//
// DatabaseImpl encapsulates the database connection and provides methods
// for executing queries and managing the connection while integrating
// with the application's tracing system.

type DatabaseImpl struct {
	Client    *sql.DB
	AppTracer appTracer.AppTracer
}

// NewDatabase creates and initializes a new Database instance.
//
// Parameters:
//   - dbConfig: Configuration for the database connection.
//   - appTracer: An instance of AppTracer for tracing database operations.
//
// Returns:
//   - Database: A new Database instance.
//   - error: An error if the database connection fails.
//
// The function performs the following steps:
// 1. Constructs a data source name (DSN) string from the provided configuration.
// 2. Opens a database connection using the specified driver and DSN.
// 3. Pings the database to verify the connection.
// 4. If successful, returns a new DatabaseImpl instance.
// 5. If any step fails, it returns an error and closes any opened connection.

func NewDatabase(dbConfig config.DatabaseConfig, appTracer appTracer.AppTracer) (Database, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.SSLMode)
	db, err := sql.Open(dbConfig.Driver, dsn)
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

// Close closes the database connection.
//
// This method should be called when the database is no longer needed to release
// any open resources and connections. It's typically used in a defer statement
// after creating a new database instance.
//
// Example usage:
//
//	db, err := NewDatabase(config, tracer)
//	if err != nil {
//	    // handle error
//	}
//	defer db.Close()

func (db *DatabaseImpl) Close() {
	db.Client.Close()
}

// GetClient returns the underlying sql.DB client.
//
// This method provides access to the raw database client, which can be used
// for operations not covered by the Database interface methods.
//
// Returns:
//   - *sql.DB: The underlying database client.
//
// Example usage:
//
//	rawDB := db.GetClient()
//	// Use rawDB for custom operations

// ExecContext executes a SQL query with tracing and returns the result.
//
// This method executes a SQL query or statement that doesn't return rows. It creates
// a new span for tracing, executes the query, and records the database call in the
// client context.
//
// Parameters:
//   - serviceName: The name of the service making the database call.
//   - ctx: The context for the database operation.
//   - query: The SQL query to execute.
//   - args: Optional arguments for the SQL query.
//
// Returns:
//   - *sql.Result: The result of the SQL execution.
//   - error: An error if the execution fails.
//
// Example usage:
//
//	result, err := db.ExecContext("UserService", ctx, "INSERT INTO users (name, email) VALUES (?, ?)", "John Doe", "john@example.com")
//	if err != nil {
//	    // handle error
//	}
//	// Use result for further operations

// QueryContext executes a query with a new span and saves results to ClientContext
//
// This method executes a SQL query that returns rows. It creates a new span for
// tracing, executes the query, and records the database call in the client context.
//
// Parameters:
//   - serviceName: The name of the service making the database call.
//   - ctx: The context for the database operation.
//   - query: The SQL query to execute.
//   - args: Optional arguments for the SQL query.
//
// Returns:
//   - *sql.Rows: The result rows from the query.
//   - error: An error if the query execution fails.
//
// Example usage:
//
//	rows, err := db.QueryContext("UserService", ctx, "SELECT id, name FROM users WHERE age > ?", 18)
//	if err != nil {
//	    // handle error
//	}
//	defer rows.Close()
//	// Process rows

func (db *DatabaseImpl) GetClient() *sql.DB {
	return db.Client
}

// ExecContext executes a SQL query or statement that doesn't return rows
//
// This method executes a SQL query or statement that doesn't return rows. It creates
// a new span for tracing, executes the query, and records the database call in the
// client context.
//
// Parameters:
//   - serviceName: The name of the service making the database call. Used for tracing.
//   - ctx: The context for the database operation.
//   - query: The SQL query to execute.
//   - args: Optional arguments for the SQL query.
//
// Returns:
//   - *sql.Result: The result of the SQL execution.
//   - error: An error if the execution fails.
//
// Example usage:
//
//	result, err := db.ExecContext("UserService", ctx, "INSERT INTO users (name, email) VALUES (?, ?)", "John Doe", "john@example.com")
//	if err != nil {
//	    // handle error
//	}
//	// Use result for further operations

func (db *DatabaseImpl) ExecContext(serviceName string, ctx context.Context, query string, args ...any) (*sql.Result, error) {
	startTime := time.Now()
	spanCtx, span := db.AppTracer.CreateSpan(ctx, serviceName)
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

// QueryContext executes a SQL query that returns rows
//
// This method executes a SQL query that returns rows. It creates a new span for tracing,
// executes the query, and records the database call in the client context.
//
// Parameters:
//   - serviceName: The name of the service making the database call. Used for tracing.
//   - ctx: The context for the database operation.
//   - query: The SQL query to execute.
//   - args: Optional arguments for the SQL query.
//
// Returns:
//   - *sql.Rows: The result set of the SQL query.
//   - error: An error if the execution fails.
//
// Example usage:
//
//	rows, err := db.QueryContext("UserService", ctx, "SELECT id, name, email FROM users WHERE id = ?", userId)
//	if err != nil {
//	    // handle error
//	}
//	defer rows.Close()
//	// Process the rows

func (db *DatabaseImpl) QueryContext(serviceName string, ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	startTime := time.Now()
	spanCtx, span := db.AppTracer.CreateSpan(ctx, serviceName)
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
