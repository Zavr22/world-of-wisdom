package main

import (
	"context"
	"fmt"
	"github.com/Zavr22/world-of-wisdom/internal/client"
	"github.com/Zavr22/world-of-wisdom/internal/pkg/config"
	"log"
	"net"
	"strconv"
)

// main function is the entry point for the client application
func main() {
	fmt.Println("Starting client...")

	// Load the configuration from file.
	cfg, err := config.Load("config/config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}

	// Create a new context with the configuration value.
	ctx := context.WithValue(context.Background(), "config", cfg)

	// Generate the server address.
	address := net.JoinHostPort(cfg.ServerHost, strconv.Itoa(cfg.ServerPort))

	// Run the client.
	if err := client.Run(ctx, address); err != nil {
		log.Fatalf("Failed to run client: %s", err)
	}

	fmt.Println("Client exited.")
}
