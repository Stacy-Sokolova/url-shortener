package storage

import (
	"context"
	"url-server/internal/storage/memdb"
	"url-server/internal/storage/pgdb"

	"github.com/dgraph-io/badger"
	"github.com/jmoiron/sqlx"
)

type Storage interface {
	GetFullURL(ctx context.Context, url string) (string, error)
	CreateShortURL(ctx context.Context, url string) (string, error)
	Close() error
}

type MyStorage struct {
	Storage
}

func NewStorage(db any) *MyStorage {
	switch db := db.(type) {
	case *sqlx.DB:
		return &MyStorage{
			Storage: pgdb.NewSQLStorage(db),
		}
	case *badger.DB:
		return &MyStorage{
			Storage: memdb.NewInmemStorage(db),
		}
	}
	return nil
}
