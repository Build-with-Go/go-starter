# Architecture Documentation

## Overview

Go Starter is a production-ready Go project template that follows idiomatic Go patterns and best practices. It provides a solid foundation for building web applications and microservices with minimal configuration and maximum flexibility.

## Design Philosophy

### Core Principles

1. **Idiomatic Go**: Follow Go conventions and standard library patterns
2. **Minimal Magic**: Explicit configuration over implicit behavior
3. **Production Defaults**: Secure and performant defaults out of the box
4. **Extensibility**: Easy to extend without modifying core components
5. **Testability**: Designed for comprehensive testing from day one

### Non-Goals

- Heavy frameworks or opinionated architectures
- Business logic or domain-specific patterns
- Support for Go versions < 1.21
- Hardcoded configuration values

## Architecture Layers

```
┌─────────────────────────────────────────┐
│              HTTP Layer                │
│  ┌─────────────┐  ┌─────────────────┐  │
│  │   Router    │  │   Middleware    │  │
│  │   (chi)     │  │   (custom)     │  │
│  └─────────────┘  └─────────────────┘  │
└─────────────────────────────────────────┘
┌─────────────────────────────────────────┐
│             Application Layer           │
│  ┌─────────────┐  ┌─────────────────┐  │
│  │   Server    │  │    Handlers     │  │
│  │             │  │   (health)       │  │
│  └─────────────┘  └─────────────────┘  │
└─────────────────────────────────────────┘
┌─────────────────────────────────────────┐
│              Core Layer                 │
│  ┌─────────────┐  ┌─────────────────┐  │
│  │   Config    │  │    Logger       │  │
│  │  (viper)    │  │   (zerolog)     │  │
│  └─────────────┘  └─────────────────┘  │
└─────────────────────────────────────────┘
```

## Component Architecture

### Configuration System

**Location**: `internal/config/`

**Technology**: Viper + Struct Validation

**Design Decisions**:
- **Hierarchical Loading**: File → Environment → Defaults
- **Type Safety**: Struct-based configuration with validation
- **Environment Prefix**: `APP_` prefix to avoid conflicts
- **Multiple Sources**: YAML files, environment variables, command-line flags

```go
type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Logger   LoggerConfig   `mapstructure:"logger"`
}
```

**Benefits**:
- Clear separation of concerns
- Type-safe configuration
- Environment-aware deployment
- Validation at startup

### Logging System

**Location**: `internal/logger/`

**Technology**: zerolog (structured JSON logging)

**Design Decisions**:
- **Structured Logging**: JSON format for production, console for development
- **Contextual Methods**: `WithRequestID()`, `WithUserID()` for tracing
- **Specialized Methods**: `HTTPRequest()`, `DatabaseQuery()` for common patterns
- **Performance**: Zero allocation logging paths

```go
type Logger struct {
    zerolog.Logger
}

func (l *Logger) WithRequestID(requestID string) *Logger
func (l *Logger) HTTPRequest(method, path, remoteAddr string, statusCode int, duration int64)
```

**Benefits**:
- High performance structured logging
- Easy integration with log aggregation systems
- Context propagation for distributed tracing
- Consistent log format across services

### HTTP Server

**Location**: `internal/server/`

**Technology**: chi/v5 router + custom middleware

**Design Decisions**:
- **Lightweight Router**: chi for stdlib-compatible routing
- **Middleware Chain**: Request ID, logging, recovery, timeouts
- **Graceful Shutdown**: Context-based shutdown with timeout
- **Health Checks**: `/healthz` and `/ready` endpoints

```go
type Server struct {
    config *config.Config
    logger *logger.Logger
    router *chi.Mux
    server *http.Server
}
```

**Middleware Stack**:
1. Request ID generation
2. Real IP extraction
3. Panic recovery
4. Request logging
5. Request timeout
6. Content type validation

**Benefits**:
- Minimal overhead
- Composable middleware
- Production-ready defaults
- Easy to extend

## Project Structure

### Directory Layout

