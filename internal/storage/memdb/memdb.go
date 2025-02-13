package memdb

import (
	"github.com/dgraph-io/badger"
)

func NewMemDB() (*badger.DB, error) {
	opts := badger.DefaultOptions("./internal/storage/memdb")
	//opts.Dir = "" // Указываем путь к in-memory хранилищу
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return db, nil
}
