package dependencies

import (
	"jphilipstevens/web-service-gin/app/appTracer"
	"jphilipstevens/web-service-gin/app/cache"
	"jphilipstevens/web-service-gin/app/db"

	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	Cache  cache.Cacher
	DB     db.Database
	Router *gin.Engine
	Tracer appTracer.AppTracer
}
