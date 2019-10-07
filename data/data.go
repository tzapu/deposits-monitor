package data

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("module", "data")

type Data struct {
	db *bolt.DB
}

// Put a key/value pair into target bucket
func (d *Data) Put(bucket string, key string, value []byte) error {
	err := d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Put([]byte(key), value)
		return err
	})

	return err
}

// Encode data as JSON and put it into the target bucket under key
func (d *Data) PutStruct(bucket string, key string, data interface{}) error {
	// Put a key/value pair into target bucket
	err := d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		value, err := json.Marshal(data)
		if err != nil {
			return err
		}
		return b.Put([]byte(key), value)
	})
	return err
}

// Get value from bucket by key
func (d *Data) Get(bucket string, key string) ([]byte, error) {
	var value []byte
	err := d.db.View(func(tx *bolt.Tx) error {
		v := tx.Bucket([]byte(bucket)).Get([]byte(key))
		if v != nil {
			value = make([]byte, len(v))
			copy(value, v)
		}
		return nil
	})
	return value, err
}

func (d *Data) String(bucket string, key string) (string, error) {
	v, err := d.Get(bucket, key)
	return string(v), err
}

func (d *Data) SettingString(key string) (string, error) {
	return d.String(Settings, key)
}

func (d *Data) SettingBool(key string) (bool, error) {
	v, err := d.SettingString(key)
	if err != nil || v == "" {
		return false, err
	}
	b, err := strconv.ParseBool(v)

	return b, nil
}

// Close the database connection
func (d *Data) Close() error {
	return d.db.Close()
}

// New returns a new BoltDB connection
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
