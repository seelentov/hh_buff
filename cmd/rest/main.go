package main

import (
	_ "hh_buff/docs"
	"hh_buff/internal/db"
	htt "hh_buff/internal/http"
	"hh_buff/internal/repo"
	"hh_buff/pkg/hh"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// @title           HH Buff API
// @version         1.0
// @description     Сервис для мониторинга вакансий hh.ru.
// @BasePath        /rest

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

	hhClient := hh.NewClient()

	queryRepo := repo.NewDBQueryRepo(database)
	snapshotRepo := repo.NewDBSnapshotRepo(database)

	ctrl := htt.NewRestController(queryRepo, snapshotRepo, hhClient)

	r := gin.Default()

	r.GET("/rest/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/rest/data", ctrl.Data)
	r.GET("/rest/queries", ctrl.Queries)
	r.POST("/rest/queries", ctrl.UploadQuery)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	restPort := os.Getenv("REST_PORT")
	if err := r.Run(":" + restPort); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
