# Automated Integration Test

This directory now contains an automated Go integration test that replaces the manual integration testing process.

## What the Test Does

The automated integration test (`test_integration/integration_test.go`) performs the following steps:

1. **Prerequisites Check**: Verifies that all required tools are available (Go, Mix, IEx, Elixir, EPMD)
2. **EPMD Setup**: Starts the Erlang Port Mapper Daemon if not already running
3. **Elixir App Startup**: Launches the Elixir application with the correct node name and cookie
4. **Go-Elixir Communication**: Runs the Go program to test communication with the Elixir node
5. **Verification**: Checks that the expected output messages are present
6. **Cleanup**: Properly terminates all processes

## Test Functions

### `TestIntegrationWithElixirApp`
The main integration test that:
- Starts the Elixir application
- Runs the Go program to connect and send messages
- Verifies successful communication
- Captures and logs all output for debugging

### `TestElixirNodeConnection`
Tests direct Elixir-to-Elixir communication to verify the Elixir app is working correctly.

### `TestBuildAndRun`
Verifies that the Go program can be built successfully.

### `TestElixirAppCompilation`
Verifies that the Elixir application can be compiled successfully.

## Running the Tests

### Using Make (Recommended)
```bash
# Run the main integration test
make itest-automated

# Run all tests
make itest-all

# Run build test only
make itest-build

# Run compilation test only
make itest-compile
```

### Using Go Test Directly
```bash
# Run the main integration test
cd test_integration && go test -v -run TestIntegrationWithElixirApp

# Run all tests
cd test_integration && go test -v

# Run a specific test
cd test_integration && go test -v -run TestElixirNodeConnection
```

## Prerequisites

Before running the tests, ensure you have:

1. **Go** installed and in your PATH
2. **Elixir** and **Erlang** installed
3. **Mix** (comes with Elixir)
4. **IEx** (comes with Elixir)
5. **EPMD** (comes with Erlang)

The test will automatically check for these tools and fail with a clear error message if any are missing.

## Manual vs Automated Testing

### Before (Manual)
```bash
# Terminal 1: Start epmd
make itest-run-epmd

# Terminal 2: Start Elixir app
make start-itest-app

# Terminal 3: Test Elixir communication
make itest-example

# Terminal 4: Run Go program
go run main.go
```

### After (Automated)
```bash
# Single command runs everything
make itest-automated
```

## Benefits of Automation

1. **Reproducible**: Same steps every time
2. **Faster**: No manual coordination between terminals
3. **Reliable**: Automatic cleanup and error handling
4. **CI/CD Ready**: Can be run in automated environments
5. **Debugging**: Captures all output for troubleshooting
6. **Validation**: Checks for expected output and behavior

## Troubleshooting

### Common Issues

1. **"Required tool not found"**: Install the missing tool (Go, Elixir, etc.)
2. **"Failed to start Elixir application"**: Check that Mix and IEx are available
3. **"Go program failed"**: Check the captured output for CGO compilation errors
4. **"Expected output not found"**: The integration may have failed - check logs

### Debug Mode

To see more detailed output, run with verbose logging:
```bash
go test -v -run TestIntegrationWithElixirApp
```

The test captures and displays all output from both the Elixir app and Go program, making it easy to diagnose issues.

## Integration with CI/CD

The automated test can be easily integrated into CI/CD pipelines:

```yaml
# Example GitHub Actions step
- name: Run Integration Tests
  run: |
    make itest-automated
```

The test will fail fast if prerequisites are missing and provide clear error messages for debugging.