package memdb

import (
	"context"

	pb "url-server/internal/service/proto"
	fnc "url-server/internal/storage/data_processing"

	"github.com/dgraph-io/badger"
)

type Storage struct {
	db *badger.DB
}

func NewInmemStorage(db *badger.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) CreateShortURL(ctx context.Context, r *pb.Request) (*pb.Response, error) {
	fullURL := r.GetUrl()
	shortURL := fnc.URLShortener(fullURL)
	val := []byte(fullURL)
	key := []byte(shortURL)
	err := s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, val)
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

func (s *Storage) Close() error {
	return s.db.Close()
}
