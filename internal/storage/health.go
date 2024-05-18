package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

func (s Storage) GetDatabaseStatus() (string, error) {
	query := "select 'hello'"
	var status string

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query).Scan(&status)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return "", errors.New("record not found")
		default:
			return "", err
		}
	}

	return status, nil
}
