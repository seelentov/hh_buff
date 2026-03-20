package main

import (
	"hh_buff/internal/daemon"
	"hh_buff/internal/db"
	"hh_buff/internal/repo"
	"hh_buff/pkg/hh"
	"log"
	"os"
	"strconv"
	"time"

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

	hhClient := hh.NewClient()

	queryRepo := repo.NewDBQueryRepo(database)
	snapshotRepo := repo.NewDBSnapshotRepo(database)

	scDaemonIntervalSec := os.Getenv("SNAPSHOT_CREATOR_INTERVAL_SEC")
	scDaemonIntervalSecI, err := strconv.ParseInt(scDaemonIntervalSec, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	scDaemonInterval := time.Duration(scDaemonIntervalSecI) * time.Second

	scDaemonPeriodSec := os.Getenv("SNAPSHOT_CREATOR_PERIOD_SEC")
	scDaemonPeriodSecI, err := strconv.ParseInt(scDaemonPeriodSec, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	scDaemonPeriod := time.Duration(scDaemonPeriodSecI) * time.Second

	scDaemon := daemon.NewSnapshotCreator(queryRepo, snapshotRepo, hhClient, scDaemonInterval, scDaemonPeriod)
	scDaemon.Start()

	select {}
}
