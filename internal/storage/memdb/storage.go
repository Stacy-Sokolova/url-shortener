package memdb

import (
	"context"

	pb "url-server/internal/service/proto"

	"github.com/dgraph-io/badger"
)

type Storage struct {
	db *badger.DB
}

func NewStorage(db *badger.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) CreateShortURL(ctx context.Context, r *pb.Request) (*pb.Response, error) {
	fullURL := r.GetUrl()
	shortURL := "generateFunc"
	val := []byte(fullURL)
	key := []byte(shortURL)
	err := s.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, val)
		return err
	})
	if err != nil {
		return nil, err
	}

	return &pb.Response{Url: shortURL}, nil
}

func (s *Storage) GetFullURL(ctx context.Context, r *pb.Request) (*pb.Response, error) {
	shortURL := r.GetUrl()
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
		return nil, err
	}

	return &pb.Response{Url: string(v)}, nil
}
