package albums

import (
	"context"
	"example/web-service-gin/app/db"
	"example/web-service-gin/app/dependencies"
	"testing"
	"time"

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

	t.Run("Successful Albums Module Init", func(t *testing.T) {
		// Setup
		client, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer client.Close()

		database := db.Database{Client: client}

		mockCache := new(MockCache)

		router := gin.Default()

		deps := &dependencies.Dependencies{
			DB:     &database,
			Cache:  mockCache,
			Router: router,
		}

		// Execute
		Init(deps)

		// Assert
		routes := router.Routes()
		assert.Len(t, routes, 1, "Should have 1 route")

		route := routes[0]
		assert.Equal(t, "GET", route.Method, "Route method should be GET")
		assert.Equal(t, "/v1/albums", route.Path, "Route path should be /v1/albums")
	})

}
