package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		configPath  string
		envVars     map[string]string
		wantErr     bool
		expected    *Config
		description string
	}{
		{
			name:       "Default configuration",
			configPath: "",
			envVars:    nil,
			wantErr:    false,
			expected: &Config{
				Server: ServerConfig{
					Host:         "localhost",
					Port:         8080,
					ReadTimeout:  30,
					WriteTimeout: 30,
					IdleTimeout:  120,
				},
				Database: DatabaseConfig{
					Host:     "localhost",
					Port:     5432,
					Name:     "go_starter",
					User:     "postgres",
					Password: "password",
					SSLMode:  "disable",
				},
				Logger: LoggerConfig{
					Level:  "info",
					Format: "json",
				},
			},
			description: "Should load default configuration when no file provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Load configuration
			got, err := Load(tt.configPath)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected.Server.Host, got.Server.Host)
			assert.Equal(t, tt.expected.Server.Port, got.Server.Port)
			assert.Equal(t, tt.expected.Logger.Level, got.Logger.Level)
		})
	}
}

func TestConfigGetAddr(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
	}

	expected := "localhost:8080"
	assert.Equal(t, expected, cfg.GetAddr())
}

func TestConfigGetDSN(t *testing.T) {
	cfg := &Config{
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			Name:     "test_db",
			User:     "test_user",
			Password: "test_pass",
			SSLMode:  "disable",
		},
	}

	expected := "host=localhost port=5432 user=test_user password=test_pass dbname=test_db sslmode=disable"
	assert.Equal(t, expected, cfg.GetDSN())
}

func TestValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		wantErr     bool
		description string
	}{
		{
			name: "Valid configuration",
			config: &Config{
				Server: ServerConfig{
					Host: "localhost",
					Port: 8080,
				},
				Database: DatabaseConfig{
					Host: "localhost",
					Name: "test",
					User: "user",
				},
				Logger: LoggerConfig{
					Level: "info",
				},
			},
			wantErr:     false,
			description: "Should pass validation with valid config",
		},
		{
			name: "Invalid server port",
			config: &Config{
				Server: ServerConfig{
					Host: "localhost",
					Port: 0, // Invalid port
				},
			},
			wantErr:     true,
			description: "Should fail validation with invalid port",
		},
		{
			name: "Missing database host",
			config: &Config{
				Server: ServerConfig{
					Host: "localhost",
					Port: 8080,
				},
				Database: DatabaseConfig{
					Host: "", // Missing host
					Name: "test",
					User: "user",
				},
			},
			wantErr:     true,
			description: "Should fail validation with missing database host",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
