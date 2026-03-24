package repo

import (
	"encoding/json"
	"errors"
	"hh_buff/internal/models"
	"strings"

	"gorm.io/gorm"
)

var ErrAlreadyExists = errors.New("already exists")
var ErrInvalidQuery = errors.New("invalid query")

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

func (r *DBQueryRepo) Save(query *models.DBQuery) (bool, error) {
	if query == nil || strings.TrimSpace(query.Query.Text) == "" {
		return false, ErrInvalidQuery
	}

	queryBytes, err := json.Marshal(query.Query)
	if err != nil {
		return false, err
	}

	err = r.db.Where("query = ?", string(queryBytes)).First(&models.DBQuery{}).Error

	if err == nil {
		return false, ErrAlreadyExists
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}

	if err := r.db.Create(query).Error; err != nil {
		return false, err
	}

	return true, nil
}

func (r *DBQueryRepo) Delete(id uint) error {
	return r.db.Delete(&models.DBQuery{}, id).Error
}
