package data

import (
	"github.com/boltdb/bolt"
)

type Data struct {
	db *bolt.DB
}

func (d *Data) Close() {
	d.db.Close()
}

func New(fn string) (*Data, error) {
	db, err := bolt.Open(fn, 0600, nil)
	if err != nil {
		return nil, err
	}

	return &Data{
		db: db,
	}, nil
}
