package config

import (
	"database/sql"
	"example/web-service-gin/app/cache"

	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	Cache  cache.Cacher
	DB     *sql.DB
	Router *gin.Engine
}
