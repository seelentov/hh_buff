package db_test

import (
	"hh_buff/internal/db"
	"hh_buff/internal/models"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSeedDefaultQueries(t *testing.T) {
	d, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	err = d.AutoMigrate(&models.DBQuery{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	t.Run("Initial seed should create 14 records", func(t *testing.T) {
		err := db.SeedDefaultQueries(d)
		if err != nil {
		}

		var count int64
		d.Model(&models.DBQuery{}).Count(&count)

		if int64(14) != count {
			t.Fatalf("expected 14, got %d", count)
		}

		var q models.DBQuery
		err = d.Where("name = ?", "Go between1And3 remote").First(&q).Error
		if err != nil {
			t.Fatalf("failed to find Go between1And3 query: %v", err)
		}
		if q.Query.Text != "Go" {
			ql := q.Query.Text
			t.Fatalf("expected Go, got %s", ql)
		}
	})

	t.Run("Running seed again should not duplicate records", func(t *testing.T) {
		err := db.SeedDefaultQueries(d)
		if err != nil {
			t.Fatalf("failed to seed: %v", err)
		}

		var count int64
		d.Model(&models.DBQuery{}).Count(&count)

		if int64(14) != count {
			t.Errorf("expected 14, got %d", count)
		}
	})
}
