package albums

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jphilipstevens/web-service-gin/app/dependencies"
	"github.com/jphilipstevens/web-service-gin/testUtils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCache struct {
	Client mock.Mock
}

func (rc *MockCache) Get(serviceName string, ctx context.Context, key string) (string, error) {
	return "", nil
}

func (rc *MockCache) Set(serviceName string, ctx context.Context, key string, value string, expiration time.Duration) error {
	return nil
}

func TestInit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Middleware enforced", func(t *testing.T) {
		// Setup
		client, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer client.Close()

		database := testUtils.NewDatabase(client)

		mockCache := new(MockCache)

		router := gin.Default()

		deps := &dependencies.Dependencies{
			DB:     database,
			Cache:  mockCache,
			Router: router,
		}

		Init(deps)

		routes := router.Routes()
		assert.Len(t, routes, 1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/v1/albums", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

}
