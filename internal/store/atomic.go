package store

import (
	"encoding/json"

	"go.etcd.io/bbolt"
	"inventoryseal/internal/domain"
)

func (db *DB) UpdateRecord(record domain.Record, expectedVersion int) error {
	return db.bolt.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("records"))
		existing := bucket.Get([]byte(record.ID))
		if existing == nil {
			return ErrNotFound
		}
		var current domain.Record
		if err := json.Unmarshal(existing, &current); err != nil {
			return err
		}
		if current.Version != expectedVersion {
			return ErrConflict
		}
		record.Version = expectedVersion + 1
		data, err := json.Marshal(record)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(record.ID), data)
	})
}

func (db *DB) SetMeta(key, value string) error {
	return db.bolt.Update(func(tx *bbolt.Tx) error { return tx.Bucket([]byte("meta")).Put([]byte(key), []byte(value)) })
}

func (db *DB) GetMeta(key string) (string, error) {
	var result string
	err := db.bolt.View(func(tx *bbolt.Tx) error {
		data := tx.Bucket([]byte("meta")).Get([]byte(key))
		if data == nil {
			return ErrNotFound
		}
		result = string(data)
		return nil
	})
	return result, err
}
