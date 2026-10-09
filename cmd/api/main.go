package main

import (
	"log"

	"github.com/yoyojr404/internal/bootstrap"
)

func main() {
	// Initialize and wire all dependencies via the bootstrap loader
	app, err := bootstrap.NewApp()
	if err != nil {
		log.Fatalf("[FATAL] Application bootstrap failed: %v", err)
	}

	// Start the application server with graceful shutdown handling
	if err := app.Run(); err != nil {
		log.Fatalf("[FATAL] Application error: %v", err)
	}
}
