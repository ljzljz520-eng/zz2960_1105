package store

import (
	"encoding/json"
	"go.etcd.io/bbolt"
	"inventoryseal/internal/domain"
)

type Transaction struct{ tx *bbolt.Tx }

func (db *DB) ReadTransaction(fn func(Transaction) error) error {
	return db.bolt.View(func(tx *bbolt.Tx) error { return fn(Transaction{tx: tx}) })
}
func (db *DB) WriteTransaction(fn func(Transaction) error) error {
	return db.bolt.Update(func(tx *bbolt.Tx) error { return fn(Transaction{tx: tx}) })
}
func (t Transaction) PutBatch(batch domain.Batch) error {
	data, err := json.Marshal(batch)
	if err != nil {
		return err
	}
	return t.tx.Bucket([]byte("batches")).Put([]byte(batch.ID), data)
}
func (t Transaction) PutRecord(record domain.Record) error {
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}
	return t.tx.Bucket([]byte("records")).Put([]byte(record.ID), data)
}
func (t Transaction) DeleteRecord(id string) error {
	return t.tx.Bucket([]byte("records")).Delete([]byte(id))
}
func (t Transaction) ClearBucket(name string) error {
	bucket := t.tx.Bucket([]byte(name))
	var keys [][]byte
	if err := bucket.ForEach(func(key, _ []byte) error { keys = append(keys, append([]byte(nil), key...)); return nil }); err != nil {
		return err
	}
	for _, key := range keys {
		if err := bucket.Delete(key); err != nil {
			return err
		}
	}
	return nil
}
