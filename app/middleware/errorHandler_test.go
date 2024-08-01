package middleware

import (
	"example/web-service-gin/app/apiErrors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("API Error", func(t *testing.T) {
		router := gin.New()
		router.Use(ErrorHandler)
		router.GET("/test", func(c *gin.Context) {
			apiErr := &apiErrors.APIError{
				Status:  http.StatusBadRequest,
				Message: "Bad Request",
				Code:    "BAD_REQUEST",
			}
			c.Error(apiErr)
			if c.Errors.Last() == nil {
				c.Status(http.StatusOK)
			}
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, `{"error":{"code":"BAD_REQUEST","message":"Bad Request"}}`, w.Body.String())
	})

	t.Run("Generic Error", func(t *testing.T) {
		router := gin.New()
		router.Use(ErrorHandler)
		router.GET("/test", func(c *gin.Context) {
			c.Error(gin.Error{Err: gin.Error{Err: gin.Error{Err: nil}}})
			if c.Errors.Last() == nil {
				c.Status(http.StatusOK)
			}
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, `{"error":"Internal Server Error"}`, w.Body.String())
	})

	t.Run("No Error", func(t *testing.T) {
		router := gin.New()
		router.Use(ErrorHandler)
		router.GET("/test", func(c *gin.Context) {
			if c.Errors.Last() == nil {
				c.Status(http.StatusOK)
			}
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "", w.Body.String())
	})

}
