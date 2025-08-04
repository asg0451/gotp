package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestIntegrationWithElixirApp(t *testing.T) {
	// Check if required tools are available
	if err := checkRequiredTools(); err != nil {
		t.Skipf("Skipping integration test - %v", err)
	}

	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Path to the Elixir app (go up one level from test_integration)
	elixirAppPath := filepath.Join(cwd, "..", "itest", "itest_elixir_app")

	// Start epmd (Erlang Port Mapper Daemon) if not already running
	startEpmd(t)

	// Create a context with timeout for the entire test
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Start the Elixir application
	elixirCmd := exec.CommandContext(ctx, "iex", "--sname", "itestapp@localhost", "--cookie", "super_secret", "-S", "mix", "run")
	elixirCmd.Dir = elixirAppPath
	
	// Capture output for debugging
	var elixirOutput strings.Builder
	elixirCmd.Stdout = &elixirOutput
	elixirCmd.Stderr = &elixirOutput

	t.Log("Starting Elixir application...")
	if err := elixirCmd.Start(); err != nil {
		t.Fatalf("Failed to start Elixir application: %v", err)
	}

	// Give the Elixir app time to start up
	time.Sleep(5 * time.Second)

	// Ensure the process is cleaned up
	defer func() {
		t.Log("Cleaning up Elixir process...")
		if elixirCmd.Process != nil {
			// Try graceful shutdown first
			if err := elixirCmd.Process.Signal(os.Interrupt); err != nil {
				t.Logf("Failed to send interrupt signal: %v", err)
			}
			
			// Wait a bit for graceful shutdown
			done := make(chan error, 1)
			go func() {
				done <- elixirCmd.Wait()
			}()
			
			select {
			case <-done:
				t.Log("Elixir process terminated gracefully")
			case <-time.After(5 * time.Second):
				t.Log("Force killing Elixir process...")
				if err := elixirCmd.Process.Kill(); err != nil {
					t.Logf("Failed to kill process: %v", err)
				}
				elixirCmd.Wait() // Wait for the process to actually terminate
			}
		}
		t.Logf("Elixir app output: %s", elixirOutput.String())
	}()

	// Test the connection by running our Go program
	t.Log("Testing Go-Elixir communication...")
	goCmd := exec.CommandContext(ctx, "go", "run", "./cmd/gotp")
	goCmd.Dir = filepath.Join(cwd, "..") // Run from parent directory
	
	// Capture Go program output
	var goOutput strings.Builder
	goCmd.Stdout = &goOutput
	goCmd.Stderr = &goOutput

	if err := goCmd.Run(); err != nil {
		// Check if it's a CGO/Erlang library issue
		if strings.Contains(goOutput.String(), "ei.h: No such file or directory") {
			t.Skipf("Skipping integration test - Erlang ei library not found. This is expected in test environments without Erlang development headers.\nOutput: %s", goOutput.String())
		}
		t.Fatalf("Go program failed: %v\nOutput: %s", err, goOutput.String())
	}

	// Check for expected output from Go program
	output := goOutput.String()
	if !strings.Contains(output, "Connected to remote Erlang node") {
		t.Errorf("Expected 'Connected to remote Erlang node' in output, got: %s", output)
	}
	if !strings.Contains(output, "Sent message to remote Erlang node") {
		t.Errorf("Expected 'Sent message to remote Erlang node' in output, got: %s", output)
	}

	// Check for message receipt in Elixir process stdout
	elixirOutputStr := elixirOutput.String()
	if !strings.Contains(elixirOutputStr, "Received message: \"Hello world\"") {
		t.Errorf("Expected Elixir process to receive message 'Hello world', but got: %s", elixirOutputStr)
	}

	t.Log("Integration test completed successfully")
	t.Logf("Go program output: %s", output)
	t.Logf("Elixir process output: %s", elixirOutputStr)
}



func startEpmd(t *testing.T) {
	// Check if epmd is already running
	checkCmd := exec.Command("epmd", "-names")
	if err := checkCmd.Run(); err == nil {
		t.Log("epmd is already running")
		return
	}

	// Start epmd in daemon mode
	t.Log("Starting epmd...")
	epmdCmd := exec.Command("epmd", "-d")
	if err := epmdCmd.Run(); err != nil {
		t.Fatalf("Failed to start epmd: %v", err)
	}

	// Give epmd time to start
	time.Sleep(2 * time.Second)
}



func TestElixirAppCompilation(t *testing.T) {
	// Check if required tools are available
	if err := checkRequiredTools(); err != nil {
		t.Skipf("Skipping compilation test - %v", err)
	}

	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Path to the Elixir app
	elixirAppPath := filepath.Join(cwd, "..", "itest", "itest_elixir_app")

	// Test that the Elixir app can be compiled
	t.Log("Testing Elixir app compilation...")
	compileCmd := exec.Command("mix", "compile")
	compileCmd.Dir = elixirAppPath
	
	output, err := compileCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to compile Elixir app: %v, output: %s", err, string(output))
	}

	t.Log("Elixir app compiles successfully")
	t.Logf("Compilation output: %s", string(output))
}

func checkRequiredTools() error {
	tools := []string{"go", "mix", "iex", "elixir", "epmd"}
	
	for _, tool := range tools {
		if _, err := exec.LookPath(tool); err != nil {
			return fmt.Errorf("required tool '%s' not found in PATH", tool)
		}
	}
	
	return nil
}