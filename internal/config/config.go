package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

// Config holds all configuration for our program
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	IAM      IAMConfig      `mapstructure:"iam"`
	Security SecurityConfig `mapstructure:"security"`
	Logging  LogConfig      `mapstructure:"logging"`
}

// AppConfig holds all application configuration
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Environment string `mapstructure:"environment"`
	Port        int    `mapstructure:"port"`
	Debug       bool   `mapstructure:"debug"`
}

// KeycloakConfig holds Keycloak-specific configuration
type KeycloakConfig struct {
	BaseURL      string `mapstructure:"base_url"`
	Realm        string `mapstructure:"realm"`
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

// CORSConfig holds CORS-related configuration
type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Enabled          bool `mapstructure:"enabled"`
	RequestPerSecond int  `mapstructure:"request_per_second"`
}

// LogConfig holds logging-related configuration
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	CORS            CORSConfig      `mapstructure:"cors"`
	RateLimitConfig RateLimitConfig `mapstructure:"rate_limit"`
}

// IAMConfig holds the configuration for IAM providers
type IAMConfig struct {
	Provider string         `mapstructure:"provider"`
	Keycloak KeycloakConfig `mapstructure:"keycloak"`
}

// LoadConfig reads configuration from file or environment variables
func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			return nil, fmt.Errorf("config file not found: %w", err)
		}
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	return &config, nil
}

// CurrentProvider returns the configured IAM provider name
func (c *IAMConfig) CurrentProvider() string {
	return strings.ToLower(c.Provider)
}

// IsDebug returns true if the application is in debug mode
func (c *Config) IsDebug() bool {
	return c.App.Debug
}

// IsDevelopment returns true if the application is in development mode
func (c *AppConfig) IsDevelopment() bool {
	return strings.ToLower(c.Environment) == "development"
}
