package dependencies

import (
	"github.com/jphilipstevens/web-service-gin/app/appTracer"
	"github.com/jphilipstevens/web-service-gin/app/cache"
	"github.com/jphilipstevens/web-service-gin/app/db"

	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	Cache  cache.Cacher
	DB     db.Database
	Router *gin.Engine
	Tracer appTracer.AppTracer
}
