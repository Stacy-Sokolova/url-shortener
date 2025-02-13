package memdb

import (
	"github.com/dgraph-io/badger"
)

func NewMemDB() (*badger.DB, error) {
	opts := badger.DefaultOptions("./tmp")
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return db, nil
}
