# Go Starter

[![Doctor Report](https://img.shields.io/badge/Doctor%20Report-🩺-green)](https://github.com/Build-with-Go/go-starter/actions)
[![CI](https://github.com/Build-with-Go/go-starter/workflows/CI/badge.svg)](https://github.com/Build-with-Go/go-starter/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)

A production-ready, batteries-included Go project template for the Build-with-Go organization.

## 🚀 Quick Start

```bash
# Clone the repository
git clone https://github.com/Build-with-Go/go-starter.git
cd go-starter

# Copy configuration
cp configs/config.example.yaml configs/config.yaml

# Install dependencies
make deps

# Run the application
make run
```

Or run directly:

```bash
go run cmd/server/main.go
```

## 📋 Features

- ✅ **Go 1.21+** with idiomatic patterns
- ✅ **Configuration** with Viper (file + env + validation)
- ✅ **Structured Logging** with zerolog (JSON output)
- ✅ **HTTP Router** with chi/v5
- ✅ **Health Checks** for monitoring
- ✅ **Graceful Shutdown** with context handling
- ✅ **Testing** with testify utilities
- ✅ **Linting** with golangci-lint
- ✅ **Docker** multi-stage builds
- ✅ **CI/CD** with GitHub Actions

## 🏗️ Project Structure

```
cmd/server/          # Application entry point
internal/            # Private application code
  ├── config/        # Configuration management
  ├── logger/        # Structured logging
  └── server/        # HTTP server and handlers
configs/             # Configuration files
pkg/                 # Public library code
test/                # Test utilities
scripts/             # Build and deployment scripts
.github/             # CI/CD workflows
```

## ⚙️ Configuration

The application loads configuration from:
1. `configs/config.yaml` (or custom path with `-config` flag)
2. Environment variables with `APP_` prefix
3. Default values

Example environment variables:
```bash
export APP_SERVER_PORT=8080
export APP_DATABASE_HOST=localhost
export APP_LOGGER_LEVEL=debug
```

## 📊 Code Quality

This project maintains high code quality standards:

- **Doctor Report**: Comprehensive code health analysis
- **golangci-lint**: Comprehensive static analysis
- **Test Coverage**: >80% coverage on core packages
- **Documentation**: Full API documentation and examples

### Doctor Report

The Doctor Report provides comprehensive code health analysis, running the same checks that Go Report Card would perform:

```bash
# Doctor Report runs automatically on:
- Every push to main/master
- Pull requests
- Daily at 00:00 UTC
```

The report checks for:
- ✅ Go formatting (gofmt)
- ✅ Go vet analysis
- ✅ Cyclomatic complexity (gocyclo)
- ✅ Code duplication (dupl)
- ✅ Package documentation
- ✅ Test coverage analysis
- ✅ Import organization
- ✅ Build verification
- ✅ And more...

## ��️ Development

### Commands

```bash
make help          # Show all available commands
make run           # Run the application
make test          # Run tests
make lint          # Run linter
make doctorreport  # Run Doctor Report analysis
make build         # Build binary
make clean         # Clean artifacts
```

### Hot Reload

```bash
make dev-setup     # Install development tools
make dev           # Run with hot reload
```

### Testing

```bash
make test          # Run all tests
make test-coverage # Generate coverage report
```

## 🐳 Docker

```bash
# Build image
make docker-build

# Run container
make docker-run
```

## 📊 Health Checks

- `GET /healthz` - Basic health check
- `GET /ready` - Readiness probe (database connectivity)

## 📝 Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_SERVER_HOST` | `localhost` | Server bind address |
| `APP_SERVER_PORT` | `8080` | Server port |
| `APP_DATABASE_HOST` | `localhost` | Database host |
| `APP_DATABASE_PORT` | `5432` | Database port |
| `APP_DATABASE_NAME` | `go_starter` | Database name |
| `APP_DATABASE_USER` | `postgres` | Database user |
| `APP_DATABASE_PASSWORD` | `password` | Database password |
| `APP_LOGGER_LEVEL` | `info` | Log level (trace, debug, info, warn, error) |
| `APP_LOGGER_FORMAT` | `json` | Log format (json, console) |

## 🧪 Testing

The project uses table-driven tests and aims for >80% coverage:

```bash
# Run specific package tests
go test -v ./internal/config

# Run with race detection
go test -race ./...

# Generate coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## 🔧 Linting

Uses [golangci-lint](https://golangci-lint.run/) with comprehensive rules:

```bash
make lint  # Run linter
```

## 📦 Dependencies

- [chi/v5](https://github.com/go-chi/chi) - HTTP router
- [viper](https://github.com/spf13/viper) - Configuration management
- [zerolog](https://github.com/rs/zerolog) - Structured logging
- [testify](https://github.com/stretchr/testify) - Testing utilities

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📚 Architecture

See [ARCHITECTURE.md](ARCHITECTURE.md) for detailed design decisions and patterns.
