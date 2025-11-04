package main

import (
	"flag"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	// Get action.go
	actionString := flag.String("action", "up", "up, down")

	flag.Parse()

	action, err := ParseAction(*actionString)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Load environment
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	dbURL := os.Getenv("DB_URL")

	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		log.Fatalf("error connecting to database for migrations: %v", err)
	}

	switch action {
	case Up:
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}
		return
	case Down:
		if err := m.Down(); err != nil {
			log.Fatal(err)
		}
		return
	}

}
