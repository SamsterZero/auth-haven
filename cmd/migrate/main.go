package main

import (
	"auth-haven/internal/config"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	m, err := migrate.New("file://migrations", cfg.DBUrl)
	if err != nil {
		log.Fatalf("failed to init migrate: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration failed: %v", err)
	}

	log.Println("✅ database migrated successfully")
}
