package datastore

import "context"

// DataStore defines a generic interface for data storage backends.
// Implementations must record a clientContext.DataStoreCall for
// each operation to provide traceability.
type DataStore interface {
	// ExecContext performs a write operation such as INSERT or SET.
	ExecContext(serviceName string, ctx context.Context, operation string, args ...any) (any, error)

	// QueryContext performs a read operation such as SELECT or GET.
	QueryContext(serviceName string, ctx context.Context, operation string, args ...any) (any, error)

	// Close releases any resources held by the store.
	Close()
}
