package clientContext

import "time"

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
	TraceId     string
	SpanId      string
	Client      ClientInfo
	Request     RequestInfo
	Response    ResponseInfo
	Downstreams []DownstreamCall
	Database    []DatabaseCall
	Cache       []CacheCall
}
