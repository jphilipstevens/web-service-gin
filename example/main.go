package main

// @title Web Service Gin Example
// @version 1.0
// @description Example API demonstrating middleware and Swagger docs.
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

import (
	"flag"
	"os"

	"github.com/jphilipstevens/web-service-gin/app"
	"github.com/jphilipstevens/web-service-gin/app/dependencies"
	"github.com/jphilipstevens/web-service-gin/config"
	"github.com/jphilipstevens/web-service-gin/example/features/albums"
	"github.com/jphilipstevens/web-service-gin/example/seed"

	"github.com/jphilipstevens/web-service-gin/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func registerRoutes(deps *dependencies.Dependencies) {
	albums.Init(deps)
	deps.Router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func RunApp() {
	docs.SwaggerInfo.BasePath = "/"

	app.RunServer(app.ServerParams{
		Routes: registerRoutes,
		ConfigOptions: config.ConfigOptions{
			Path: "./example/config", // or from ENV, flags, etc
			Name: "config",           // without extension
			Type: "yaml",             // optional
		},
	})
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
