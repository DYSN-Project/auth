package db

import (
	"fmt"
	"github.com/DYSN-Project/auth/config"
	"github.com/DYSN-Project/auth/internal/packages/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

const dbUri = "host=%s user=%s dbname=%s port=%s sslmode=disable password=%s"

func newDb(cfg *config.Config, logger *log.Logger) *gorm.DB {
	url := fmt.Sprintf(dbUri,
		cfg.GetDbHost(),
		cfg.GetDbUsername(),
		cfg.GetDbName(),
		cfg.GetDbPort(),
		cfg.GetDbPassword(),
	)

	db, err := gorm.Open("postgres", url)
	if err != nil {
		logger.ErrorLog.Panic(err)
	}

	return db
}

func StartDB(cfg *config.Config, logger *log.Logger) *gorm.DB {
	if db == nil {
		db = newDb(cfg, logger)
	}
	logger.InfoLog.Println("Connecting to database...")

	return db
}

func CloseDB(db *gorm.DB, logger *log.Logger) {
	logger.InfoLog.Println("Close database Connection")

	if err := db.Close(); err != nil {
		logger.ErrorLog.Panic(err)
	}
}