```
cmd/server/          # Application entry points
├── main.go         # Main application with graceful shutdown

internal/           # Private application code
├── config/         # Configuration management
├── logger/         # Structured logging
└── server/         # HTTP server and handlers
    ├── server.go   # Server implementation
    └── handlers/   # HTTP handlers
        └── health.go

configs/            # Configuration files
├── config.example.yaml
└── config.yaml     # Runtime configuration

pkg/                # Public library code
└── errors/         # Common error types

test/               # Test utilities
└── testutils/      # Test helpers and mocks

scripts/            # Build and deployment scripts
└── build.sh        # Build automation

.github/            # CI/CD workflows
├── workflows/
│   └── ci.yml      # GitHub Actions
└── dependabot.yml  # Dependency updates
```

### Package Organization

**cmd/**: Application entry points following Go conventions
**internal/**: Private code that should not be imported by other projects
**pkg/**: Public library code that can be imported by other projects
**configs/**: Configuration files and examples
**test/**: Test utilities and shared test code

## Dependency Management

### Core Dependencies

| Component | Library | Version | Reason |
|-----------|---------|---------|---------|
| Router | chi/v5 | v5.2.5 | Lightweight, stdlib-compatible |
| Config | viper | v1.16.0 | Multi-source configuration |
| Logging | zerolog | v1.31.0 | High-performance structured logging |
| Testing | testify | v1.8.4 | Assertions and test utilities |

### Development Dependencies

| Tool | Purpose | Version |
|------|---------|---------|
| golangci-lint | Linting aggregator | Latest |
| air | Hot reload for development | Latest |
| go-mod-tidy | Dependency management | Built-in |

## Configuration Management

### Loading Priority

1. Command-line flags (`-config` flag)
2. Environment variables (`APP_*` prefix)
3. Configuration file (`configs/config.yaml`)
4. Default values

### Environment Variables

All configuration can be overridden via environment variables:

```bash
APP_SERVER_HOST=0.0.0.0
APP_SERVER_PORT=8080
APP_DATABASE_HOST=localhost
APP_LOGGER_LEVEL=debug
```

### Validation

Configuration is validated at startup with clear error messages:

```go
if err := validate(config); err != nil {
    return fmt.Errorf("config validation failed: %w", err)
}
```

## Error Handling Strategy

### Error Types

```go
// Wrapped errors for context
return fmt.Errorf("database connection failed: %w", err)

// Structured errors with types
type ValidationError struct {
    Field   string
    Message string
}
```

### Error Propagation

- Use wrapped errors (`fmt.Errorf("%w")`) for context
- Log errors at appropriate levels
- Return user-friendly error messages
- Include correlation IDs for debugging

## Testing Strategy

### Test Organization

```
internal/
├── config/
│   ├── config.go
│   └── config_test.go    # Unit tests
├── logger/
│   ├── logger.go
│   └── logger_test.go    # Unit tests
└── server/
    ├── server.go
    └── server_test.go    # Integration tests
```

### Test Types

1. **Unit Tests**: Individual component testing
2. **Integration Tests**: Component interaction testing
3. **End-to-End Tests**: Full application testing
4. **Performance Tests**: Load and stress testing

### Coverage Goals

- **Internal packages**: >80% coverage
- **Critical paths**: >95% coverage
- **Error handling**: 100% coverage

## Deployment Architecture

### Container Strategy

**Multi-stage Dockerfile**:
1. **Build stage**: Compile Go binary
2. **Runtime stage**: Minimal Alpine image
3. **Security**: Non-root user, minimal attack surface

### Environment Support

- **Development**: Docker Compose with hot reload
- **Staging**: Containerized deployment
- **Production**: Kubernetes or cloud-native deployment

### Health Checks

- **Liveness**: `/healthz` - Basic service health
- **Readiness**: `/ready` - Dependency connectivity
- **Startup**: Container startup validation

## Security Considerations

### Default Security

- **Non-root containers**: Run as non-privileged user
- **Minimal base images**: Reduce attack surface
- **No secrets in code**: Environment-driven configuration
- **Input validation**: Struct validation and sanitization

### Security Headers

```go
// Security middleware (to be added)
func securityMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        next.ServeHTTP(w, r)
    })
}
```

## Performance Considerations

### Logging Performance

- **Zero allocation**: zerolog's design minimizes allocations
- **Structured fields**: Efficient JSON serialization
- **Level filtering**: Early log level filtering

### HTTP Performance

- **Connection pooling**: Default Go HTTP client optimizations
- **Timeouts**: Configurable timeouts for all operations
- **Middleware efficiency**: Minimal overhead middleware chain

### Memory Management

- **Object pooling**: Reuse objects where appropriate
- **Buffer management**: Efficient buffer handling
- **Garbage collection**: Minimize allocations

## Monitoring and Observability

### Logging Strategy

- **Structured logs**: JSON format for machine parsing
- **Correlation IDs**: Request tracing across services
- **Error tracking**: Detailed error context and stack traces

### Metrics (Future)

- **Prometheus integration**: Standard metrics exposition
- **Custom metrics**: Business and application metrics
- **Health metrics**: Service health and performance

### Tracing (Future)

- **OpenTelemetry**: Distributed tracing integration
- **Context propagation**: Request tracing across services
- **Performance profiling**: Request latency analysis

## Extension Points

### Adding New Routes

```go
// In server.go setupRoutes()
s.router.Route("/api/v1", func(r chi.Router) {
    r.Get("/users", s.handleListUsers)
    r.Post("/users", s.handleCreateUser)
})
```

### Adding New Middleware

```go
// In server.go setupMiddleware()
s.router.Use(middleware.Compress(5))
s.router.Use(securityMiddleware)
```

### Adding New Configuration

```go
// In config/config.go
type Config struct {
    // existing fields...
    Redis RedisConfig `mapstructure:"redis"`
}

type RedisConfig struct {
    Host string `mapstructure:"host"`
    Port int    `mapstructure:"port"`
}
```

## Migration Guide

### From Standard Library HTTP

1. Replace `http.ServeMux` with `chi.NewRouter()`
2. Add middleware chain for logging and recovery
3. Implement graceful shutdown with context
4. Add structured logging throughout

### From Other Frameworks

1. Extract business logic from framework code
2. Replace framework-specific routing with chi
3. Migrate configuration to Viper
4. Replace logging with zerolog

## Best Practices

### Code Organization

- **Package by feature**: Group related functionality
- **Interface segregation**: Small, focused interfaces
- **Dependency injection**: Constructor-based injection
- **Explicit dependencies**: No global state

### Error Handling

- **Early returns**: Handle errors immediately
- **Context preservation**: Include context in errors
- **User messages**: Separate internal and external errors
- **Logging**: Log errors at appropriate levels

### Performance

- **Avoid allocations**: Reuse buffers and objects
- **Connection reuse**: Use connection pooling
- **Timeouts**: Set appropriate timeouts for all operations
- **Profiling**: Regular performance profiling

## Future Enhancements

### Planned Features

1. **Metrics Integration**: Prometheus metrics exposition
2. **Tracing**: OpenTelemetry distributed tracing
3. **Circuit Breaker**: Resilience patterns
4. **Rate Limiting**: Request rate limiting
5. **Authentication**: JWT-based authentication
6. **Database Integration**: Connection pooling and migrations

### Architecture Evolution

- **Microservices**: Service decomposition patterns
- **Event-driven**: Message queue integration
- **CQRS**: Command query responsibility segregation
- **Event Sourcing**: Event-based persistence

## Conclusion

Go Starter provides a solid foundation for building production-ready Go applications. Its architecture emphasizes simplicity, performance, and maintainability while following Go best practices and idioms.

The template is designed to be:

- **Minimal**: Only essential dependencies and patterns
- **Extensible**: Easy to add new features without breaking existing code
- **Production-ready**: Security, performance, and observability built-in
- **Maintainable**: Clear structure and comprehensive testing

This architecture enables teams to focus on business logic while having a solid technical foundation for their applications.
