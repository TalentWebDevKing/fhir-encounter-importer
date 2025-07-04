package fhir

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const baseURL = "https://hapi.fhir.org/baseR4"

type FHIREncounterListResponse struct {
    Entry []FHIREncounterEntry `json:"entry"`
}

type FHIREncounterEntry struct {
    Resource struct {
        ID          string `json:"id"`
        Subject     struct {
            Reference string `json:"reference"`
        } `json:"subject"`
        Participant []struct {
            Individual struct {
                Reference string `json:"reference"`
            } `json:"individual"`
        } `json:"participant"`
        Period struct {
            Start string `json:"start"`
        } `json:"period"`
    } `json:"resource"`
}

func FetchEncounters(page int) ([]FHIREncounterEntry, error) {
    url := fmt.Sprintf("%s/Encounter?_count=10&_page=%d", baseURL, page)
    resp, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("fetch error: %w", err)
    }
    defer resp.Body.Close()

    var data FHIREncounterListResponse
    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        return nil, fmt.Errorf("decode error: %w", err)
    }

    return data.Entry, nil
}


type FHIREncounterListResponseEntry = struct {
	Resource struct {
		ID        string `json:"id"`
		Subject   struct {
			Reference string `json:"reference"`
		} `json:"subject"`
		Participant []struct {
			Individual struct {
				Reference string `json:"reference"`
			} `json:"individual"`
		} `json:"participant"`
		Period struct {
			Start string `json:"start"`
		} `json:"period"`
	}
}
