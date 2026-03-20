package db_test

import (
	"hh_buff/internal/db"
	"testing"
)

func TestNewSQLiteInMemory(t *testing.T) {
	db, err := db.NewSQLiteInMemory()
	if err != nil {
		t.Fatalf("failed to create in-memory SQLite database: %v", err)
	}

	if db == nil {
		t.Fatalf("db is nil")
	}
}
