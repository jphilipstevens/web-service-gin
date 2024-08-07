package seed

import (
	"example/web-service-gin/app/db"
	"example/web-service-gin/config"
	"example/web-service-gin/testUtils"
	"fmt"
)

func Init() {
	config.Init()
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
