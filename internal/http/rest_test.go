package http_test

import (
	"encoding/json"
	"hh_buff/internal/db"
	htp "hh_buff/internal/http"
	"hh_buff/internal/models"
	"hh_buff/internal/repo"
	"hh_buff/pkg/hh"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

var restCurrentPath = "/rest/current"
var restDataPath = "/rest/data"

func setupTestRestController(t *testing.T) *gin.Engine {
	t.Helper()

	d, err := db.NewSQLiteInMemory()
	if err != nil {
		t.Fatal(err)
	}

	q := []*models.DBQuery{
		{
			Name: "Go",
			Query: hh.GetVacanciesRequest{
				Text: "Go",
			},
		},
		{
			Name: "Java",
			Query: hh.GetVacanciesRequest{
				Text: "Java",
			},
		},
	}

	if err := d.Create(q).Error; err != nil {
		t.Fatal(err)
	}

	s := []*models.DBSnapshot{
		{
			Count:     250,
			QueryID:   q[0].ID,
			CreatedAt: time.Now().Add(-time.Hour),
		},
		{
			Count:     150,
			QueryID:   q[1].ID,
			CreatedAt: time.Now().Add(-time.Hour),
		},
		{
			Count:     150,
			QueryID:   q[0].ID,
			CreatedAt: time.Now(),
		},
		{
			Count:     50,
			QueryID:   q[1].ID,
			CreatedAt: time.Now(),
		},
	}

	if err := d.Create(s).Error; err != nil {
		t.Fatal(err)
	}

	qr := repo.NewDBQueryRepo(d)
	sr := repo.NewDBSnapshotRepo(d)
	rc := htp.NewRestController(qr, sr)

	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.GET(restCurrentPath, rc.Current)
	router.GET(restDataPath, rc.Data)

	return router
}

func TestRestControllerCurrent(t *testing.T) {
	router := setupTestRestController(t)

	req, _ := http.NewRequest(http.MethodGet, restCurrentPath, nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var ens map[string]int
	if err := json.Unmarshal(w.Body.Bytes(), &ens); err != nil {
		t.Fatal(err)
	}

	if ens["Go"] != 150.0 {
		t.Errorf("Go got: %v, want: %v", ens["Go"], 150)
	}

	if ens["Java"] != 50 {
		t.Errorf("Java got: %v, want: %v", ens["Java"], 50)
	}
}

func TestRestControllerData(t *testing.T) {
	router := setupTestRestController(t)

	req, _ := http.NewRequest(http.MethodGet, restDataPath, nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var ens []models.DBSnapshot
	if err := json.Unmarshal(w.Body.Bytes(), &ens); err != nil {
		t.Fatal(err)
	}

	for i, en := range []models.DBSnapshot{
		{
			Count:   250,
			QueryID: 1,
			Query: &models.DBQuery{
				Name: "Go",
			},
		},
		{
			Count:   150,
			QueryID: 2,
			Query: &models.DBQuery{
				Name: "Go",
			},
		},
		{
			Count:   150,
			QueryID: 1,
			Query: &models.DBQuery{
				Name: "Java",
			},
		},
		{
			Count:   50,
			QueryID: 2,
			Query: &models.DBQuery{
				Name: "Java",
			},
		},
	} {
		if ens[i].Count != en.Count {
			t.Errorf("Encounter got: %v, want: %v", ens[i].Count, en.Count)
		}

		if ens[i].QueryID != ens[i].QueryID {
			t.Errorf("Encounter got: %v, want: %v", ens[i].QueryID, ens[i].QueryID)
		}

		if ens[i].Query.Name != ens[i].Query.Name {
			t.Errorf("Encounter got: %v, want: %v", ens[i].Query.Name, ens[i].Query.Name)
		}
	}
}
