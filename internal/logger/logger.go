package logger

import (
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

// LoggerConfig holds logging configuration
type LoggerConfig struct {
	Level  string `mapstructure:"level" validate:"required,oneof=trace debug info warn error"`
	Format string `mapstructure:"format" validate:"required,oneof=json console"`
}

// Logger wraps zerolog logger with additional functionality
type Logger struct {
	zerolog.Logger
}

// New creates a new logger instance based on configuration
func New(cfg *LoggerConfig) (*Logger, error) {
	var output io.Writer = os.Stdout

	// Configure output format
	switch strings.ToLower(cfg.Format) {
	case "json":
		output = os.Stdout
	case "console":
		output = zerolog.ConsoleWriter{Out: os.Stdout}
	default:
		output = os.Stdout
	}

	// Parse log level
	level, err := zerolog.ParseLevel(strings.ToLower(cfg.Level))
	if err != nil {
		level = zerolog.InfoLevel // fallback to info
	}

	// Create logger
	zerolog.SetGlobalLevel(level)
	logger := zerolog.New(output).
		With().
		Timestamp().
		Caller().
		Logger()

	return &Logger{Logger: logger}, nil
}

// WithRequestID adds request ID to logger context
func (l *Logger) WithRequestID(requestID string) *Logger {
	return &Logger{Logger: l.Logger.With().Str("request_id", requestID).Logger()}
}

// WithUserID adds user ID to logger context
func (l *Logger) WithUserID(userID string) *Logger {
	return &Logger{Logger: l.Logger.With().Str("user_id", userID).Logger()}
}

// WithError adds error to logger context
func (l *Logger) WithError(err error) *Logger {
	return &Logger{Logger: l.Logger.With().Err(err).Logger()}
}

// WithField adds a custom field to logger context
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{Logger: l.Logger.With().Interface(key, value).Logger()}
}

// WithFields adds multiple fields to logger context
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	ctx := l.Logger.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	return &Logger{Logger: ctx.Logger()}
}

// HTTPRequest logs HTTP request information
func (l *Logger) HTTPRequest(method, path, remoteAddr string, statusCode int, duration int64) {
	l.Logger.Info().
		Str("method", method).
		Str("path", path).
		Str("remote_addr", remoteAddr).
		Int("status_code", statusCode).
		Int64("duration_ms", duration).
		Msg("HTTP request")
}

// HTTPError logs HTTP error information
func (l *Logger) HTTPError(method, path, remoteAddr string, statusCode int, err error) {
	l.Logger.Error().
		Err(err).
		Str("method", method).
		Str("path", path).
		Str("remote_addr", remoteAddr).
		Int("status_code", statusCode).
		Msg("HTTP request error")
}

// DatabaseQuery logs database query information
func (l *Logger) DatabaseQuery(query string, args []interface{}, duration int64, err error) {
	event := l.Logger.Info().
		Str("query", query).
		Int64("duration_ms", duration).
		Int("args_count", len(args))

	if err != nil {
		event.Err(err)
	}

	event.Msg("Database query")
}

// GracefulShutdown logs graceful shutdown information
func (l *Logger) GracefulShutdown(component string, err error) {
	if err != nil {
		l.Logger.Error().
			Err(err).
			Str("component", component).
			Msg("Graceful shutdown error")
	} else {
		l.Logger.Info().
			Str("component", component).
			Msg("Graceful shutdown completed")
	}
}

// Startup logs application startup information
func (l *Logger) Startup(appName, version string, serverAddr, logLevel string) {
	l.Logger.Info().
		Str("app_name", appName).
		Str("version", version).
		Str("server_addr", serverAddr).
		Str("log_level", logLevel).
		Msg("Application starting")
}

// Shutdown logs application shutdown information
func (l *Logger) Shutdown(reason string) {
	l.Logger.Info().
		Str("reason", reason).
		Msg("Application shutting down")
}
