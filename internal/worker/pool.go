package worker

import (
	"database/sql"
	"log"
	"sync"
	"time"
	"strings"

	"github.com/TalentWebDevKing/fhir-encounter-importer/internal/db"
	"github.com/TalentWebDevKing/fhir-encounter-importer/internal/model"
	"github.com/TalentWebDevKing/fhir-encounter-importer/internal/fhir"
	"github.com/TalentWebDevKing/fhir-encounter-importer/internal/metrics"
)

// StartPool launches worker goroutines
func StartPool(conn *sql.DB, tasks []Task, numWorkers int) {
	jobs := make(chan Task)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for task := range jobs {
				processEncounter(conn, task)
			}
		}(i)
	}

	// Feed tasks into the channel
	for _, task := range tasks {
		jobs <- task
	}
	close(jobs)

	wg.Wait()
}

func processEncounter(conn *sql.DB, task Task) {
	e := task.Encounter.Resource

	if len(e.Participant) == 0 {
		log.Printf("⚠️  Encounter %s has no participants. Skipping.", e.ID)
		return
	}

	// Extract IDs
	encounterID := e.ID
	patientID := parseFHIRRef(e.Subject.Reference)
	providerID := parseFHIRRef(e.Participant[0].Individual.Reference)
	serviceDate, _ := time.Parse(time.RFC3339, e.Period.Start)

	// Fetch names
	patientName, err := fhir.FetchPatient(patientID)
	if err != nil {
		log.Printf("⚠️  Patient fetch error (%s): %v", patientID, err)
		patientName = "Unknown Patient"
	}

	providerName, err := fhir.FetchPractitioner(providerID)
	if err != nil {
		log.Printf("⚠️  Practitioner fetch error (%s): %v", providerID, err)
		providerName = "Unknown Provider"
	}

	// Insert into DB
	err = db.CreatePatient(conn, model.Patient{ID: patientID, Name: patientName})
	if err != nil {
		log.Printf("❌ Patient insert error: %v", err)
	} else {
		metrics.Incr("patients.created")
	}

	err = db.CreateProvider(conn, model.Provider{ID: providerID, Name: providerName})
	if err != nil {
		log.Printf("❌ Provider insert error: %v", err)
	} else {
		metrics.Incr("providers.created")
	}

	err = db.UpsertEncounter(conn, model.Encounter{
		ID:          encounterID,
		PatientID:   patientID,
		ProviderID:  providerID,
		ServiceDate: serviceDate,
	})
	if err != nil {
		log.Printf("❌ Encounter insert error: %v", err)
		metrics.Incr("encounters.import.error")
	} else {
		log.Printf("✅ Imported encounter %s", encounterID)
		metrics.Incr("encounters.import.success")
	}
}

func parseFHIRRef(ref string) string {
	parts := strings.Split(ref, "/")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}
