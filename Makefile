# Development build with CGO enabled (Linux)
dev:
	CGO_ENABLED=1 CGO_CFLAGS="-I/usr/lib/erlang/lib/erl_interface-5.5.2/include -Wall -g" \
	CGO_LDFLAGS="-L/usr/lib/erlang/lib/erl_interface-5.5.2/lib -lei -lpthread" \
	go build -o gotp ./cmd/gotp

# Development build with CGO enabled (macOS Homebrew)
dev-macos:
	CGO_ENABLED=1 CGO_CFLAGS="-I/opt/homebrew/Cellar/erlang/28.0.1/lib/erlang/lib/erl_interface-5.6/include -Wall -g" \
	CGO_LDFLAGS="-L/opt/homebrew/Cellar/erlang/28.0.1/lib/erlang/lib/erl_interface-5.6/lib -L. -lei -lpthread" \
	go build -o gotp ./cmd/gotp

# Auto-detect platform for development build
dev-auto:
	@if [ "$$(uname)" = "Darwin" ]; then \
		echo "Detected macOS, using Homebrew paths"; \
		$(MAKE) dev-macos; \
	else \
		echo "Detected Linux, using system paths"; \
		$(MAKE) dev; \
	fi

# CLI build for production (Linux)
cli:
	CGO_ENABLED=1 CGO_CFLAGS="-I/usr/lib/erlang/lib/erl_interface-5.5.2/include -Wall -O2" \
	CGO_LDFLAGS="-L/usr/lib/erlang/lib/erl_interface-5.5.2/lib -lei -lpthread" \
	go build -ldflags="-s -w" -o gotp ./cmd/gotp

# CLI build for production (macOS Homebrew)
cli-macos:
	CGO_ENABLED=1 CGO_CFLAGS="-I/opt/homebrew/Cellar/erlang/28.0.1/lib/erlang/lib/erl_interface-5.6/include -Wall -O2" \
	CGO_LDFLAGS="-L/opt/homebrew/Cellar/erlang/28.0.1/lib/erlang/lib/erl_interface-5.6/lib -L. -lei -lpthread" \
	go build -ldflags="-s -w" -o gotp ./cmd/gotp

# Auto-detect platform for production build
cli-auto:
	@if [ "$$(uname)" = "Darwin" ]; then \
		echo "Detected macOS, using Homebrew paths"; \
		$(MAKE) cli-macos; \
	else \
		echo "Detected Linux, using system paths"; \
		$(MAKE) cli; \
	fi

# Clean build artifacts
clean:
	rm -f gotp

# Run tests (Linux)
test:
	CGO_ENABLED=1 CGO_CFLAGS="-I/usr/lib/erlang/lib/erl_interface-5.5.2/include -Wall -g" \
	CGO_LDFLAGS="-L/usr/lib/erlang/lib/erl_interface-5.5.2/lib -lei -lpthread" \
	go test ./...

# Run tests (macOS Homebrew)
test-macos:
	CGO_ENABLED=1 CGO_CFLAGS="-I/opt/homebrew/Cellar/erlang/28.0.1/lib/erlang/lib/erl_interface-5.6/include -Wall -g" \
	CGO_LDFLAGS="-L/opt/homebrew/Cellar/erlang/28.0.1/lib/erlang/lib/erl_interface-5.6/lib -L. -lei -lpthread" \
	go test ./...

# Auto-detect platform for tests
test-auto:
	@if [ "$$(uname)" = "Darwin" ]; then \
		echo "Detected macOS, using Homebrew paths"; \
		$(MAKE) test-macos; \
	else \
		echo "Detected Linux, using system paths"; \
		$(MAKE) test; \
	fi

# Docker build
docker-build:
	docker build -t gotp .

# Docker run
docker-run:
	docker run --rm -it gotp

# Start itest app
start-itest-app:
	cd itest/itest_elixir_app && iex --sname itestapp@localhost --cookie 'super_secret' -S mix run

# Run itest example
itest-example:
	elixir --sname itest@localhost --cookie 'super_secret' -e 'Node.connect(:"itestapp@localhost");  Node.spawn(:"itestapp@localhost", fn -> send(ItestElixirApp.Worker, "hi") end)'

# Run epmd
itest-run-epmd:
	epmd -d

# Automated integration tests
itest-automated:
	cd test_integration && go test -v -run TestIntegrationWithElixirApp

itest-all:
	cd test_integration && go test -v

itest-compile:
	cd test_integration && go test -v -run TestElixirAppCompilation

# Show help
help:
	@echo "Available targets:"
	@echo "  dev          - Development build (Linux)"
	@echo "  dev-macos    - Development build (macOS Homebrew)"
	@echo "  dev-auto     - Auto-detect platform for development build"
	@echo "  cli          - Production build (Linux)"
	@echo "  cli-macos    - Production build (macOS Homebrew)"
	@echo "  cli-auto     - Auto-detect platform for production build"
	@echo "  test         - Run tests (Linux)"
	@echo "  test-macos   - Run tests (macOS Homebrew)"
	@echo "  test-auto    - Auto-detect platform for tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  start-itest-app - Start itest Elixir app"
	@echo "  itest-example   - Run itest example"
	@echo "  itest-run-epmd  - Run epmd"
	@echo "  itest-automated - Run automated integration tests"
	@echo "  itest-all       - Run all integration tests"
	@echo "  itest-compile   - Test Elixir app compilation"
	@echo ""
	@echo "For local development on macOS, use: make dev-macos"
	@echo "For local development on Linux, use: make dev"
	@echo "For auto-detection, use: make dev-auto"
