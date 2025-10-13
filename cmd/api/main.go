package main

import (
	"github.com/andrew-hayworth22/critiquefi-service/internal/config"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}
}
