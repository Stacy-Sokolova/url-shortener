package storage

import (
	"context"
	pb "url-server/internal/service/proto"
	"url-server/internal/storage/memdb"
	"url-server/internal/storage/pgdb"

	"github.com/dgraph-io/badger"
	"github.com/jmoiron/sqlx"
)

type Storage interface {
	GetFullURL(ctx context.Context, r *pb.Request) (*pb.Response, error)
	CreateShortURL(ctx context.Context, r *pb.Request) (*pb.Response, error)
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
