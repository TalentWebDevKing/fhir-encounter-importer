package main

import (
	"fmt"
	"log"
	"github.com/TalentWebDevKing/fhir-encounter-importer/internal/db"
	"github.com/TalentWebDevKing/fhir-encounter-importer/internal/utils"
	"github.com/TalentWebDevKing/fhir-encounter-importer/internal/fhir"
)

func main() {
	// Load DB config from env
	cfg := db.Config{
		Host:     utils.GetEnv("DB_HOST", "localhost"),
		Port:     utils.GetEnv("DB_PORT", "5432"),
		User:     utils.GetEnv("DB_USER", "fhiruser"),
		Password: utils.GetEnv("DB_PASSWORD", "fhirpass"),
		DBName:   utils.GetEnv("DB_NAME", "fhirdb"),
	}

	// Connect to DB
	conn, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer conn.Close()
	fmt.Println("✅ Connected to PostgreSQL!")

	entries, err := fhir.FetchEncounters(1)
	if err != nil {
		log.Fatalf("❌ Failed to fetch encounters: %v", err)
	}
	fmt.Printf("✅ Got %d encounters!\n", len(entries))
}
