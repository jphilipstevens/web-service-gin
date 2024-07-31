package seed

import (
	"example/web-service-gin/app/db"
	"fmt"

	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")

	// Load configuration
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	var config struct {
		DB db.DatabaseConfig `mapstructure:"database"`
	}

	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// Initialize database connection
	dbConn, err := db.ConnectToDB(config.DB)
	if err != nil {
		// Handle error
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}

	// Create the albums table if it doesn't exist
	if err := SeedAlbums(dbConn); err != nil {
		panic(fmt.Errorf("fatal error cannot create Album Table: %w", err))
	}
}
