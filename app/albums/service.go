package albums

import (
	"context"
	"encoding/json"
	"example/web-service-gin/app/cache"
	"example/web-service-gin/app/db"
	"fmt"
)

const (
	albumsCacheKeySuffix   = "_albumsArtistFilter"
	albumsCacheTTLMinutes  = 10
	albumsCacheServiceName = "albumsCache"
)

type AlbumService interface {
	GetAlbums(ctx context.Context, params GetAlbumsParams) (*db.Paginated[Album], error)
}

type albumService struct {
	cacher           cache.Cacher
	albumsRepository AlbumRepository
}

func NewAlbumService(cacher cache.Cacher, albumsRepository AlbumRepository) AlbumService {
	return &albumService{
		cacher:           cacher,
		albumsRepository: albumsRepository,
	}
}

func (as *albumService) GetAlbums(ctx context.Context, params GetAlbumsParams) (*db.Paginated[Album], error) {

	albumSearchCacheKey := fmt.Sprintf("%s:%s", albumsCacheKeySuffix, params.Artist)
	cachedAlbums, err := as.cacher.Get(serviceName, ctx, albumSearchCacheKey)
	if err == nil && cachedAlbums != "" {
		var filteredAlbums db.Paginated[Album]
		marshallError := json.Unmarshal([]byte(cachedAlbums), &filteredAlbums)

		// TODO create better error handling
		return &filteredAlbums, marshallError
	}

	// TODO save request to cache for later fetching
	albums, err := as.albumsRepository.GetAlbums(ctx, params)
	return albums, err
}
