package main

import (
	"example/web-service-gin/app/albums"
	"example/web-service-gin/app/cache"
	"example/web-service-gin/app/db"
	"example/web-service-gin/app/middleware"
	"example/web-service-gin/config"
	"example/web-service-gin/seed"
	"flag"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

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

	var dependencies = config.Dependencies{
		Cache:  redisClient,
		DB:     dbConn,
		Router: router,
	}

	albums.Init(
		&dependencies,
	)

	router.Run("localhost:8080")
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
