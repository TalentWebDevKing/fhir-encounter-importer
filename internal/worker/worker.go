package worker

import "github.com/TalentWebDevKing/fhir-encounter-importer/internal/fhir"

type Task struct {
	Encounter fhir.FHIREncounterEntry
}
