package main

import (
	"context"
	"example/web-service-gin/app/albums"
	"example/web-service-gin/app/cache"
	"example/web-service-gin/app/db"
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
	"github.com/spf13/viper"
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")

	// Load configuration
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	var configFile struct {
		Redis cache.RedisClientConfig `mapstructure:"redis"`
		DB    db.DatabaseConfig       `mapstructure:"database"`
	}

	if err := viper.Unmarshal(&configFile); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// Initialize Redis client
	redisClient := cache.NewCacher(configFile.Redis)

	// Initialize database connection
	dbConn, err := db.ConnectToDB(configFile.DB)
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}

	router := gin.Default()
	router.Use(middleware.ErrorHandler)
	router.Use(middleware.JsonLogger())

	var dependencies = config.Dependencies{
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
