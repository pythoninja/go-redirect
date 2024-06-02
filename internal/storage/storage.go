package storage

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound                 = errors.New("record not found")
	ErrDuplicateAlias                 = errors.New("alias is already exists")
	errUniqueConstraintViolationAlias = errors.New(
		`pq: duplicate key value violates unique constraint "links_alias_key"`)
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
