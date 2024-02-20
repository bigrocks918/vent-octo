package store

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic("An error occurred when initializing mock DB: " + err.Error())
	}
	return db, mock
}
