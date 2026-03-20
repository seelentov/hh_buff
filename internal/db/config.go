package db

import (
	"hh_buff/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var entities = []interface{}{
	&models.DBQuery{},
	&models.DBSnapshot{},
}

var config = &gorm.Config{
	Logger: logger.Default.LogMode(logger.Info),
}
