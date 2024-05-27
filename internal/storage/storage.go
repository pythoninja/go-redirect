package storage

import (
	"database/sql"
	"errors"
)

var ErrRecordNotFound = errors.New("record not found")

type Storage struct {
	Links  LinksStorage
	Health HealthStorage
}

func New(db *sql.DB) *Storage {
	return &Storage{
		Links:  LinksStorage{db: db},
		Health: HealthStorage{db: db},
	}
}
