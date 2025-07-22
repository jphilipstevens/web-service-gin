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

	"yourapp/pkg/appTracer"
	"yourapp/pkg/config"
	"yourapp/pkg/dependencies"
	"yourapp/pkg/middleware"
)

const gracefulShutdownTimeout = 5 * time.Second

type RouterFunc func(deps *dependencies.Dependencies)

type Server struct {
	config config.Config
	deps   *dependencies.Dependencies
}

func (s *Server) Config() config.Config {
	return s.config
}

func (s *Server) Dependencies() *dependencies.Dependencies {
	return s.deps
}

func (s *Server) Use(mw gin.HandlerFunc) {
	s.deps.Router.Use(mw)
}

func (s *Server) RegisterRoutes(fn RouterFunc) {
	fn(s.deps)
}

func New(cfg config.Config) *Server {
	tracer := appTracer.NewAppTracer(cfg)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(otelgin.Middleware(cfg.AppName))
	router.Use(middleware.ClientContextMiddleware())
	router.Use(middleware.TraceMiddleware(cfg.AppName))
	router.Use(middleware.ErrorHandler)
	router.Use(middleware.JsonLogger())

	deps := &dependencies.Dependencies{
		Router: router,
		Tracer: tracer,
	}

	return &Server{
		config: cfg,
		deps:   deps,
	}
}

func (s *Server) Run() error {
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port),
		Handler: s.deps.Router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("failed to start server: %v", err)
		}
	}()

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
