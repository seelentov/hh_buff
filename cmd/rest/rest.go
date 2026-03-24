package main

import (
	"hh_buff/internal/db"
	htt "hh_buff/internal/http"
	"hh_buff/internal/repo"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	dbDriver := os.Getenv("DB_DRIVER")
	dbType := os.Getenv("DB_TYPE")
	dbDsn := os.Getenv("DB_DSN")

	var database *gorm.DB

	switch dbDriver {
	case "sqlite":
		switch dbType {
		case "in-memory":
			var err error
			database, err = db.NewSQLiteInMemory()
			if err != nil {
				log.Fatal(err)
			}
		case "file":
			var err error
			database, err = db.NewSQLiteFile(dbDsn)
			if err != nil {
				log.Fatal(err)
			}
		default:
			log.Fatal("Unsupported DB type with SQLite driver")
		}
	default:
		log.Fatal("Unsupported DB driver")
	}

	if err := db.SeedDefaultQueries(database); err != nil {
		log.Fatal(err)
	}

	queryRepo := repo.NewDBQueryRepo(database)
	snapshotRepo := repo.NewDBSnapshotRepo(database)

	ctrl := htt.NewRestController(queryRepo, snapshotRepo)

	r := gin.Default()

	r.GET("/rest/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/rest/data", ctrl.Data)
	r.GET("/rest/current", ctrl.Current)
	r.GET("/rest/queries", ctrl.Queries)

	restPort := os.Getenv("REST_PORT")
	if err := r.Run(":" + restPort); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
