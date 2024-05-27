package storage

import (
	"database/sql"
	"errors"
)

var ErrRecordNotFound = errors.New("record not found")

type Storage struct {
	Links  linksStorage
	Health healthStorage
}

func New(db *sql.DB) *Storage {
	return &Storage{
		Links:  linksStorage{db: db},
		Health: healthStorage{db: db},
	}
}
