package storage

import (
	"context"
	"time"
)

func (s Storage) GetDatabaseStatus() (int, error) {
	query := "select 1"
	var result int

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query).Scan(&result)
	if err != nil {
		return 0, err
	}

	return result, nil
}
