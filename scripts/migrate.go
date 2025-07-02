package main

import (
	"log"
	"os"
	"tj_techtest/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// Connect to database
	config.ConnectDB()

	// Get the underlying sql.DB from GORM
	sqlDB, err := config.DB.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	// Create migrate instance
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		log.Fatal("Failed to create postgres driver:", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal("Failed to create migrate instance:", err)
	}

	// Check command line arguments
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run scripts/migrate.go [up|down|version]")
	}

	command := os.Args[1]

	switch command {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Failed to run migrations:", err)
		}
		log.Println("Migrations applied successfully")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Failed to rollback migrations:", err)
		}
		log.Println("Migrations rolled back successfully")
	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatal("Failed to get migration version:", err)
		}
		log.Printf("Current migration version: %d, dirty: %t", version, dirty)
	default:
		log.Fatal("Unknown command. Use: up, down, or version")
	}
}
