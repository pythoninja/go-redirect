package storage

import (
	"context"
	"github.com/pythoninja/go-redirect/internal/model"
	"time"
)

func (s Storage) GetAllLinks() ([]*model.Link, error) {
	query := "select id, created_at, short_url, long_url, clicks from links"

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
