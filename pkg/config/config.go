package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// ServerConfig holds the HTTP server host and port settings.
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// UptraceConfig holds the DSN used to connect to Uptrace for tracing.
type UptraceConfig struct {
	DSN string `mapstructure:"dsn"`
}

// Config represents the minimal application configuration required to start
// the HTTP server. Additional configuration should be loaded separately by the
// modules that need it.
type Config struct {
	AppName string        `mapstructure:"app_name"`
	Server  ServerConfig  `mapstructure:"server"`
	Uptrace UptraceConfig `mapstructure:"uptrace"`
	// EnableOpenTelemetry toggles insertion of tracing middleware.
	// Defaults to true for backward compatibility.
	EnableOpenTelemetry bool `mapstructure:"enable_open_telemetry"`
}

// ConfigOptions defines how the configuration file is located and parsed.
type ConfigOptions struct {
	Path string
	Name string
	Type string // optional: json, yaml, toml
}

var cfg Config

// Get returns the loaded configuration. Init must be called before using Get
// otherwise the function panics to indicate a misconfigured application.
func Get() Config {
	if cfg == (Config{}) {
		panic(fmt.Errorf("config not initialized; call config.Init() in main"))
	}
	return cfg
}

// Init loads configuration from file and environment variables. Environment
// variables override values from the file using `KEY=value` where the key is the
// config path with dots replaced by underscores (e.g. `server.port`).
func Init(opts ConfigOptions) error {
	viper.SetConfigName(opts.Name)
	viper.AddConfigPath(opts.Path)

	if opts.Type != "" {
		viper.SetConfigType(opts.Type)
	}

	// Preserve backward compatibility by defaulting tracing to enabled.
	viper.SetDefault("enable_open_telemetry", true)

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("fatal error config file: %w", err)
	}
	cfg = Config{}
	if err := viper.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("fatal error config file: %w", err)
	}
	return nil
}
