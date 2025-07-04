package fhir

import (
	"encoding/json"
	"fmt"
	"net/http"
    "strings"
	"io"
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

func FetchPatient(id string) (string, error) {
	url := fmt.Sprintf("%s/Patient/%s", baseURL, id)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("fetch patient: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("bad response: %s", string(body))
	}

	var data struct {
		Name []struct {
			Given  []string `json:"given"`
			Family string   `json:"family"`
		} `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("decode patient: %w", err)
	}

	if len(data.Name) == 0 {
		return "Unknown Patient", nil
	}
	fullName := strings.Join(data.Name[0].Given, " ") + " " + data.Name[0].Family
	return fullName, nil
}

func FetchPractitioner(id string) (string, error) {
	url := fmt.Sprintf("%s/Practitioner/%s", baseURL, id)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("fetch practitioner: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("bad response: %s", string(body))
	}

	var data struct {
		Name []struct {
			Given  []string `json:"given"`
			Family string   `json:"family"`
		} `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("decode practitioner: %w", err)
	}

	if len(data.Name) == 0 {
		return "Unknown Doctor", nil
	}
	fullName := strings.Join(data.Name[0].Given, " ") + " " + data.Name[0].Family
	return fullName, nil
}

