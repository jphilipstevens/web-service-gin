package clientContext

import (
	"context"
	"time"
)

// ClientContext represents the client context that holds information about the client making the request.
type contextKey string

// ClientContextKey is the key used to store and retrieve the client context from the context.
const (
	ClientContextKey contextKey = "clientContext"
)

// ServiceTransaction represents a transaction within a service.
// It contains information about the service name and span ID.
type ServiceTransaction struct {
	// ServiceName is the name of the service involved in the transaction.
	ServiceName string

	// SpanId is the unique identifier for this specific span of the transaction.
	SpanId string
}

// ClientInfo represents information about the client making the request.
// It contains the client's IP address and User-Agent string.
type ClientInfo struct {
	// IP is the IP address of the client.
	IP string

	// UserAgent is the User-Agent string provided by the client's browser or application.
	UserAgent string
}

// DownstreamCall represents a call made to a downstream service.
// It contains information about the service transaction, response time,
// status code, error (if any), and cache ID.
type DownstreamCall struct {
	// ServiceTransaction contains information about the service and span ID.
	ServiceTransaction

	// ResponseTime is the duration it took for the downstream call to complete.
	ResponseTime time.Duration

	// StatusCode is the HTTP status code returned by the downstream service.
	StatusCode int

	// Error holds any error that occurred during the downstream call.
	Error error

	// CacheId is the identifier for any cache associated with this call.
	CacheId string
}

// RequestInfo represents information about the incoming HTTP request.
type RequestInfo struct {
	// Method is the HTTP method used for the request (e.g., GET, POST, PUT, DELETE).
	Method string

	// Path is the requested URL path.
	Path string
}

// ResponseInfo represents information about the HTTP response.
type ResponseInfo struct {
	// Status is the HTTP status code of the response.
	Status int
}

// DatabaseCall represents a call made to a database.
type DatabaseCall struct {
	// ServiceTransaction contains information about the service and span ID.
	ServiceTransaction

	// Query is the SQL query or operation executed on the database.
	Query string

	// ResponseTime is the duration it took for the database call to complete.
	ResponseTime time.Duration

	// Error holds any error that occurred during the database call.
	Error error
}

// CacheCall represents a call made to a cache service.
type CacheCall struct {
	// ServiceTransaction contains information about the service and span ID.
	ServiceTransaction

	// Action is the type of operation performed on the cache (e.g., GET, SET, DELETE).
	Action string

	// Key is the identifier used to access the cached data.
	Key string

	// Hit indicates whether the cache lookup was successful (true) or not (false).
	Hit bool

	// ResponseTime is the duration it took for the cache call to complete.
	ResponseTime time.Duration

	// Error holds any error that occurred during the cache call.
	Error error
}

// ClientContext represents the context information for a client request.
// It contains information about the service transaction, client, service,
// request, response, downstream calls, database calls, and cache calls.
type ClientContext struct {
	ServiceTransaction
	TraceId      string
	SpanId       string
	Client       ClientInfo
	Request      RequestInfo
	Response     ResponseInfo
	Downstreams  []DownstreamCall
	Database     []DatabaseCall
	Cache        []CacheCall
	ResponseTime time.Duration
}

func GetClientContext(ctx context.Context) *ClientContext {
	return ctx.Value(ClientContextKey).(*ClientContext)
}

// since we are saving the client context as a pointer add any modifications to the client context here and handle multiple go routines safely

func AddResponseTime(ctx context.Context, responseTime time.Duration) {
	currentContext := ctx.Value(ClientContextKey).(*ClientContext)
	currentContext.ResponseTime = responseTime
}

func AddResponseInfo(ctx context.Context, response ResponseInfo) {
	currentContext := ctx.Value(ClientContextKey).(*ClientContext)
	currentContext.Response = response
}

func AddDownstreamCall(ctx context.Context, call DownstreamCall) {
	currentContext := ctx.Value(ClientContextKey).(*ClientContext)
	currentContext.Downstreams = append(currentContext.Downstreams, call)
}

func AddDatabaseCall(ctx context.Context, call DatabaseCall) {
	currentContext := ctx.Value(ClientContextKey).(*ClientContext)
	currentContext.Database = append(currentContext.Database, call)
}

func AddCacheCall(ctx context.Context, call CacheCall) {
	currentContext := ctx.Value(ClientContextKey).(*ClientContext)
	currentContext.Cache = append(currentContext.Cache, call)
}
