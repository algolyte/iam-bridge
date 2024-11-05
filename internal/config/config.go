package config

import (
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// ProviderSet is a provider set for wire
var ProviderSet = wire.NewSet(NewConfig)

type Config struct {
	AppName            string   `mapstructure:"APP_NAME"`
	AppEnv             string   `mapstructure:"APP_ENV"`
	ServerPort         string   `mapstructure:"APP_PORT"`
	LogLevel           string   `mapstructure:"LOG_LEVEL"`
	CorsAllowedOrigins []string `mapstructure:"CORS_ALLOWED_ORIGINS"`
	CorsAllowedMethods []string `mapstructure:"CORS_ALLOWED_METHODS"`
	CorsAllowedHeaders []string `mapstructure:"CORS_ALLOWED_HEADERS"`
	CorsMaxAge         int      `mapstructure:"CORS_MAX_AGE"`
}

func NewConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigFile(".env")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "failed to read config")
	}

	config := &Config{}
	if err := v.Unmarshal(config); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal config")
	}

	// Set defaults
	if config.ServerPort == "" {
		config.ServerPort = "8080"
	}
	if config.LogLevel == "" {
		config.LogLevel = "info"
	}
	if config.AppEnv == "" {
		config.AppEnv = "development"
	}

	return config, nil
}
