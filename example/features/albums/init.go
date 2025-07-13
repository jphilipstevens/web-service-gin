package albums

import (
	"github.com/jphilipstevens/web-service-gin/app/dependencies"
	"github.com/jphilipstevens/web-service-gin/app/middleware"
)

func Init(deps *dependencies.Dependencies) {
	albumsRepository := NewAlbumRepository(deps.DB)
	albumService := NewAlbumService(deps.Cache, albumsRepository)
	albumController := NewAlbumController(albumService)

	v1 := deps.Router.Group("/v1")
	v1.Use(middleware.AuthMiddleware())
	v1.GET("/albums", middleware.RequireHeader("X-Rate", "1"), albumController.GetAlbums)
	// v1.GET("/albums/:id", getAlbum)
}
