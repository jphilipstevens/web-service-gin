package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

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
	Driver   string `mapstructure:"driver"`

	DBName  string `mapstructure:"dbname"`
	SSLMode string `mapstructure:"sslmode"`

	MaxOpenConns    int           `mapstructure:"maxOpenConns"`
	MaxIdleConns    int           `mapstructure:"maxIdleConns"`
	ConnMaxLifetime time.Duration `mapstructure:"connMaxLifetime"`
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
	Server  ServerConfig      `mapstructure:"server"`
}

type ConfigOptions struct {
	Path string
	Name string
	Type string // optional: json, yaml, toml
}

var configFile ConfigFile

func GetConfig() ConfigFile {
	if configFile == (ConfigFile{}) {
		panic(fmt.Errorf("Config File not initialized. This indicates that the main app was not setup correctly. Make sure to call config.Init() in main.go"))

	}
	return configFile
}

// Config entries can be specified in the config file and overridden using
// environment variables. To map an environment variable to a config field,
// replace dots in the key with underscores and use uppercase letters. For
// example, the key "redis.host" becomes "REDIS_HOST".
func Init(opts ConfigOptions) error {
	viper.SetConfigName(opts.Name)
	viper.AddConfigPath(opts.Path)

	if opts.Type != "" {
		viper.SetConfigType(opts.Type)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

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
