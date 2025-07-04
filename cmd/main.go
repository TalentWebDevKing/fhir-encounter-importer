package main

import (
	"fmt"
	"log"

	"github.com/TalentWebDevKing/fhir-encounter-importer/internal/db"
	"github.com/TalentWebDevKing/fhir-encounter-importer/internal/fhir"
	"github.com/TalentWebDevKing/fhir-encounter-importer/internal/utils"
	"github.com/TalentWebDevKing/fhir-encounter-importer/internal/worker"
)

func main() {
	// Load DB configuration from environment
	cfg := db.Config{
		Host:     utils.GetEnv("DB_HOST", "localhost"),
		Port:     utils.GetEnv("DB_PORT", "5432"),
		User:     utils.GetEnv("DB_USER", "fhiruser"),
		Password: utils.GetEnv("DB_PASSWORD", "fhirpass"),
		DBName:   utils.GetEnv("DB_NAME", "fhirdb"),
	}

	// Connect to PostgreSQL
	conn, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer conn.Close()
	fmt.Println("✅ Connected to PostgreSQL!")

	// Fetch encounters (page 1)
	entries, err := fhir.FetchEncounters(1)
	if err != nil {
		log.Fatalf("❌ Failed to fetch encounters: %v", err)
	}

	if len(entries) == 0 {
		log.Println("⚠️  No encounters found on page 1. Exiting.")
		return
	}

	// Create task list
	tasks := []worker.Task{}
	for _, entry := range entries {
		tasks = append(tasks, worker.Task{Encounter: entry})
	}

	// Start worker pool with 5 concurrent workers
	worker.StartPool(conn, tasks, 5)
}
