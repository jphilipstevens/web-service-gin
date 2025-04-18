package albums

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AlbumController interface {
	GetAlbums(c *gin.Context)
}

type albumController struct {
	albumService AlbumService
}

func NewAlbumController(albumService AlbumService) AlbumController {
	return &albumController{albumService}
}

func (ac *albumController) GetAlbums(c *gin.Context) {
	artist := c.Query("artist")
	ctx := c.Request.Context()
	params := GetAlbumsParams{Artist: artist, Limit: 10, Page: 1}
	albums, err := ac.albumService.GetAlbums(ctx, params)
	if err != nil {
		c.Error(err)
		return
	}
	c.IndentedJSON(http.StatusOK, albums)
}
