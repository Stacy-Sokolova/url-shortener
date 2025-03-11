package pgdb

import (
	"context"
	"fmt"

	fnc "url-server/internal/storage/data_processing"

	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

func NewSQLStorage(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) CreateShortURL(ctx context.Context, fullURL string) (string, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return "", fmt.Errorf("postgres.CreateShortURL - db.Begin: %v", err)
	}

	//check if original url is already exist then just return short url from db

	shortURL := fnc.URLShortener(fullURL)
	createItemQuery := fmt.Sprintf("INSERT INTO %s (full_url, short_url) VALUES ($1, $2)", "urltable")

	_, err = tx.Exec(createItemQuery, fullURL, shortURL)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("postgres.CreateShortURL - tx.Exec: %v", err)
	}

	tx.Commit()

	return shortURL, nil
}

func (s *Storage) GetFullURL(ctx context.Context, shortURL string) (string, error) {
	var fullURL string
	query := fmt.Sprintf(`SELECT ti.full_url FROM %s ti WHERE ti.short_url = $1`,
		"urltable")
	if err := s.db.Get(&fullURL, query, shortURL); err != nil {
		return "", fmt.Errorf("postgres.GEtFullURL - db.Get: %v", err)
	}

	return fullURL, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
