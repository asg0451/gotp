# CGO Project Makefile

# Variables
LIB_NAME = gotp
STATIC_LIB = lib$(LIB_NAME).a
SHARED_LIB = lib$(LIB_NAME).so
C_FILES = gotp.c
HEADER_FILES = gotp.h
GO_FILES = main.go
BINARY = gotp

# Compiler flags
CC = gcc
CFLAGS = -Wall -Wextra -std=c99 -fPIC -g
LDFLAGS = 

# Default target
all: $(STATIC_LIB) $(BINARY)

# Build static library
$(STATIC_LIB): $(C_FILES) $(HEADER_FILES)
	$(CC) $(CFLAGS) -c $(C_FILES)
	ar rcs $(STATIC_LIB) *.o

# Build shared library (optional)
$(SHARED_LIB): $(C_FILES) $(HEADER_FILES)
	$(CC) $(CFLAGS) -shared -o $(SHARED_LIB) $(C_FILES)

# Build Go binary
$(BINARY): $(GO_FILES) $(STATIC_LIB)
	go build -o $(BINARY) $(GO_FILES)

# Run the program
run: $(BINARY)
	./$(BINARY)

# Test the C library standalone
test-c: $(STATIC_LIB)
	$(CC) $(CFLAGS) -DTEST_MAIN -o test_c $(C_FILES) -L. -l$(LIB_NAME)
	./test_c

# Clean build artifacts
clean:
	rm -f *.o $(STATIC_LIB) $(SHARED_LIB) $(BINARY) test_c

# Install system dependencies (example for common libs)
deps:
	@echo "Add system library installation commands here"
	@echo "Example: sudo apt-get install libssl-dev"

# Format Go code
fmt:
	go fmt ./...

# Run Go tests
test:
	go test ./...

# Development build with debug info
debug: CFLAGS += -DDEBUG -O0
debug: $(BINARY)

# Release build with optimizations
release: CFLAGS += -O2 -DNDEBUG
release: $(BINARY)

.PHONY: all clean run test-c deps fmt test debug release