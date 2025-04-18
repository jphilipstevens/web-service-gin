package main

import (
	"flag"
	"os"

	"github.com/jphilipstevens/web-service-gin/app"
	"github.com/jphilipstevens/web-service-gin/config"
	"github.com/jphilipstevens/web-service-gin/example/features/albums"
	"github.com/jphilipstevens/web-service-gin/example/seed"
)

func RunApp() {

	app.RunServer(app.ServerParams{
		Routes: albums.Init,
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
