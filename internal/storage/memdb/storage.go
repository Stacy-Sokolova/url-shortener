package memdb

import (
	"context"
	"fmt"

	fnc "url-server/internal/storage/data_processing"

	"github.com/dgraph-io/badger"
)

type Storage struct {
	db *badger.DB
}

func NewInmemStorage(db *badger.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) CreateShortURL(ctx context.Context, fullURL string) (string, error) {
	shortURL := fnc.URLShortener(fullURL)
	val := []byte(fullURL)
	key := []byte(shortURL)
	err := s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, val)
	})
	if err != nil {
		return "", fmt.Errorf("inmemory.CreateShortURL - db.Update: %v", err)
	}

	return shortURL, nil
}

func (s *Storage) GetFullURL(ctx context.Context, shortURL string) (string, error) {
	key := []byte(shortURL)
	var v []byte
	err := s.db.View(func(txn *badger.Txn) error {
		i, err := txn.Get(key)
		if err != nil {
			return err
		}
		v, err = i.ValueCopy(v)
		return err
	})
	if err != nil {
		return "", fmt.Errorf("inmemory.GetFullURL - db.View: %v", err)
	}

	return string(v), nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
