package repo

import (
	"hh_buff/internal/models"
	"time"

	"gorm.io/gorm"
)

type DBSnapshotRepo struct {
	db *gorm.DB
}

func NewDBSnapshotRepo(db *gorm.DB) *DBSnapshotRepo {
	return &DBSnapshotRepo{db}
}

func (r *DBSnapshotRepo) Save(snapshot *models.DBSnapshot) error {
	return r.db.Save(snapshot).Error
}

func (r *DBSnapshotRepo) GetAll() ([]*models.DBSnapshot, error) {
	var snapshots []*models.DBSnapshot
	return snapshots, r.db.Preload("Query").Find(&snapshots).Error
}

func (r *DBSnapshotRepo) GetAllByQuery(queryID uint) ([]*models.DBSnapshot, error) {
	var snapshots []*models.DBSnapshot
	return snapshots, r.db.Preload("Query").Where("query_id = ?", queryID).Find(&snapshots).Error
}

func (r *DBSnapshotRepo) GetLatestByQuery(queryID uint) (*models.DBSnapshot, error) {
	var snapshot models.DBSnapshot
	return &snapshot, r.db.Preload("Query").Where("query_id = ?", queryID).Order("created_at DESC").First(&snapshot).Error
}

func (r *DBSnapshotRepo) GetCurrentCount(queryID uint) (int, error) {
	var snapshot models.DBSnapshot
	if err := r.db.Where("query_id = ?", queryID).Order("created_at DESC").First(&snapshot).Error; err != nil {
		return 0, err
	}

	return snapshot.Count, nil
}

func (r *DBSnapshotRepo) GetByQueryIDsAndDate(queryID []uint, startDate *time.Time, endDate *time.Time) ([]*models.DBSnapshot, error) {
	var snapshots []*models.DBSnapshot

	db := r.db.Preload("Query")

	if len(queryID) > 0 {
		db = db.Where("query_id IN (?)", queryID)
	}

	if startDate != nil && endDate != nil {
		db = db.Where("created_at BETWEEN ? AND ?", *startDate, *endDate)
	} else if startDate != nil {
		db = db.Where("created_at >= ?", *startDate)
	} else if endDate != nil {
		db = db.Where("created_at <= ?", *endDate)
	}

	err := db.Find(&snapshots).Error
	return snapshots, err
}
