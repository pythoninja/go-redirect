package storage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/pythoninja/go-redirect/internal/model"
	"time"
)

type linksStorage struct {
	db *sql.DB
}

func (s Storage) GetAllLinks() ([]*model.Link, error) {
	query := `select id, created_at, short_url, long_url, clicks from links`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := s.Links.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	links := []*model.Link{}

	for rows.Next() {
		var link model.Link

		err := rows.Scan(
			&link.Id,
			&link.CreatedAt,
			&link.ShortLink,
			&link.LongLink,
			&link.Clicks,
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

func (s Storage) GetLinkById(id int64) (*model.Link, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `select id, created_at, short_url, long_url, clicks
			from links
			where id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var link model.Link

	err := s.Links.db.QueryRowContext(ctx, query, id).Scan(
		&link.Id,
		&link.CreatedAt,
		&link.ShortLink,
		&link.LongLink,
		&link.Clicks,
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
