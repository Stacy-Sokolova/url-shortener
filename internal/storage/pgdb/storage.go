package pgdb

import (
	"context"
	"fmt"

	pb "url-server/internal/service/proto"
	fnc "url-server/internal/storage/data_processing"

	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) CreateShortURL(ctx context.Context, r *pb.Request) (*pb.Response, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	fullURL := r.GetUrl()
	//shortURL := "newshorturl"
	shortURL := fnc.URLShortener(fullURL)
	createItemQuery := fmt.Sprintf("INSERT INTO %s (full_url, short_url) VALUES ($1, $2)", "urltable")

	_, err = tx.Exec(createItemQuery, fullURL, shortURL)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &pb.Response{Url: shortURL}, nil
}

func (s *Storage) GetFullURL(ctx context.Context, r *pb.Request) (*pb.Response, error) {
	shortURL := r.GetUrl()

	var fullURL string
	query := fmt.Sprintf(`SELECT ti.full_url FROM %s ti WHERE ti.short_url = $1`,
		"urltable")
	if err := s.db.Get(&fullURL, query, shortURL); err != nil {
		return nil, err
	}

	return &pb.Response{Url: fullURL}, nil
}
