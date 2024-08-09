package config

import (
	"fmt"
	"strings"

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

type UptraceConfig struct {
	DSN      string `mapstructure:"dsn"`
	Endpoint string `mapstructure:"endpoint"`
}

type ConfigFile struct {
	AppName string            `mapstructure:"app_name"`
	Redis   RedisClientConfig `mapstructure:"redis"`
	DB      DatabaseConfig    `mapstructure:"database"`
	Uptrace UptraceConfig     `mapstructure:"uptrace"`
}

var configFile ConfigFile

func GetConfig() ConfigFile {
	if configFile == (ConfigFile{}) {
		panic(fmt.Errorf("Config File not initialized. This indicates that the main app was not setup correctly. Make sure to call config.Init() in main.go"))

	}
	return configFile
}

// Config entries can be set in the config file or as environment variables.
// When set as environment variables, the key should be in the format where the dot notation is replaced with an underscore.
// For example, the key "redis.host" can be set as the environment variable "REDIS_HOST"
func Init() error {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Load configuration
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("fatal error config file: %w", err)
	}
	configFile = ConfigFile{}

	if err := viper.Unmarshal(&configFile); err != nil {
		return fmt.Errorf("fatal error config file: %w", err)
	}
	return nil
}
