package main

import (
	"flag"
	"jphilipstevens/web-service-gin/app"
	"jphilipstevens/web-service-gin/config"
	"jphilipstevens/web-service-gin/example/features/albums"
	"jphilipstevens/web-service-gin/example/seed"
	"os"
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
