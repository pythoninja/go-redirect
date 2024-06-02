package storage

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrDuplicateAlias = errors.New("alias is already exists")
)

type Storage struct {
	Links  linksStorage
	Health HealthStorage
}

func New(db *sql.DB) *Storage {
	return &Storage{
		Links:  linksStorage{db: db},
		Health: HealthStorage{db: db},
	}
}
