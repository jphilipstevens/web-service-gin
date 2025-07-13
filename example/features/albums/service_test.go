package albums

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/jphilipstevens/web-service-gin/app/cache"
	"github.com/jphilipstevens/web-service-gin/app/db"

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
	Client mock.Mock
}

func (m *MockCacher) Get(serviceName string, ctx context.Context, key string) (string, error) {
	args := m.Client.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockCacher) Set(serviceName string, ctx context.Context, key string, value string, expiration time.Duration) error {
	args := m.Client.Called(ctx, key, value, expiration)
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
		}
		cachedData, _ := json.Marshal(expectedAlbums)

		mockCacher.Client.On("Get", ctx, "_albumsArtistFilter:Test Artist").Return(string(cachedData), nil).Once()

		albums, err := service.GetAlbums(ctx, GetAlbumsParams{
			Artist: artist,
			Limit:  10,
			Page:   0,
		})

		assert.NoError(t, err)
		assert.Equal(t, expectedAlbums, albums)
		mockCacher.Client.AssertExpectations(t)
	})

	t.Run("Cache miss", func(t *testing.T) {
		expectedAlbums := &db.Paginated[Album]{
			Items: []Album{{ID: "2", Title: "Another Album", Artist: "Test Artist", Price: 14.99}},
		}

		marshalledAlbums, _ := json.Marshal(expectedAlbums)

		mockCacher.Client.On("Get", ctx, "_albumsArtistFilter:Test Artist").Return("", cache.ErrCacheMiss).Once()
		mockCacher.Client.On("Set", ctx, "_albumsArtistFilter:Test Artist", string(marshalledAlbums), time.Minute*albumsCacheTTLMinutes).Return(nil).Once()
		mockRepo.On("GetAlbums", ctx, artist).Return(expectedAlbums, nil).Once()

		albums, err := service.GetAlbums(ctx, GetAlbumsParams{
			Artist: artist,
			Limit:  10,
			Page:   0,
		})

		assert.NoError(t, err)
		assert.Equal(t, expectedAlbums, albums)
		mockCacher.Client.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository error", func(t *testing.T) {
		expectedErr := errors.New("db error")

		mockCacher.Client.On("Get", ctx, "_albumsArtistFilter:Test Artist").Return("", cache.ErrCacheMiss).Once()
		mockRepo.On("GetAlbums", ctx, artist).Return((*db.Paginated[Album])(nil), expectedErr).Once()

		albums, err := service.GetAlbums(ctx, GetAlbumsParams{
			Artist: artist,
			Limit:  10,
			Page:   0,
		})

		assert.Error(t, err)
		assert.Nil(t, albums)
		mockCacher.Client.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})
}
