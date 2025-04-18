package seed

import (
	"fmt"

	"github.com/jphilipstevens/web-service-gin/app/db"
	"github.com/jphilipstevens/web-service-gin/config"
	"github.com/jphilipstevens/web-service-gin/testUtils"
)

func Init() {
	config.Init(config.ConfigOptions{
		Path: "./example/config", // or from ENV, flags, etc
		Name: "config",           // without extension
		Type: "yaml",             // optional
	})
	configFile := config.GetConfig()

	// Initialize database connection
	// TODO: Review using testing app tracer. Not needed for seeding. But may be needed for the future
	dbConn, err := db.NewDatabase(configFile.DB, testUtils.NewAppTracer())
	if err != nil {
		// Handle error
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}

	// Create the albums table if it doesn't exist
	if err := SeedAlbums(dbConn); err != nil {
		panic(fmt.Errorf("fatal error cannot create Album Table: %w", err))
	}
}
