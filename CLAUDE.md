# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

GOTP is a Go project that implements Erlang/OTP (Open Telecom Platform) communication using CGO bindings to the Erlang Interface library (`ei`). The project allows Go applications to connect to and communicate with Erlang/Elixir nodes as if they were part of an Erlang distributed system.

## Key Architecture

The project consists of two main components:

1. **Go C Interface** (`main.go`, `gotp.h`):
   - Uses CGO to interface with Erlang's `ei` library
   - Implements Erlang distribution protocol communication
   - Provides wrapper functions for variadic C functions that CGO cannot handle directly
   - Connects to remote Erlang nodes and sends/receives messages

2. **Elixir Test Application** (`itest/itest_elixir_app/`):
   - Simple Elixir OTP application for integration testing
   - Contains a `Worker` GenServer that can receive messages from the Go client
   - Runs as a named node (`itestapp@localhost`) for distributed communication

## Development Commands

### Build and Run
```bash
# Build and run the Go binary (requires Erlang/OTP development headers)
go run .
```

### Integration Testing
```bash
# Start the Elixir test application (run in separate terminal)
make start-itest-app

# Start EPMD (Erlang Port Mapper Daemon) if not running
make itest-run-epmd

# Run a simple Elixir integration test
make itest-example

# Run the Go application (communicates with the Elixir app)
go run .
```

## Critical Dependencies

- **Erlang/OTP**: Required for `ei` library headers and runtime
- **EPMD**: Erlang Port Mapper Daemon must be running for node discovery
- **CGO**: Enabled for C library integration

## Important Files

- `main.go`: Main Go application with Erlang communication logic
- `gotp.h`: C wrapper functions for variadic `ei` library functions
- `go.mod`: Go module definition (minimal, Go 1.24.2+)
- `Makefile`: Integration test commands for Elixir interop
- `itest/itest_elixir_app/`: Complete Elixir OTP application for testing

## Communication Protocol

The project uses Erlang's external term format and distribution protocol:
- Connects to remote nodes using `ei_connect`
- Publishes to EPMD using `ei_publish`
- Sends messages via RPC calls to named processes
- Handles tuple decoding for complex Erlang terms

## Node Configuration

- Go node name: `itest`
- Elixir test node: `itestapp@localhost`
- Shared cookie: `super_secret` (for authentication)
- Default port: 9999 (published to EPMD)
