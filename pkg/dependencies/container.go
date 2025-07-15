package dependencies

import (
	"github.com/jphilipstevens/web-service-gin/pkg/appTracer"
	"github.com/jphilipstevens/web-service-gin/pkg/cache"

	"github.com/gin-gonic/gin"
)

// Dependencies is a minimal container for cross-cutting concerns. The generic
// type parameter allows application code to specify any database client type
// without coupling to this library's database abstraction.
type Dependencies[T any] struct {
	Cache  cache.Cacher
	DB     T
	Router *gin.Engine
	Tracer appTracer.AppTracer
}
