package albums

import (
	"context"
	"encoding/json"
	"example/web-service-gin/app/cache"
	"example/web-service-gin/app/db"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock AlbumRepository
type MockAlbumRepository struct {
	mock.Mock
}

func (m *MockAlbumRepository) GetAlbums(ctx context.Context, params GetAlbumsParams) (*db.Paginated[Album], error) {
	artist := params.Artist
	args := m.Called(ctx, artist)
	return args.Get(0).(*db.Paginated[Album]), args.Error(1)
}

func (m *MockAlbumRepository) Insert(ctx context.Context, album Album) error {
	args := m.Called(ctx, album)
	return args.Error(0)
}

func (m *MockAlbumRepository) InsertBatch(ctx context.Context, albums []Album) error {
	args := m.Called(ctx, albums)
	return args.Error(0)
}

// Mock Cacher
type MockCacher struct {
	mock.Mock
}

func (m *MockCacher) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockCacher) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func TestNewAlbumService(t *testing.T) {
	mockCacher := new(MockCacher)
	mockRepo := new(MockAlbumRepository)

	service := NewAlbumService(mockCacher, mockRepo)
	assert.NotNil(t, service)
}

func TestGetAlbumsService(t *testing.T) {
	mockCacher := new(MockCacher)
	mockRepo := new(MockAlbumRepository)

	service := NewAlbumService(mockCacher, mockRepo)

	ctx := context.Background()
	artist := "Test Artist"

	t.Run("Cache hit", func(t *testing.T) {
		expectedAlbums := &db.Paginated[Album]{
			Items: []Album{{ID: "1", Title: "Test Album", Artist: "Test Artist", Price: 9.99}},
			Total: 1,
		}
		cachedData, _ := json.Marshal(expectedAlbums)

		mockCacher.On("Get", ctx, "_albumsArtistFilter:Test Artist").Return(string(cachedData), nil).Once()

		albums, err := service.GetAlbums(ctx, GetAlbumsParams{
			Artist: artist,
			Limit:  10,
			Page:   0,
		})

		assert.NoError(t, err)
		assert.Equal(t, expectedAlbums, albums)
		mockCacher.AssertExpectations(t)
	})

	t.Run("Cache miss", func(t *testing.T) {
		expectedAlbums := &db.Paginated[Album]{
			Items: []Album{{ID: "2", Title: "Another Album", Artist: "Test Artist", Price: 14.99}},
			Total: 1,
		}

		mockCacher.On("Get", ctx, "_albumsArtistFilter:Test Artist").Return("", cache.ErrCacheMiss).Once()
		mockRepo.On("GetAlbums", ctx, artist).Return(expectedAlbums, nil).Once()

		albums, err := service.GetAlbums(ctx, GetAlbumsParams{
			Artist: artist,
			Limit:  10,
			Page:   0,
		})

		assert.NoError(t, err)
		assert.Equal(t, expectedAlbums, albums)
		mockCacher.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})
}
