package models

import (
	"time"
)

type DBSnapshot struct {
	ID        uint      `gorm:"primarykey" json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at"`

	Count int `json:"count,omitempty"`

	QueryID uint     `json:"query_id,omitempty"`
	Query   *DBQuery `gorm:"foreignKey:QueryID" json:"query,omitempty"`
}
