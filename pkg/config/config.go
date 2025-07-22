// Package config provides idiomatic Go configuration handling
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

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Driver          string        `mapstructure:"driver"`
	DBName          string        `mapstructure:"dbname"`
	SSLMode         string        `mapstructure:"sslmode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type UptraceConfig struct {
	DSN      string `mapstructure:"dsn"`
	Endpoint string `mapstructure:"endpoint"`
}

type Config struct {
	AppName string         `mapstructure:"app_name"`
	Server  ServerConfig   `mapstructure:"server"`
	Redis   RedisConfig    `mapstructure:"redis"`
	DB      DatabaseConfig `mapstructure:"database"`
	Uptrace UptraceConfig  `mapstructure:"uptrace"`
}

var globalConfig Config

func Load(path, name, filetype string) (Config, error) {
	viper.SetConfigName(name)
	viper.SetConfigType(filetype)
	viper.AddConfigPath(path)

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("could not read config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("could not unmarshal config: %w", err)
	}

	globalConfig = cfg
	return cfg, nil
}

func Get() Config {
	if globalConfig == (Config{}) {
		panic("config not initialized: call config.Load first")
	}
	return globalConfig
}
