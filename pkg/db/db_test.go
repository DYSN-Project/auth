package db

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestStartDB(t *testing.T) {
	db := getDbMock()
	defer db.Close()
	_, err := gorm.Open("postgres", db)

	assert.NoError(t, err)
}

func getDbMock() *sql.DB {
	dbMock, _, err := sqlmock.New()
	defer dbMock.Close()

	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return dbMock
}
