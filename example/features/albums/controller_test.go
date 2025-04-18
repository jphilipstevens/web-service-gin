package albums

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jphilipstevens/web-service-gin/app/db"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAlbumService struct {
	mock.Mock
}

func (m *MockAlbumService) GetAlbums(ctx context.Context, params GetAlbumsParams) (*db.Paginated[Album], error) {
	args := m.Called(ctx, params)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	paginated := result.(*db.Paginated[Album])

	return paginated, args.Error(1)
}

func TestGetAlbums(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Successful retrieval", func(t *testing.T) {
		mockService := new(MockAlbumService)
		controller := NewAlbumController(mockService)

		expectedAlbums := db.Paginated[Album]{
			Items: []Album{
				{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
			},
		}

		mockService.On("GetAlbums", mock.Anything, GetAlbumsParams{
			Artist: "John Coltrane",
			Limit:  10,
			Page:   1,
		}).Return(&expectedAlbums, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/albums?artist=John Coltrane", nil)

		controller.GetAlbums(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response db.Paginated[Album]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedAlbums, response)

		mockService.AssertExpectations(t)
	})

	t.Run("Error from service", func(t *testing.T) {
		mockService := new(MockAlbumService)
		controller := NewAlbumController(mockService)

		mockService.On("GetAlbums", mock.Anything, GetAlbumsParams{
			Artist: "Unknown",
			Limit:  10,
			Page:   1,
		}).Return(nil, errors.New("service error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/albums?artist=Unknown", nil)

		controller.GetAlbums(c)

		assert.Equal(t, len(c.Errors), 1)
		error := c.Errors[0]
		assert.Equal(t, error.Err.Error(), "service error")

		mockService.AssertExpectations(t)
	})
}
