package repo_test

import (
	db "hh_buff/internal/db"
	"testing"

	"hh_buff/internal/models"
	"hh_buff/internal/repo"

	"gorm.io/gorm"
)

func setupSnapshotTestDB(t *testing.T) (*gorm.DB, *repo.DBSnapshotRepo, *repo.DBQueryRepo) {
	database, err := db.NewSQLiteInMemory()
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}

	err = database.AutoMigrate(&models.DBQuery{}, &models.DBSnapshot{})
	if err != nil {
		t.Fatalf("migration failed: %v", err)
	}

	return database, repo.NewDBSnapshotRepo(database), repo.NewDBQueryRepo(database)
}

func TestDBSnapshotRepo_Save(t *testing.T) {
	_, r, qRepo := setupSnapshotTestDB(t)

	query := &models.DBQuery{Name: "Go Test"}
	qRepo.Save(query)

	snapshot := &models.DBSnapshot{
		QueryID: query.ID,
		Count:   100,
	}

	if err := r.Save(snapshot); err != nil {
		t.Errorf("Save() error: %v", err)
	}

	if snapshot.ID == 0 {
		t.Error("expected non-zero ID for snapshot")
	}
}

func TestDBSnapshotRepo_GetAllByQuery(t *testing.T) {
	_, r, qRepo := setupSnapshotTestDB(t)

	q1 := &models.DBQuery{Name: "Query 1"}
	q2 := &models.DBQuery{Name: "Query 2"}
	qRepo.Save(q1)
	qRepo.Save(q2)

	r.Save(&models.DBSnapshot{QueryID: q1.ID, Count: 10})
	r.Save(&models.DBSnapshot{QueryID: q1.ID, Count: 20})
	r.Save(&models.DBSnapshot{QueryID: q2.ID, Count: 30})

	t.Run("Filter by Q1", func(t *testing.T) {
		results, err := r.GetAllByQuery(q1.ID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if len(results) != 2 {
			t.Errorf("expected 2 snapshots for Q1, got %d", len(results))
		}

		for _, s := range results {
			if s.QueryID != q1.ID {
				t.Errorf("expected QueryID %d, got %d", q1.ID, s.QueryID)
			}
		}
	})

	t.Run("Empty results", func(t *testing.T) {
		results, err := r.GetAllByQuery(999)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if len(results) != 0 {
			t.Errorf("expected 0 snapshots, got %d", len(results))
		}
	})
}
