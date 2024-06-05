package database

import (
	"context"
	"database/sql"
	"github.com/pythoninja/go-redirect/internal/config"
	"time"

	// Importing "lib/pq" for side effects only.
	_ "github.com/lib/pq"
)

func NewConnectionPool(app *config.Application) (*sql.DB, error) {
	duration, err := time.ParseDuration(app.Config.Database.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", app.Config.Database.Dsn)
	db.SetMaxIdleConns(app.Config.Database.MaxIdleConns)
	db.SetMaxOpenConns(app.Config.Database.MaxOpenConns)
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
