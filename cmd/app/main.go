package main

import (
	"context"
	"log"
	"memory-db/internal/bootstrap"
)

func main() {
	app, err := bootstrap.NewApp()
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
		return
	}

	if err = app.RunTCPServer(context.Background()); err != nil {
		log.Fatalf("server exited with error: %v", err)
	}
}
