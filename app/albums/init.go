package albums

import (
	"example/web-service-gin/config"
)

func Init(deps *config.Dependencies) {
	albumsRepository := NewAlbumRepository(deps.DB)
	albumService := NewAlbumService(deps.Cache, albumsRepository)
	albumController := NewAlbumController(albumService)

	v1 := deps.Router.Group("/v1")
	v1.GET("/albums", albumController.GetAlbums)
	// v1.GET("/albums/:id", getAlbum)
}