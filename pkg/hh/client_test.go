package hh

import "testing"

func TestClient_GetVacancies(t *testing.T) {
	client := NewClient()

	_, err := client.GetVacancies(GetVacanciesRequest{})
	if err != nil {
		t.Errorf("GetVacancies failed: %v", err)
	}
}
