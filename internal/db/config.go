package db

import (
	"hh_buff/internal/models"

	"gorm.io/gorm"
)

var entities = []interface{}{
	&models.DBQuery{},
	&models.DBSnapshot{},
}

var config = &gorm.Config{
	// Logger: logger.Default.LogMode(logger.Info),
}
