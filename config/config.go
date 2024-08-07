package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type RedisClientConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`

	DBName  string `mapstructure:"dbname"`
	SSLMode string `mapstructure:"sslmode"`
}

type ConfigFile struct {
	AppName string            `mapstructure:"app_name"`
	Redis   RedisClientConfig `mapstructure:"redis"`
	DB      DatabaseConfig    `mapstructure:"database"`
}

var configFile ConfigFile

func GetConfig() ConfigFile {
	if configFile == (ConfigFile{}) {
		panic(fmt.Errorf("Config File not initialized. This indicates that the main app was not setup correctly. Make sure to call config.Init() in main.go"))

	}
	return configFile
}

func Init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")

	// Load configuration
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	configFile = ConfigFile{}

	if err := viper.Unmarshal(&configFile); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

}
