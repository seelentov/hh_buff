package models

import (
	"hh_buff/pkg/hh"
	"time"
)

type DBQuery struct {
	ID        uint      `gorm:"primarykey" json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at"`

	Name  string                 `json:"name,omitempty"`
	Query hh.GetVacanciesRequest `gorm:"serializer:json" json:"query"`

	Snapshots []*DBSnapshot `gorm:"foreignKey:QueryID" json:"snapshots,omitempty"`
}
