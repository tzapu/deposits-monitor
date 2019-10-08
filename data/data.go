package data

import (
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("module", "data")

type Data struct {
	DB *bolt.DB
}

func (d *Data) Stats(bucket string) (bolt.BucketStats, error) {
	var stats bolt.BucketStats
	err := d.DB.View(func(tx *bolt.Tx) error {
		stats = tx.Bucket([]byte(bucket)).Stats()
		return nil
	})
	return stats, err
}

func (d *Data) Count(bucket string) (int, error) {
	stats, err := d.Stats(bucket)
	return stats.KeyN, err
}

// Get value from bucket by key
func (d *Data) Get(bucket string, key string) ([]byte, error) {
	var value []byte
	err := d.DB.View(func(tx *bolt.Tx) error {
		v := tx.Bucket([]byte(bucket)).Get([]byte(key))
		if v != nil {
			value = make([]byte, len(v))
			copy(value, v)
		}
		return nil
	})
	return value, err
}

// Put a key/value pair into target bucket
func (d *Data) Put(bucket string, key string, value []byte) error {
	err := d.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Put([]byte(key), value)
		return err
	})

	return err
}

// Encode data as JSON and put it into the target bucket under key
func (d *Data) PutStruct(bucket string, key string, data interface{}) error {
	// Put a key/value pair into target bucket
	err := d.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		value, err := json.Marshal(data)
		if err != nil {
			return err
		}
		return b.Put([]byte(key), value)
	})
	return err
}

func (d *Data) String(bucket string, key string) (string, error) {
	v, err := d.Get(bucket, key)
	return string(v), err
}

func (d *Data) PutString(bucket string, key string, value string) error {
	return d.Put(bucket, key, []byte(value))
}

func (d *Data) Bool(bucket string, key string) (bool, error) {
	v, err := d.Get(bucket, key)
	if err != nil {
		return false, err
	}
	if len(v) == 1 && v[0] == 1 {
		return true, nil
	}

	return false, nil
}

func (d *Data) PutBool(bucket string, key string, value bool) error {
	var b byte = 0
	if value {
		b = 1
	}
	return d.Put(bucket, key, []byte{b})
}

// Close the database connection
func (d *Data) Close() error {
	return d.DB.Close()
}

// New returns a new BoltDB connection
func New(fn string, buckets []string) (*Data, error) {
	db, err := bolt.Open(fn, 0600, &bolt.Options{Timeout: 15 * time.Second})
	if err != nil {
		return nil, err
	}

	// making sure all buckets exist
	err = db.Update(func(tx *bolt.Tx) error {
		for _, value := range buckets {
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
		DB: db,
	}, nil
}
