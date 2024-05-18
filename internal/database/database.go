package database

import (
	"context"
	"database/sql"
	"github.com/pythoninja/go-redirect/internal/config"
	"time"

	_ "github.com/lib/pq"
)

func NewConnectionPool(cfg *config.Config) (*sql.DB, error) {
	duration, err := time.ParseDuration(cfg.Database.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", cfg.Database.Dsn)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetConnMaxIdleTime(duration)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, err
}
