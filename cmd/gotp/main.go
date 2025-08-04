package main

import (
	"fmt"
	"os"

	"go.coldcutz.net/gotp/pkg/erlang"
)

func main() {
	if err := run(); err != nil {
		fmt.Println("ERROR: ", err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Println("CGO Skeleton Example")

	config := erlang.Config{
		Cookie:         "super_secret",
		MyNodeName:     "itest",
		RemoteNodeName: "itestapp@localhost",
		Creation:       1,
		Port:           9999,
	}

	conn, err := erlang.NewConnection(config)
	if err != nil {
		return fmt.Errorf("failed to create connection: %w", err)
	}

	_, err = conn.Publish(config.Port)
	if err != nil {
		return fmt.Errorf("failed to publish port: %w", err)
	}
	fmt.Println("Published port 9999 to epmd")

	sockFd, err := conn.Connect(config.RemoteNodeName)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	fmt.Println("Connected to remote Erlang node")

	response, err := conn.SendRPC(sockFd, "erlang", "node", "[]")
	if err != nil {
		return fmt.Errorf("failed to send RPC: %w", err)
	}

	fmt.Println("Response: ", response)
	fmt.Println("Sent message to remote Erlang node")

	return nil
}