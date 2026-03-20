package repo

import (
	"hh_buff/internal/models"

	"gorm.io/gorm"
)

type DBQueryRepo struct {
	db *gorm.DB
}

func NewDBQueryRepo(db *gorm.DB) *DBQueryRepo {
	return &DBQueryRepo{db}
}

func (r *DBQueryRepo) Get(id uint) (*models.DBQuery, error) {
	var query models.DBQuery
	return &query, r.db.First(&query, id).Error
}

func (r *DBQueryRepo) GetByText(text string) ([]*models.DBQuery, error) {
	var query []*models.DBQuery
	return query, r.db.Where("LOWER(name) LIKE LOWER(?)", "%"+text+"%").Find(&query).Error
}

func (r *DBQueryRepo) GetAll() ([]*models.DBQuery, error) {
	var query []*models.DBQuery
	return query, r.db.Find(&query).Error
}

func (r *DBQueryRepo) Save(query *models.DBQuery) error {
	return r.db.Save(query).Error
}

func (r *DBQueryRepo) Delete(id uint) error {
	return r.db.Delete(&models.DBQuery{}, id).Error
}
