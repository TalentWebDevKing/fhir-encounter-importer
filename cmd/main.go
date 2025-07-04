package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/TalentWebDevKing/fhir-encounter-importer/internal/db"
	"github.com/TalentWebDevKing/fhir-encounter-importer/internal/fhir"
	"github.com/TalentWebDevKing/fhir-encounter-importer/internal/model"
	"github.com/TalentWebDevKing/fhir-encounter-importer/internal/utils"
)

func main() {
	cfg := db.Config{
		Host:     utils.GetEnv("DB_HOST", "localhost"),
		Port:     utils.GetEnv("DB_PORT", "5432"),
		User:     utils.GetEnv("DB_USER", "fhiruser"),
		Password: utils.GetEnv("DB_PASSWORD", "fhirpass"),
		DBName:   utils.GetEnv("DB_NAME", "fhirdb"),
	}

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

	if len(entries) == 0 {
		log.Println("⚠️  No encounters found on page 1. Exiting.")
		return
	}

	log.Printf("📦 Fetched %d encounters", len(entries))

	for _, entry := range entries {
		if len(entry.Resource.Participant) == 0 {
			log.Printf("⚠️  Encounter %s has no participants. Skipping.", entry.Resource.ID)
			continue
		}

		// Parse IDs
		encounterID := entry.Resource.ID
		ref := entry.Resource.Subject.Reference
		patientID := parseFHIRRef(ref)
		providerID := parseFHIRRef(entry.Resource.Participant[0].Individual.Reference)

		serviceDate, _ := time.Parse(time.RFC3339, entry.Resource.Period.Start)

		patientName, err := fhir.FetchPatient(patientID)
		if err != nil {
			log.Printf("⚠️  Could not fetch patient name: %v", err)
			patientName = "Unknown Patient"
		}

		providerName, err := fhir.FetchPractitioner(providerID)
		if err != nil {
			log.Printf("⚠️  Could not fetch practitioner name: %v", err)
			providerName = "Unknown Provider"
		}
		patient := model.Patient{ID: patientID, Name: patientName}
		provider := model.Provider{ID: providerID, Name: providerName}
		encounter := model.Encounter{
			ID:          encounterID,
			PatientID:   patientID,
			ProviderID:  providerID,
			ServiceDate: serviceDate,
		}

		if err := db.CreatePatient(conn, patient); err != nil {
			log.Fatal("insert patient:", err)
		}
		if err := db.CreateProvider(conn, provider); err != nil {
			log.Fatal("insert provider:", err)
		}
		if err := db.UpsertEncounter(conn, encounter); err != nil {
			log.Fatal("insert encounter:", err)
		}

		fmt.Printf("✅ Inserted encounter: %s\n", encounterID)
	}
}

func parseFHIRRef(ref string) string {
	parts := strings.Split(ref, "/")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}
