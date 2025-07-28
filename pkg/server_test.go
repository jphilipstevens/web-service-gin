package server

import (
	"net"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/jphilipstevens/web-service-gin/v2/pkg/config"
)

func TestMiddlewareSlicesAreImmutable(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := config.Config{AppName: "test", EnableOpenTelemetry: false}
	s := New(cfg)

	m1 := func(c *gin.Context) {}
	s.UseBefore(m1)

	got := s.PreMiddleware()
	got = append(got, func(c *gin.Context) {})

	assert.Equal(t, 1, len(s.PreMiddleware()))
}

func TestBuiltInMiddlewareToggle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := config.Config{AppName: "test", EnableOpenTelemetry: false}
	s := New(cfg)

	assert.Equal(t, 5, len(s.BuiltInMiddleware()))

	cfg.EnableOpenTelemetry = true
	s = New(cfg)
	assert.Equal(t, 6, len(s.BuiltInMiddleware()))
}

func TestGracefulShutdown(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := config.Config{AppName: "test", EnableOpenTelemetry: false}
	s := New(cfg)
	s.RegisterRoutes(func(r *gin.Engine) {
		r.GET("/ping", func(c *gin.Context) {})
	})
	s.ApplyMiddleware()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	assert.NoError(t, err)

	httpSrv := &http.Server{Handler: s.Router()}
	go httpSrv.Serve(ln)

	done := make(chan error)
	go func() { done <- GracefulShutdown(httpSrv, time.Millisecond) }()

	time.Sleep(10 * time.Millisecond)
	p, _ := os.FindProcess(os.Getpid())
	p.Signal(syscall.SIGTERM)

	assert.NoError(t, <-done)
}
