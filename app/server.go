package app

import (
	"context"
	"example/web-service-gin/app/appTracer"
	"example/web-service-gin/app/cache"
	"example/web-service-gin/app/db"
	"example/web-service-gin/app/dependencies"
	"example/web-service-gin/app/middleware"
	"example/web-service-gin/config"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

const (
	gracefulShutdownTimeout = 5 * time.Second
)

// RouterFunc will pass all dependencies and let you initialize your routes
// you have access to dependencies.Router which is a *gin.Engine
type RouterFunc func(dependencies *dependencies.Dependencies)

// initGracefulShutdown sets up a graceful shutdown mechanism for the HTTP server.
// It listens for interrupt signals (SIGINT and SIGTERM) and initiates a
// graceful shutdown when received. The server is given a timeout of 5 seconds
// to finish ongoing requests before forcefully shutting down.
//
// Parameters:
//   - srv: A pointer to the http.Server that should be gracefully shut down.
//
// The function blocks until the shutdown is complete or the timeout is reached.
// It logs the shutdown process and any errors that occur during shutdown.

func initGracefulShutdown(srv *http.Server) {

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)

	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		logrus.Info("timeout of 5 seconds.")
	}
	logrus.Info("Server exiting")
}

// ServerParams is a struct that contains all the dependencies needed to run the server
// if ServerParams.Dependencies is nil, it will be initialized with default values
//   - DB is postgres
//   - Cache is redis
//   - Router is gin.Default()
type ServerParams struct {
	Routes       RouterFunc
	Dependencies *dependencies.Dependencies
}

func RunServer(ServerParams ServerParams) {

	err := config.Init()
	if err != nil {
		panic(err)
	}
	configFile := config.GetConfig()

	// Initialize Redis client
	appTracer := appTracer.NewAppTracer(configFile)
	defer uptrace.Shutdown(context.Background())
	redisClient := cache.NewCacher(configFile.Redis, appTracer)

	// Initialize database connection
	dbConn, err := db.NewDatabase(configFile.DB, appTracer)
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}

	router := gin.Default()
	router.Use(otelgin.Middleware(configFile.AppName))
	router.Use(middleware.ClientContextMiddleware())
	router.Use(middleware.TraceMiddleware(configFile.AppName))
	router.Use(middleware.ErrorHandler)
	router.Use(middleware.JsonLogger())

	var serverDependencies *dependencies.Dependencies
	if ServerParams.Dependencies == nil {
		serverDependencies = &dependencies.Dependencies{
			Cache:  redisClient,
			DB:     dbConn,
			Router: router,
		}
	} else {
		serverDependencies = ServerParams.Dependencies
	}

	ServerParams.Routes(serverDependencies)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", configFile.Server.Host, configFile.Server.Port),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Errorf("failed to start server: %w", err))
		}
	}()

	initGracefulShutdown(srv)
}
