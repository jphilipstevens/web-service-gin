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

// GetAlbums lists albums.
// @Summary List albums
// @Description Retrieves albums. Requires Authorization header and X-Rate header.
// @Tags albums
// @Produce json
// @Param artist query string false "Artist name"
// @Success 200 {array} Album
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Router /v1/albums [get]
// @Security ApiKeyAuth
// @Middleware AuthMiddleware
// @Middleware RequireHeader
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
