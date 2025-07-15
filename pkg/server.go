package app

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

	"github.com/jphilipstevens/web-service-gin/pkg/appTracer"
	"github.com/jphilipstevens/web-service-gin/pkg/cache"
	"github.com/jphilipstevens/web-service-gin/pkg/config"
	"github.com/jphilipstevens/web-service-gin/pkg/dependencies"
	"github.com/jphilipstevens/web-service-gin/pkg/middleware"
)

const gracefulShutdownTimeout = 5 * time.Second

// RouterFunc registers routes using the provided dependency container. The
// generic type parameter represents the application's database client.
type RouterFunc[T any] func(deps *dependencies.Dependencies[T])

// Server exposes a composable HTTP server instance with optional components.
// The generic type parameter indicates the application's database client type.
type Server[T any] struct {
	config config.ConfigFile
	deps   *dependencies.Dependencies[T]
}

// Config returns the loaded application configuration.
func (s *Server[T]) Config() config.ConfigFile {
	return s.config
}

// Dependencies returns the dependency container for advanced customization.
func (s *Server[T]) Dependencies() *dependencies.Dependencies[T] {
	return s.deps
}

// WithCache injects a custom caching implementation. Passing nil removes the
// cache dependency.
func (s *Server[T]) WithCache(c cache.Cacher) {
	s.deps.Cache = c
}

// WithDatabase injects a custom database implementation. Passing the zero value
// removes the dependency.
func (s *Server[T]) WithDatabase(d T) {
	s.deps.DB = d
}

// Use registers middleware to run after the built-in middleware stack.
func (s *Server[T]) Use(mw gin.HandlerFunc) {
	s.deps.Router.Use(mw)
}

// RegisterRoutes allows modules to add routes to the server.
func (s *Server[T]) RegisterRoutes(fn RouterFunc[T]) {
	fn(s.deps)
}

// NewServer creates a new Server using the provided configuration options. Only
// the router and tracer are initialized by default, leaving cache and database
// setup to the caller.
func NewServer[T any](opts config.ConfigOptions) (*Server[T], error) {
	if err := config.Init(opts); err != nil {
		return nil, err
	}
	cfg := config.GetConfig()

	tracer := appTracer.NewAppTracer(cfg)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(otelgin.Middleware(cfg.AppName))
	router.Use(middleware.ClientContextMiddleware())
	router.Use(middleware.TraceMiddleware(cfg.AppName))
	router.Use(middleware.ErrorHandler)
	router.Use(middleware.JsonLogger())

	deps := &dependencies.Dependencies[T]{
		Router: router,
		Tracer: tracer,
	}

	return &Server[T]{config: cfg, deps: deps}, nil
}

// Run starts the HTTP server and blocks until a shutdown signal is received.
func (s *Server[T]) Run() error {
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port),
		Handler: s.deps.Router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("failed to start server: %v", err)
		}
	}()

	// Listen for interrupt and terminate signals to shut down gracefully.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("Server Shutdown:", err)
	}

	select {
	case <-ctx.Done():
		logrus.Info("timeout of 5 seconds.")
	}
	logrus.Info("Server exiting")
	return nil
}
