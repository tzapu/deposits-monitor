package data

import (
	"time"

	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("module", "data")

type Data struct {
	db *bolt.DB
}

func (d *Data) Close() error {
	return d.db.Close()
}

func New(fn string) (*Data, error) {
	db, err := bolt.Open(fn, 0600, &bolt.Options{Timeout: 15 * time.Second})
	if err != nil {
		return nil, err
	}

	// making sure all buckets exist
	err = db.Update(func(tx *bolt.Tx) error {
		for _, value := range Buckets {
			_, err := tx.CreateBucketIfNotExists([]byte(value))
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &Data{
		db: db,
	}, nil
}
