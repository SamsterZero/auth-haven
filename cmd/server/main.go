package main

import (
	"auth-haven/internal/config"
	"auth-haven/internal/db"
	"auth-haven/internal/server"
	"log"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Connect to DB
	conn, err := db.Connect(cfg.DBUrl)
	if err != nil {
		log.Fatalf("failed to connect DB: %v", err)
	}
	defer conn.Close()

	// Start gRPC server
	if err := server.StartGRPC(cfg, conn); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
