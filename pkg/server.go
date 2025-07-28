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
	router            *gin.Engine
	config            config.Config
	preMiddleware     []gin.HandlerFunc
	postMiddleware    []gin.HandlerFunc
	finalMiddleware   []gin.HandlerFunc
	middlewareApplied bool
}

// New creates a Server with the provided configuration. Middleware is applied
// when Run is invoked so callers can register their own hooks beforehand.
func New(cfg config.Config) *Server {
	middleware.SetupLogger()

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

// UseFinal registers middleware that executes after the route handlers.
// Each final middleware must call c.Next() to ensure the request flows to the
// handlers before running its cleanup logic.
func (s *Server) UseFinal(mw ...gin.HandlerFunc) {
	s.finalMiddleware = append(s.finalMiddleware, mw...)
}

// PreMiddleware returns the middleware registered with UseBefore.
// The returned slice must not be modified by the caller.
func (s *Server) PreMiddleware() []gin.HandlerFunc {
	return append([]gin.HandlerFunc(nil), s.preMiddleware...)
}

// PostMiddleware returns the middleware registered with UseAfter.
// The returned slice must not be modified by the caller.
func (s *Server) PostMiddleware() []gin.HandlerFunc {
	return append([]gin.HandlerFunc(nil), s.postMiddleware...)
}

// FinalMiddleware returns the middleware registered with UseFinal.
// The returned slice must not be modified by the caller.
func (s *Server) FinalMiddleware() []gin.HandlerFunc {
	return append([]gin.HandlerFunc(nil), s.finalMiddleware...)
}

// BuiltInMiddleware returns the core middleware stack in the order it will be
// applied. The slice must not be modified by the caller.
func (s *Server) BuiltInMiddleware() []gin.HandlerFunc {
	stack := []gin.HandlerFunc{
		gin.Recovery(),
	}
	if s.config.EnableOpenTelemetry {
		stack = append(stack, otelgin.Middleware(s.config.AppName))
	}
	stack = append(stack,
		middleware.ClientContextMiddleware(),
		middleware.TraceMiddleware(s.config.AppName),
		middleware.ErrorHandler,
		middleware.JsonLogger(),
	)
	return stack
}

// FullMiddlewareChain returns the ordered chain that will be applied to the router.
// It consists of PreMiddleware, the built-in stack, PostMiddleware and FinalMiddleware.
// The returned slice must not be modified by the caller.
func (s *Server) FullMiddlewareChain() []gin.HandlerFunc {
	chain := make([]gin.HandlerFunc, 0,
		len(s.preMiddleware)+len(s.postMiddleware)+len(s.finalMiddleware)+6)
	chain = append(chain, s.preMiddleware...)
	chain = append(chain, s.BuiltInMiddleware()...)
	chain = append(chain, s.postMiddleware...)
	chain = append(chain, s.finalMiddleware...)
	return chain
}

// Router exposes the underlying Gin engine for advanced usage.
// Handlers must call ApplyMiddleware before serving requests.
func (s *Server) Router() *gin.Engine {
	return s.router
}

// ApplyMiddleware attaches registered middleware to the router in
// the correct order. It is idempotent so repeated calls have no effect.
func (s *Server) ApplyMiddleware() {
	if s.middlewareApplied {
		return
	}

	for _, mw := range s.preMiddleware {
		s.router.Use(mw)
	}

	for _, mw := range s.BuiltInMiddleware() {
		s.router.Use(mw)
	}

	for _, mw := range s.postMiddleware {
		s.router.Use(mw)
	}

	for _, mw := range s.finalMiddleware {
		s.router.Use(mw)
	}

	s.middlewareApplied = true
}

// RegisterRoutes allows modules to add routes to the server.
func (s *Server) RegisterRoutes(fn func(r *gin.Engine)) {
	s.ApplyMiddleware()
	fn(s.router)
}

// Run starts the HTTP server, applies middleware in the correct order, and
// waits for a shutdown signal.
func (s *Server) Run() error {
	s.ApplyMiddleware()

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port),
		Handler: s.router,
	}

	logrus.Infof("🚀 Listening on host %s, port %d", s.config.Server.Host, s.config.Server.Port)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("failed to start server: %v", err)
		}
	}()

	return GracefulShutdown(srv, gracefulShutdownTimeout)
}

// GracefulShutdown blocks until an interrupt signal is received and then attempts
// to shut down the provided server within the given timeout.
func GracefulShutdown(srv *http.Server, timeout time.Duration) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	logrus.Info("server exited properly")
	return nil
}
