package repo_test

import (
	"errors"
	"hh_buff/internal/db"
	"testing"

	"hh_buff/internal/models"
	"hh_buff/internal/repo"
	"hh_buff/pkg/hh"

	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, *repo.DBQueryRepo) {
	database, err := db.NewSQLiteInMemory()
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}

	err = database.AutoMigrate(&models.DBQuery{})
	if err != nil {
		t.Fatalf("migration failed: %v", err)
	}

	return database, repo.NewDBQueryRepo(database)
}

func TestDBQueryRepo_SaveAndGet(t *testing.T) {
	_, r := setupTestDB(t)

	newQuery := &models.DBQuery{
		Name: "Golang Senior",
		Query: hh.GetVacanciesRequest{
			Text: "Go",
		},
	}

	if _, err := r.Save(newQuery); err != nil {
		t.Errorf("Save() error: %v", err)
	}

	if newQuery.ID == 0 {
		t.Error("expected non-zero ID")
	}

	found, err := r.Get(newQuery.ID)
	if err != nil {
		t.Errorf("Get() error: %v", err)
	}

	if found.Name != "Golang Senior" {
		t.Errorf("expected %s, got %s", "Golang Senior", found.Name)
	}
}

func TestDBQueryRepo_GetByText(t *testing.T) {
	_, r := setupTestDB(t)

	r.Save(&models.DBQuery{Name: "Python Developer", Query: hh.GetVacanciesRequest{Text: "Python Developer"}})

	t.Run("Exists", func(t *testing.T) {
		res, err := r.GetByText("Python Developer")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if len(res) == 0 || res[0].Name != "Python Developer" {
			t.Error("record not found")
		}
	})
}

func TestDBQueryRepo_GetAll(t *testing.T) {
	_, r := setupTestDB(t)

	_, err := r.Save(&models.DBQuery{Name: "Q1", Query: hh.GetVacanciesRequest{Text: "Q1"}})
	if err != nil {
		t.Errorf("Save() error: %v", err)
	}
	_, err = r.Save(&models.DBQuery{Name: "Q2", Query: hh.GetVacanciesRequest{Text: "Q2"}})
	if err != nil {
		t.Errorf("Save() error: %v", err)
	}

	all, err := r.GetAll()
	if err != nil {
		t.Errorf("GetAll() error: %v", err)
	}

	if len(all) < 2 {
		t.Errorf("expected 2 records, got %d", len(all))
	}
}

func TestDBQueryRepo_Delete(t *testing.T) {
	_, r := setupTestDB(t)

	q := &models.DBQuery{Name: "DeleteMe", Query: hh.GetVacanciesRequest{Text: "DeleteMe"}}
	r.Save(q)

	if err := r.Delete(q.ID); err != nil {
		t.Errorf("Delete() error: %v", err)
	}

	_, err := r.Get(q.ID)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("expected ErrRecordNotFound, got: %v", err)
	}
}
