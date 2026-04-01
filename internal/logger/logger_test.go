package logger

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		config      *LoggerConfig
		wantErr     bool
		description string
	}{
		{
			name: "JSON format logger",
			config: &LoggerConfig{
				Level:  "info",
				Format: "json",
			},
			wantErr:     false,
			description: "Should create JSON logger",
		},
		{
			name: "Console format logger",
			config: &LoggerConfig{
				Level:  "debug",
				Format: "console",
			},
			wantErr:     false,
			description: "Should create console logger",
		},
		{
			name: "Invalid level fallback",
			config: &LoggerConfig{
				Level:  "invalid",
				Format: "json",
			},
			wantErr:     false,
			description: "Should fallback to info level for invalid level",
		},
		{
			name: "Invalid format fallback",
			config: &LoggerConfig{
				Level:  "info",
				Format: "invalid",
			},
			wantErr:     false,
			description: "Should fallback to JSON format for invalid format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := New(tt.config)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, logger)
		})
	}
}

func TestLoggerWithContext(t *testing.T) {
	config := &LoggerConfig{
		Level:  "info",
		Format: "json",
	}

	logger, err := New(config)
	require.NoError(t, err)

	// Test WithRequestID
	withReqID := logger.WithRequestID("test-123")
	assert.NotNil(t, withReqID)

	// Test WithUserID
	withUserID := logger.WithUserID("user-456")
	assert.NotNil(t, withUserID)

	// Test WithField
	withField := logger.WithField("key", "value")
	assert.NotNil(t, withField)

	// Test WithFields
	withFields := logger.WithFields(map[string]interface{}{
		"field1": "value1",
		"field2": "value2",
	})
	assert.NotNil(t, withFields)

	// Test chaining
	chained := logger.WithRequestID("req-123").WithUserID("user-456").WithField("action", "test")
	assert.NotNil(t, chained)
}

func TestLoggerMethods(t *testing.T) {
	// Create a logger that writes to a buffer for testing
	var buf bytes.Buffer
	config := &LoggerConfig{
		Level:  "info",
		Format: "json",
	}

	logger, err := New(config)
	require.NoError(t, err)

	// Redirect logger output to buffer for testing
	logger.Logger = logger.Logger.Output(&buf)

	// Test HTTPRequest
	logger.HTTPRequest("GET", "/test", "127.0.0.1", 200, 150)
	output := buf.String()
	assert.Contains(t, output, "HTTP request")
	assert.Contains(t, output, "GET")
	assert.Contains(t, output, "/test")
	assert.Contains(t, output, "200")

	// Reset buffer
	buf.Reset()

	// Test HTTPError
	logger.HTTPError("POST", "/error", "127.0.0.1", 500, assert.AnError)
	output = buf.String()
	assert.Contains(t, output, "HTTP request error")
	assert.Contains(t, output, "POST")
	assert.Contains(t, output, "500")

	// Reset buffer
	buf.Reset()

	// Test DatabaseQuery
	logger.DatabaseQuery("SELECT * FROM users", []interface{}{}, 25, nil)
	output = buf.String()
	assert.Contains(t, output, "Database query")
	assert.Contains(t, output, "SELECT * FROM users")

	// Reset buffer
	buf.Reset()

	// Test GracefulShutdown
	logger.GracefulShutdown("test-component", nil)
	output = buf.String()
	assert.Contains(t, output, "Graceful shutdown completed")
	assert.Contains(t, output, "test-component")

	// Reset buffer
	buf.Reset()

	// Test GracefulShutdown with error
	logger.GracefulShutdown("test-component", assert.AnError)
	output = buf.String()
	assert.Contains(t, output, "Graceful shutdown error")
	assert.Contains(t, output, "test-component")

	// Reset buffer
	buf.Reset()

	// Test Startup
	logger.Startup("test-app", "1.0.0", "localhost:8080", "info")
	output = buf.String()
	assert.Contains(t, output, "Application starting")
	assert.Contains(t, output, "test-app")
	assert.Contains(t, output, "1.0.0")

	// Reset buffer
	buf.Reset()

	// Test Shutdown
	logger.Shutdown("test reason")
	output = buf.String()
	assert.Contains(t, output, "Application shutting down")
	assert.Contains(t, output, "test reason")
}

func TestLoggerFormat(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		expected string
	}{
		{
			name:     "JSON format",
			format:   "json",
			expected: "{",
		},
		{
			name:     "Console format",
			format:   "console",
			expected: "info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			config := &LoggerConfig{
				Level:  "info",
				Format: tt.format,
			}

			logger, err := New(config)
			require.NoError(t, err)

			// Redirect output to buffer
			logger.Logger = logger.Logger.Output(&buf)

			// Log a message
			logger.Logger.Info().Msg("test message")

			output := buf.String()
			if tt.expected == "{" {
				assert.True(t, strings.HasPrefix(output, "{"))
			} else {
				assert.Contains(t, output, tt.expected)
			}
		})
	}
}
