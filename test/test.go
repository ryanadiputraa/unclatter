package test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewMockDB(t *testing.T) (*gorm.DB, *sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("fail to create mock db conn: ", err.Error())
	}

	gormDB, err := gorm.Open(postgres.New(
		postgres.Config{
			Conn:       db,
			DriverName: "postgres",
		},
	), &gorm.Config{})
	if err != nil {
		t.Fatal("fail to open db conn: ", err.Error())
	}
	return gormDB, db, mock
}
