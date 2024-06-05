package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type HealthStorage struct {
	db *sql.DB
}

func (s HealthStorage) GetDatabaseStatus() (int, error) {
	query := `select 1`

	var result int

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query).Scan(&result)
	if err != nil {
		return 0, fmt.Errorf("failed to get database status: %w", err)
	}

	return result, nil
}
