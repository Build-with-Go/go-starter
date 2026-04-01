// Package config provides configuration management for the Go Starter application.
package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig   `mapstructure:"server" validate:"required"`
	Database DatabaseConfig `mapstructure:"database" validate:"required"`
	Logger   LoggerConfig   `mapstructure:"logger" validate:"required"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Host         string `mapstructure:"host" validate:"required"`
	Port         int    `mapstructure:"port" validate:"required,min=1,max=65535"`
	ReadTimeout  int    `mapstructure:"read_timeout" validate:"required,min=1"`
	WriteTimeout int    `mapstructure:"write_timeout" validate:"required,min=1"`
	IdleTimeout  int    `mapstructure:"idle_timeout" validate:"required,min=1"`
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host     string `mapstructure:"host" validate:"required"`
	Port     int    `mapstructure:"port" validate:"required,min=1,max=65535"`
	Name     string `mapstructure:"name" validate:"required"`
	User     string `mapstructure:"user" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
	SSLMode  string `mapstructure:"ssl_mode" validate:"required,oneof=disable enable require"`
}

// LoggerConfig holds logging configuration
type LoggerConfig struct {
	Level  string `mapstructure:"level" validate:"required,oneof=trace debug info warn error"`
	Format string `mapstructure:"format" validate:"required,oneof=json console"`
}

// Load loads configuration from file and environment variables
func Load(configPath string) (*Config, error) {
	config := &Config{}

	// Set config file path and name
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./configs")
		viper.AddConfigPath(".")
	}

	// Set environment variable prefix
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()

	// Set default values
	setDefaults()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; use defaults and env vars
			fmt.Printf("Config file not found, using defaults and environment variables\n")
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Unmarshal config
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Validate config
	if err := validate(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", 30)
	viper.SetDefault("server.write_timeout", 30)
	viper.SetDefault("server.idle_timeout", 120)

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.name", "go_starter")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "") // Empty default - must be set via environment
	viper.SetDefault("database.ssl_mode", "disable")

	// Logger defaults
	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.format", "json")
}

// validate validates the configuration
func validate(config *Config) error {
	// Basic validation - in production, you might want to use a validation library
	if config.Server.Host == "" {
		return fmt.Errorf("server host is required")
	}
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("server port must be between 1 and 65535")
	}
	if config.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if config.Database.Name == "" {
		return fmt.Errorf("database name is required")
	}
	if config.Database.User == "" {
		return fmt.Errorf("database user is required")
	}
	if config.Logger.Level == "" {
		return fmt.Errorf("logger level is required")
	}

	return nil
}

// GetAddr returns the server address in host:port format
func (c *Config) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode,
	)
}
