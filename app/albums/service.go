package albums

import (
	"context"
	"encoding/json"
	"example/web-service-gin/app/cache"
	"example/web-service-gin/app/db"
	"fmt"
)

type AlbumService interface {
	GetAlbums(ctx context.Context, artist string) (*db.Paginated[Album], error)
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

func (as *albumService) GetAlbums(ctx context.Context, artist string) (*db.Paginated[Album], error) {
	albumSearchCacheKey := fmt.Sprintf("_albumsArtistFilter:%s", artist)
	cachedAlbums, err := as.cacher.Get(ctx, albumSearchCacheKey)
	if err == nil && cachedAlbums != "" {
		var filteredAlbums db.Paginated[Album]
		marshallError := json.Unmarshal([]byte(cachedAlbums), &filteredAlbums)

		// TODO create better error handling
		return &filteredAlbums, marshallError
	}

	// TODO save request to cache for later fetching
	albums, err := as.albumsRepository.GetAlbums(ctx, artist)
	return albums, err
}
