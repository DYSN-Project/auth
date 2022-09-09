package db

import (
	"dysn/auth/config"
	"dysn/auth/pkg/log"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	database, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		logger.ErrorLog.Panic(err)
	}
	return database
}

func StartDB(cfg *config.Config, logger *log.Logger) *gorm.DB {
	if db == nil {
		db = newDb(cfg, logger)
	}
	sql, err := db.DB()
	if err != nil {
		logger.ErrorLog.Panic(err)
	}
	if err = sql.Ping(); err != nil {
		logger.ErrorLog.Panic(err)
	}
	logger.InfoLog.Println("Connecting to database...")

	return db
}

func CloseDB(db *gorm.DB, logger *log.Logger) {
	logger.InfoLog.Println("Close database Connection")

	sql, err := db.DB()
	if err != nil {
		logger.ErrorLog.Panic(err)
	}
	err = sql.Close()
	if err != nil {
		logger.ErrorLog.Panic(err)
	}
}
