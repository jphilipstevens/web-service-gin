package dependencies

import (
	"example/web-service-gin/app/cache"
	"example/web-service-gin/app/db"

	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	Cache  cache.Cacher
	DB     *db.Database
	Router *gin.Engine
}
