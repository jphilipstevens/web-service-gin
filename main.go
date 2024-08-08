package main

import (
	"context"
	"example/web-service-gin/app/albums"
	"example/web-service-gin/app/appTracer"
	"example/web-service-gin/app/cache"
	"example/web-service-gin/app/db"
	"example/web-service-gin/app/dependencies"
	"example/web-service-gin/app/middleware"
	"example/web-service-gin/config"
	"example/web-service-gin/seed"
	"flag"
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

func RunApp() {

	err := config.Init()
	if err != nil {
		panic(err)
	}
	configFile := config.GetConfig()

	// Initialize Redis client
	appTracer := appTracer.NewDownstreamSpan(configFile)
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

	var dependencies = dependencies.Dependencies{
		Cache:  redisClient,
		DB:     dbConn,
		Router: router,
	}

	albums.Init(
		&dependencies,
	)

	srv := &http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Errorf("failed to start server: %w", err))
		}
	}()

	initGracefulShutdown(srv)
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) > 0 {
		switch args[0] {
		case "seed":
			seed.Init()
			os.Exit(0)
		default:
			RunApp()
		}
	} else {
		RunApp()
	}
}
