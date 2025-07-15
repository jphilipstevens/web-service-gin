package dependencies

import (
	"github.com/jphilipstevens/web-service-gin/v2/pkg/appTracer"

	"github.com/gin-gonic/gin"
)

// Dependencies is a minimal container for cross-cutting concerns. The generic
// type parameter allows application code to specify any database client type
// without coupling to this library's database abstraction.
// Dependencies is a minimal container for cross-cutting concerns.
type Dependencies struct {
	Router *gin.Engine
	Tracer appTracer.AppTracer
}
