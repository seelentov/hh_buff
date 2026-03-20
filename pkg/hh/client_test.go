package hh

import "testing"

func TestClient_GetVacancies(t *testing.T) {
	client := NewClient()

	data, err := client.GetVacancies(GetVacanciesRequest{})
	if err != nil {
		t.Errorf("GetVacancies failed: %v", err)
	}

	if data == nil {
		t.Errorf("GetVacancies returned nil")
	}
}
