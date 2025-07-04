package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/yourusername/fhir-encounter-importer/internal/model"
)

// CreatePatient inserts a patient if not exists
func CreatePatient(db *sql.DB, patient model.Patient) error {
	_, err := db.Exec(`
		INSERT INTO patients (id, name)
		VALUES ($1, $2)
		ON CONFLICT (id) DO NOTHING
	`, patient.ID, patient.Name)
	return err
}

// CreateProvider inserts a provider if not exists
func CreateProvider(db *sql.DB, provider model.Provider) error {
	_, err := db.Exec(`
		INSERT INTO providers (id, name)
		VALUES ($1, $2)
		ON CONFLICT (id) DO NOTHING
	`, provider.ID, provider.Name)
	return err
}

// UpsertEncounter inserts or updates an encounter
func UpsertEncounter(db *sql.DB, encounter model.Encounter) error {
	_, err := db.Exec(`
		INSERT INTO encounters (id, patient_id, provider_id, service_date)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			patient_id = EXCLUDED.patient_id,
			provider_id = EXCLUDED.provider_id,
			service_date = EXCLUDED.service_date
	`, encounter.ID, encounter.PatientID, encounter.ProviderID, encounter.ServiceDate)
	return err
}