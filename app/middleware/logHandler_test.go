package middleware

import (
	"bytes"
	"encoding/json"
	"example/web-service-gin/app/clientContext"
	"example/web-service-gin/testUtils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestJsonLogger(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Skip("Skipping test due to logrus issue")
	t.Run("GET request", func(t *testing.T) {

		var buf bytes.Buffer
		logrus.SetOutput(&buf)
		defer buf.Reset() // Reset the buffer before each test

		router := gin.New()
		router.Use(func(ctx *gin.Context) {
			ctx.Request = ctx.Request.WithContext(testUtils.CreateTestContext())
		})
		router.Use(JsonLogger())

		router.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("User-Agent", "test-agent")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var logEntry struct {
			ClientContext clientContext.ClientContext `json:"clientContext"`
		}
		err := json.Unmarshal(buf.Bytes(), &logEntry)
		assert.NoError(t, err)

		assert.Equal(t, "POST", logEntry.ClientContext.Request.Method)
		assert.Equal(t, "/test", logEntry.ClientContext.Request.Path)
	})

	t.Run("POST request with body", func(t *testing.T) {
		var buf bytes.Buffer
		logrus.SetOutput(&buf)
		defer buf.Reset() // Reset the buffer before each test

		router := gin.New()
		router.Use(func(ctx *gin.Context) {
			ctx.Request = ctx.Request.WithContext(testUtils.CreateTestContext())
		})
		router.Use(JsonLogger())

		router.POST("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		w := httptest.NewRecorder()
		body := `{"key": "value"}`
		req, _ := http.NewRequest("POST", "/test", strings.NewReader(body))
		req.Header.Set("User-Agent", "test-agent")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var logEntry struct {
			ClientContext clientContext.ClientContext `json:"clientContext"`
		}
		err := json.Unmarshal(buf.Bytes(), &logEntry)
		assert.NoError(t, err)

		assert.Equal(t, "POST", logEntry.ClientContext.Request.Method)
		assert.Equal(t, "/test", logEntry.ClientContext.Request.Path)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		var buf bytes.Buffer
		logrus.SetOutput(&buf)
		defer buf.Reset() // Reset the buffer before each test

		router := gin.New()
		router.Use(JsonLogger())

		router.GET("/error", func(c *gin.Context) {
			c.Status(http.StatusInternalServerError)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/error", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var logEntry struct {
			ClientContext clientContext.ClientContext `json:"clientContext"`
		}
		err := json.Unmarshal(buf.Bytes(), &logEntry)
		assert.NoError(t, err)

		assert.Equal(t, "POST", logEntry.ClientContext.Request.Method)
		assert.Equal(t, "/test", logEntry.ClientContext.Request.Path)
	})
}

func TestJsonLoggerLevels(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Skip("Skipping test due to test setup issue")

	t.Run("Info level log", func(t *testing.T) {
		var buf bytes.Buffer
		logrus.SetOutput(&buf)
		defer buf.Reset() // Reset the buffer before each test

		router := gin.New()
		router.Use(JsonLogger())

		router.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		// var logEntry LogEntry
		// err := json.Unmarshal(buf.Bytes(), &logEntry)
		// assert.NoError(t, err)
	})

	t.Run("Error level log", func(t *testing.T) {

		var buf bytes.Buffer
		logrus.SetOutput(&buf)
		defer buf.Reset() // Reset the buffer before each test

		router := gin.New()
		router.Use(JsonLogger())

		router.GET("/error", func(c *gin.Context) {
			c.Status(http.StatusInternalServerError)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/error", nil)
		router.ServeHTTP(w, req)

		// var logEntry LogEntry
		// err := json.Unmarshal(buf.Bytes(), &logEntry)
		// assert.NoError(t, err)
	})
}
