package seed

import (
	"example/web-service-gin/app/db"
	"example/web-service-gin/config"
	"fmt"
)

func Init() {
	config.Init()
	configFile := config.GetConfig()

	// Initialize database connection
	dbConn, err := db.ConnectToDB(configFile.DB)
	if err != nil {
		// Handle error
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}

	// Create the albums table if it doesn't exist
	if err := SeedAlbums(dbConn); err != nil {
		panic(fmt.Errorf("fatal error cannot create Album Table: %w", err))
	}
}
