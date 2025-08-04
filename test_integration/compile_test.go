package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestGoCodeCompilation(t *testing.T) {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Path to the main project directory
	projectDir := filepath.Join(cwd, "..")

	// Test that the Go code compiles
	t.Log("Testing Go code compilation...")
	
	// Set CGO environment variables for compilation test
	env := os.Environ()
	env = append(env, "CGO_ENABLED=1")
	env = append(env, "CGO_CFLAGS=-I/usr/lib/erlang/lib/erl_interface-5.5.2/include -Wall -g")
	env = append(env, "CGO_LDFLAGS=-L/usr/lib/erlang/lib/erl_interface-5.5.2/lib -lei -lpthread")

	compileCmd := exec.Command("go", "build", "./cmd/gotp")
	compileCmd.Dir = projectDir
	compileCmd.Env = env
	
	output, err := compileCmd.CombinedOutput()
	if err != nil {
		// Check if it's a missing Erlang library issue
		if strings.Contains(string(output), "ei.h: No such file or directory") {
			t.Skipf("Skipping compilation test - Erlang ei library not found. This is expected in test environments without Erlang development headers.\nOutput: %s", string(output))
		}
		t.Fatalf("Failed to compile Go code: %v, output: %s", err, string(output))
	}

	t.Log("Go code compiles successfully")
}

func TestLibraryPackageCompilation(t *testing.T) {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Path to the main project directory
	projectDir := filepath.Join(cwd, "..")

	// Test that the library package compiles
	t.Log("Testing library package compilation...")
	
	// Set CGO environment variables for compilation test
	env := os.Environ()
	env = append(env, "CGO_ENABLED=1")
	env = append(env, "CGO_CFLAGS=-I/usr/lib/erlang/lib/erl_interface-5.5.2/include -Wall -g")
	env = append(env, "CGO_LDFLAGS=-L/usr/lib/erlang/lib/erl_interface-5.5.2/lib -lei -lpthread")

	compileCmd := exec.Command("go", "build", "./pkg/erlang")
	compileCmd.Dir = projectDir
	compileCmd.Env = env
	
	output, err := compileCmd.CombinedOutput()
	if err != nil {
		// Check if it's a missing Erlang library issue
		if strings.Contains(string(output), "ei.h: No such file or directory") {
			t.Skipf("Skipping library compilation test - Erlang ei library not found. This is expected in test environments without Erlang development headers.\nOutput: %s", string(output))
		}
		t.Fatalf("Failed to compile library package: %v, output: %s", err, string(output))
	}

	t.Log("Library package compiles successfully")
}