package storage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/pythoninja/go-redirect/internal/model"
	"time"
)

type LinksStorage struct {
	db *sql.DB
}

func (s LinksStorage) GetAllLinks() ([]*model.Link, error) {
	query := `select id, alias, long_url, clicks, created_at from links`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	links := []*model.Link{}

	for rows.Next() {
		var link model.Link

		err := rows.Scan(
			&link.Id,
			&link.Alias,
			&link.Url,
			&link.Clicks,
			&link.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		links = append(links, &link)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return links, nil
}

func (s LinksStorage) GetLinkById(id int64) (*model.Link, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `select id, alias, long_url, clicks, created_at from links where id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var link model.Link

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&link.Id,
		&link.Alias,
		&link.Url,
		&link.Clicks,
		&link.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &link, nil
}

func (s LinksStorage) GetUrlByAlias(alias string) (string, error) {
	query := `select long_url from links where alias = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var url string

	err := s.db.QueryRowContext(ctx, query, alias).Scan(&url)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return "", ErrRecordNotFound
		default:
			return "", err
		}
	}

	return url, nil
}

func (s LinksStorage) UpdateClicksByAlias(alias string) error {
	query := `update links set clicks = clicks + 1 where alias = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, alias)

	return err
}
