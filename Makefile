# Development build with CGO enabled
dev:
	CGO_ENABLED=1 CGO_CFLAGS="-I/usr/lib/erlang/lib/erl_interface-5.5.2/include -Wall -g" \
	CGO_LDFLAGS="-L/usr/lib/erlang/lib/erl_interface-5.5.2/lib -lei -lpthread" \
	go build -o gotp .

# CLI build for production
cli:
	CGO_ENABLED=1 CGO_CFLAGS="-I/usr/lib/erlang/lib/erl_interface-5.5.2/include -Wall -O2" \
	CGO_LDFLAGS="-L/usr/lib/erlang/lib/erl_interface-5.5.2/lib -lei -lpthread" \
	go build -ldflags="-s -w" -o gotp .

# Clean build artifacts
clean:
	rm -f gotp

# Run tests
test:
	CGO_ENABLED=1 CGO_CFLAGS="-I/usr/lib/erlang/lib/erl_interface-5.5.2/include -Wall -g" \
	CGO_LDFLAGS="-L/usr/lib/erlang/lib/erl_interface-5.5.2/lib -lei -lpthread" \
	go test ./...

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
