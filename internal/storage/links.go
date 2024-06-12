package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/pythoninja/go-redirect/internal/model"
	"time"
)

type linksStorage struct {
	db *sql.DB
}

func (s linksStorage) GetAll() ([]*model.Link, error) {
	query := `select id, alias, target_url, clicks, created_at from links`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch links from database: %w", err)
	}

	defer rows.Close()

	links := []*model.Link{}

	for rows.Next() {
		var link model.Link

		err := rows.Scan(
			&link.ID,
			&link.Alias,
			&link.URL,
			&link.Clicks,
			&link.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan database rows for links: %w", err)
		}

		links = append(links, &link)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during row scanning for links: %w", err)
	}

	return links, nil
}

func (s linksStorage) GetByID(id int64) (*model.Link, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `select id, alias, target_url, clicks, created_at from links where id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var link model.Link

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&link.ID,
		&link.Alias,
		&link.URL,
		&link.Clicks,
		&link.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, fmt.Errorf("failed to fetch link by ID: %w", err)
		}
	}

	return &link, nil
}

func (s linksStorage) GetURLByAlias(alias string) (string, error) {
	query := `select target_url from links where alias = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var url string

	err := s.db.QueryRowContext(ctx, query, alias).Scan(&url)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return "", ErrRecordNotFound
		default:
			return "", fmt.Errorf("failed to fetch URL by alias: %w", err)
		}
	}

	return url, nil
}

func (s linksStorage) UpdateClicksByAlias(alias string) error {
	query := `update links set clicks = clicks + 1 where alias = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, alias)
	if err != nil {
		return fmt.Errorf("failed to increment clicks by alias.: %w", err)
	}

	return nil
}

func (s linksStorage) Insert(link *model.Link) error {
	query := `insert into links (target_url, alias)
		values ($1, $2)
		returning id, created_at, clicks`

	args := []any{link.URL, link.Alias}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, args...).Scan(&link.ID, &link.CreatedAt, &link.Clicks)
	if err != nil {
		switch {
		case err.Error() == errUniqueConstraintViolationAlias.Error():
			return ErrDuplicateAlias
		default:
			return fmt.Errorf("failed to create new link: %w", err)
		}
	}

	return nil
}

func (s linksStorage) Update(link *model.Link) error {
	query := `update links
		set target_url = $1, alias = $2
		where id = $3`

	args := []any{
		link.URL,
		link.Alias,
		link.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		switch {
		case err.Error() == errUniqueConstraintViolationAlias.Error():
			return ErrDuplicateAlias
		default:
			return fmt.Errorf("failed to update link: %w", err)
		}
	}

	return nil
}

func (s linksStorage) Delete(id int64) error {
	query := `delete from links where id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete link: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to count delete link operation affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
