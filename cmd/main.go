package main

import (
	"fmt"
	"log"
	"os"

	"github.com/TalentWebDevKing/fhir-encounter-importer/internal/db"
)

func main() {
	// Load DB config from env
	cfg := db.Config{
		Host:     utils.GetEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "fhiruser"),
		Password: getEnv("DB_PASSWORD", "fhirpass"),
		DBName:   getEnv("DB_NAME", "fhirdb"),
	}

	// Connect to DB
	conn, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer conn.Close()
	fmt.Println("✅ Connected to PostgreSQL!")

	// TODO: Start importer here
}
