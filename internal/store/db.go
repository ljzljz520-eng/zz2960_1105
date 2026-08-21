package store

import (
	"errors"
	"os"

	"go.etcd.io/bbolt"
)

var ErrNotFound = errors.New("entity not found")
var ErrConflict = errors.New("version conflict")

var bucketNames = [][]byte{[]byte("batches"), []byte("records"), []byte("audits"), []byte("notes"), []byte("snapshots"), []byte("meta")}

type DB struct {
	bolt *bbolt.DB
	path string
}

func Open(path string) (*DB, error) {
	bolt, err := bbolt.Open(path, 0o600, &bbolt.Options{NoSync: true})
	if err != nil {
		return nil, err
	}
	db := &DB{bolt: bolt, path: path}
	if err := db.initialize(); err != nil {
		bolt.Close()
		return nil, err
	}
	return db, nil
}

func (db *DB) initialize() error {
	return db.bolt.Update(func(tx *bbolt.Tx) error {
		for _, name := range bucketNames {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return err
			}
		}
		return nil
	})
}

func (db *DB) Close() error {
	if db == nil || db.bolt == nil {
		return nil
	}
	return db.bolt.Close()
}
func (db *DB) Path() string { return db.path }
func (db *DB) Exists() bool { _, err := os.Stat(db.path); return err == nil }
