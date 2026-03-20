package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSQLiteInMemory() (*gorm.DB, error) {
	return newSQLite("file::memory:?cache=shared")
}

func NewSQLiteFile(path string) (*gorm.DB, error) {
	return newSQLite(path)
}

func newSQLite(dsn string) (*gorm.DB, error) {
	return newDB(sqlite.Open(dsn))
}

func newDB(dialector gorm.Dialector) (*gorm.DB, error) {
	db, err := gorm.Open(dialector, config)
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(entities...); err != nil {
		return nil, err
	}

	return db, nil
}
