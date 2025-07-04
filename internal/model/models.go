package model

import "time"

type Patient struct {
	ID   string
	Name string
}

type Provider struct {
	ID   string
	Name string
}

type Encounter struct {
	ID         string
	PatientID  string
	ProviderID string
	ServiceDate time.Time
}
