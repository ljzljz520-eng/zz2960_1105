package store

import (
	"fmt"
	"time"

	"go.etcd.io/bbolt"
)

type Health struct {
	Path      string
	Open      bool
	Buckets   int
	Timestamp string
}

func (db *DB) Health() Health {
	open := db != nil && db.bolt != nil
	return Health{Path: db.path, Open: open, Buckets: len(bucketNames), Timestamp: time.Unix(0, 0).UTC().Format(time.RFC3339)}
}
func (db *DB) String() string {
	health := db.Health()
	return fmt.Sprintf("%s open=%t buckets=%d", health.Path, health.Open, health.Buckets)
}
func (db *DB) BucketNames() []string {
	names := make([]string, 0, len(bucketNames))
	for _, name := range bucketNames {
		names = append(names, string(name))
	}
	return names
}
func (db *DB) CountBucket(name string) (int, error) {
	count := 0
	err := db.bolt.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(name))
		if bucket == nil {
			return ErrNotFound
		}
		return bucket.ForEach(func(_, value []byte) error {
			if value != nil {
				count++
			}
			return nil
		})
	})
	return count, err
}
