# GOTP - Go Erlang/OTP Interface

A Go project that implements Erlang/OTP (Open Telecom Platform) communication using CGO bindings to the Erlang Interface library (`ei`). The project allows Go applications to connect to and communicate with Erlang/Elixir nodes as if they were part of an Erlang distributed system.

## Features

- Uses CGO to interface with Erlang's `ei` library
- Provides wrapper functions for variadic C functions that CGO cannot handle directly
- Supports both Linux and macOS (Homebrew) environments
- Docker support for containerized builds and deployment

## Prerequisites

### Linux
```bash
sudo apt-get install erlang-dev
```

### macOS (Homebrew)
```bash
brew install erlang
```

## Build Options

The project supports multiple build targets for different environments:

### Development Builds
- `make dev` - Linux development build
- `make dev-macos` - macOS Homebrew development build  
- `make dev-auto` - Auto-detect platform for development build

### Production Builds
- `make cli` - Linux production build (optimized)
- `make cli-macos` - macOS Homebrew production build (optimized)
- `make cli-auto` - Auto-detect platform for production build

### Testing
- `make test` - Run tests (Linux)
- `make test-macos` - Run tests (macOS Homebrew)
- `make test-auto` - Auto-detect platform for tests

### Docker
- `make docker-build` - Build Docker image
- `make docker-run` - Run Docker container

### Utilities
- `make clean` - Clean build artifacts
- `make help` - Show all available targets

## Quick Start

### Local Development

**Linux:**
```bash
make dev
./gotp
```

**macOS:**
```bash
make dev-macos
./gotp
```

**Auto-detect:**
```bash
make dev-auto
./gotp
```

### Docker

```bash
make docker-build
make docker-run
```

## CGO Configuration

The project uses environment-based CGO configuration instead of hardcoded paths:

- **Linux**: Uses system Erlang installation paths
- **macOS**: Uses Homebrew Erlang installation paths  
- **Docker**: Uses Debian package paths

This approach provides better cross-platform compatibility and eliminates the need to modify source code for different environments.

## Testing with Erlang/Elixir

The project includes integration tests with an Elixir application:

```bash
# Start the test Elixir app
make start-itest-app

# In another terminal, run the Go application
make dev-auto
./gotp
```

## Project Structure

```
.
├── main.go          # Main Go application with CGO bindings
├── gotp.h           # C header with wrapper functions
├── Makefile         # Build targets for different environments
├── Dockerfile       # Multi-stage Docker build
├── go.mod           # Go module definition
└── itest/           # Integration test Elixir application
```

## Troubleshooting

### Build Issues

1. **Missing Erlang headers**: Install `erlang-dev` package
2. **Wrong paths**: Use the appropriate target for your platform (`dev-macos` for macOS, `dev` for Linux)
3. **CGO not enabled**: Ensure `CGO_ENABLED=1` is set

### Runtime Issues

1. **Connection refused**: Ensure an Erlang node is running and accessible
2. **Library not found**: Verify Erlang runtime libraries are installed

## License

[Add your license information here]