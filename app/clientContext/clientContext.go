package clientContext

import (
	"context"
	"time"
)

type contextKey string

const (
	ClientContextKey contextKey = "clientContext"
)

type ServiceTransaction struct {
	ServiceName string
	SpanId      string
}

type ClientInfo struct {
	IP        string
	UserAgent string
}

type DownstreamCall struct {
	ServiceTransaction
	ResponseTime time.Duration
	StatusCode   int
	Error        error
	CacheId      string
}

type RequestInfo struct {
	Method string
	Path   string
}

type ResponseInfo struct {
	Status int
}

type DatabaseCall struct {
	ServiceTransaction
	Query        string
	ResponseTime time.Duration
	Error        error
}

type CacheCall struct {
	ServiceTransaction
	Action       string
	Key          string
	Hit          bool
	ResponseTime time.Duration
	Error        error
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
