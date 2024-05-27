package storage

import (
	"database/sql"
	"errors"
)

var ErrRecordNotFound = errors.New("record not found")

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}
