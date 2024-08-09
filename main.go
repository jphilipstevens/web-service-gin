package main

import (
	"example/web-service-gin/app"
	"example/web-service-gin/features/albums"
	"example/web-service-gin/seed"
	"flag"
	"os"
)

func RunApp() {

	app.RunServer(app.ServerParams{
		Routes: albums.Init,
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
