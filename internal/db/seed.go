package db

import (
	"hh_buff/internal/models"
	"hh_buff/pkg/hh"

	"gorm.io/gorm"
)

func SeedDefaultQueries(db *gorm.DB) error {
	pls := []string{"Go", "Java", "C++", "Python Backend", "Rust", "C#", "PHP"}
	queries := make([]models.DBQuery, 0, len(pls)*2)
	for _, pl := range pls {
		exps := []string{"between1And3", "between3And6"}
		for _, exp := range exps {
			queries = append(queries, models.DBQuery{
				Name: pl + " " + exp + " " + "remote",
				Query: hh.GetVacanciesRequest{
					Text:       pl,
					Experience: exp,
					Schedule:   []string{"remote"},
				},
			})
		}
	}

	for _, q := range queries {
		var qs []*models.DBQuery
		if err := db.Where("name = ?", q.Name).Limit(1).Find(&qs).Error; err != nil {
			return err
		}

		if len(qs) > 0 {
			continue
		}

		if err := db.Create(&q).Error; err != nil {
			return err
		}
	}

	return nil
}
