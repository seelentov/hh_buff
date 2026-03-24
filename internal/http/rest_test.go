package http_test

import (
	"bytes"
	"encoding/json"
	"fmt"
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

var restDataPath = "/rest/data"
var restQueriesPath = "/rest/queries"
var restUploadPath = "/rest/upload"

var restController *gin.Engine

func setupTestRestController(t *testing.T) *gin.Engine {
	t.Helper()

	if restController != nil {
		return restController
	}

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
	hc := hh.NewClient()
	rc := htp.NewRestController(qr, sr, hc)

	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.GET(restDataPath, rc.Data)
	router.GET(restQueriesPath, rc.Queries)
	router.POST(restUploadPath, rc.UploadQuery)

	restController = router
	return restController
}

func TestRestControllerData(t *testing.T) {
	router := setupTestRestController(t)

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?start_date=%s&end_date=%s", restDataPath, time.Now().Add(time.Hour*24*-1).Format("2006-01-02"), time.Now().Add(time.Hour*24).Format("2006-01-02")), nil)
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

func TestRestControllerQueries(t *testing.T) {
	router := setupTestRestController(t)

	req, _ := http.NewRequest(http.MethodGet, restQueriesPath, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var queries []models.DBQuery
	if err := json.Unmarshal(w.Body.Bytes(), &queries); err != nil {
		t.Fatal(err)
	}

	expectedCount := 2
	if len(queries) != expectedCount {
		t.Errorf("Expected %d queries, got %d", expectedCount, len(queries))
	}

	if queries[0].Name != "Go" || queries[1].Name != "Java" {
		t.Errorf("Unexpected query names: %s, %s", queries[0].Name, queries[1].Name)
	}
}

func TestRestControllerUploadQuery(t *testing.T) {
	router := setupTestRestController(t)

	uploadReq := htp.UploadQueryReq{
		Name: "Python",
		Query: hh.GetVacanciesRequest{
			Text: "Python",
			Area: []string{"1", "2"},
		},
	}

	body, _ := json.Marshal(uploadReq)
	req, _ := http.NewRequest(http.MethodPost, restUploadPath, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	reqCheck, _ := http.NewRequest(http.MethodGet, restQueriesPath, nil)
	wCheck := httptest.NewRecorder()
	router.ServeHTTP(wCheck, reqCheck)

	var queries []models.DBQuery
	json.Unmarshal(wCheck.Body.Bytes(), &queries)

	found := false
	for _, q := range queries {
		if q.Name == "Python" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Uploaded query 'Python' not found in database")
	}
}

func TestRestControllerUploadQueryInvalid(t *testing.T) {
	router := setupTestRestController(t)

	body, _ := json.Marshal(map[string]string{"invalid": "data"})
	req, _ := http.NewRequest(http.MethodPost, restUploadPath, bytes.NewBuffer(body))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid payload, got %d", w.Code)
	}
}

func TestRestControllerUploadQueryExists(t *testing.T) {
	router := setupTestRestController(t)

	uploadReq := htp.UploadQueryReq{
		Name: "Go",
		Query: hh.GetVacanciesRequest{
			Text: "Go",
			Area: []string{"1", "2"},
		},
	}

	body, _ := json.Marshal(uploadReq)
	req, _ := http.NewRequest(http.MethodPost, restUploadPath, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	body, _ = json.Marshal(uploadReq)
	req, _ = http.NewRequest(http.MethodPost, restUploadPath, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status 409, got %d. Body: %s", w.Code, w.Body.String())
	}
}
