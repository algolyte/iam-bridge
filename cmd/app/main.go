package main

import (
	"log"
	"os"

	"github.com/zahidhasanpapon/iam-bridge/internal/server"
)

func main() {
	// Initialize the server
	srv, err := server.NewServer()
	if err != nil {
		log.Printf("Error initializing server: %v\n", err)
		os.Exit(1)
	}

	// Start the server
	if err := srv.Start(); err != nil {
		log.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}
