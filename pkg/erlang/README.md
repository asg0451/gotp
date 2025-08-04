# Erlang Library Package

This package provides a Go interface to Erlang's ei (Erlang Interface) library, allowing Go applications to connect to and communicate with Erlang nodes.

## Features

- Connect to remote Erlang nodes
- Publish ports to epmd (Erlang Port Mapper Daemon)
- Send RPC calls to remote processes
- Decode Erlang terms and tuples

## Usage

```go
package main

import (
    "fmt"
    "log"
    
    "go.coldcutz.net/gotp/pkg/erlang"
)

func main() {
    config := erlang.Config{
        Cookie:         "super_secret",
        MyNodeName:     "myapp",
        RemoteNodeName: "remote@localhost",
        Creation:       1,
        Port:           9999,
    }

    conn, err := erlang.NewConnection(config)
    if err != nil {
        log.Fatal(err)
    }

    // Publish a port to epmd
    _, err = conn.Publish(config.Port)
    if err != nil {
        log.Fatal(err)
    }

    // Connect to remote node
    sockFd, err := conn.Connect(config.RemoteNodeName)
    if err != nil {
        log.Fatal(err)
    }

    // Send RPC call
    response, err := conn.SendRPC(sockFd, "MyModule", "my_function", "[~s]", "Hello")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Response:", response)
}
```

## Configuration

The `Config` struct contains all necessary connection parameters:

- `Cookie`: The Erlang cookie for authentication
- `MyNodeName`: The name of the local node
- `RemoteNodeName`: The name of the remote node to connect to
- `Creation`: The creation number for the node
- `Port`: The port to publish to epmd

## Dependencies

This package requires:
- Erlang/OTP with ei library
- CGO enabled
- Proper CGO flags for ei library paths

See the main Makefile for examples of how to set up the build environment.