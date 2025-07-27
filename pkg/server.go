package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/jphilipstevens/web-service-gin/v2/pkg/config"
	"github.com/jphilipstevens/web-service-gin/v2/pkg/middleware"
)

const gracefulShutdownTimeout = 5 * time.Second

// Server exposes a composable HTTP server instance with optional components.
// Middleware can be registered before or after the built-in set using
// UseBefore and UseAfter.
type Server struct {
	router         *gin.Engine
	config         config.Config
	preMiddleware  []gin.HandlerFunc
	postMiddleware []gin.HandlerFunc
}

// New creates a Server with the provided configuration. Middleware is applied
// when Run is invoked so callers can register their own hooks beforehand.
func New(cfg config.Config) *Server {
	s := &Server{
		router: gin.New(),
		config: cfg,
	}
	return s
}

// UseBefore registers middleware that executes before the built-in stack.
func (s *Server) UseBefore(mw ...gin.HandlerFunc) {
	s.preMiddleware = append(s.preMiddleware, mw...)
}

// UseAfter registers middleware that executes after the built-in stack.
func (s *Server) UseAfter(mw ...gin.HandlerFunc) {
	s.postMiddleware = append(s.postMiddleware, mw...)
}

// RegisterRoutes allows modules to add routes to the server.
func (s *Server) RegisterRoutes(fn func(r *gin.Engine)) {
	fn(s.router)
}

// Run starts the HTTP server, applies middleware in the correct order, and
// waits for a shutdown signal.
func (s *Server) Run() error {
	if len(s.preMiddleware) > 0 {
		s.router.Use(s.preMiddleware...)
	}

	s.router.Use(gin.Recovery())
	s.router.Use(otelgin.Middleware(s.config.AppName))
	s.router.Use(middleware.ClientContextMiddleware())
	s.router.Use(middleware.TraceMiddleware(s.config.AppName))
	s.router.Use(middleware.ErrorHandler)
	s.router.Use(middleware.JsonLogger())

	if len(s.postMiddleware) > 0 {
		s.router.Use(s.postMiddleware...)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port),
		Handler: s.router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatalf("server shutdown failed: %v", err)
	}

	logrus.Info("server exited properly")
	return nil
}
