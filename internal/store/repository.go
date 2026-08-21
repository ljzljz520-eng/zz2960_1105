package store

import (
	"encoding/json"
	"fmt"
	"sort"

	"go.etcd.io/bbolt"
	"inventoryseal/internal/domain"
)

func put[T any](db *DB, bucket, key string, value T) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return db.bolt.Update(func(tx *bbolt.Tx) error { return tx.Bucket([]byte(bucket)).Put([]byte(key), data) })
}

func get[T any](db *DB, bucket, key string) (T, error) {
	var value T
	err := db.bolt.View(func(tx *bbolt.Tx) error {
		data := tx.Bucket([]byte(bucket)).Get([]byte(key))
		if data == nil {
			return ErrNotFound
		}
		return json.Unmarshal(append([]byte(nil), data...), &value)
	})
	return value, err
}

func all[T any](db *DB, bucket string) ([]T, error) {
	values := make([]T, 0)
	err := db.bolt.View(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte(bucket)).ForEach(func(_, data []byte) error {
			if data == nil {
				return nil
			}
			var value T
			if err := json.Unmarshal(data, &value); err != nil {
				return err
			}
			values = append(values, value)
			return nil
		})
	})
	sort.Slice(values, func(i, j int) bool { return fmt.Sprint(values[i]) < fmt.Sprint(values[j]) })
	return values, err
}

func (db *DB) SaveBatch(batch domain.Batch) error       { return put(db, "batches", batch.ID, batch) }
func (db *DB) GetBatch(id string) (domain.Batch, error) { return get[domain.Batch](db, "batches", id) }
func (db *DB) ListBatches() ([]domain.Batch, error)     { return all[domain.Batch](db, "batches") }

func (db *DB) SaveRecord(record domain.Record) error { return put(db, "records", record.ID, record) }
func (db *DB) GetRecord(id string) (domain.Record, error) {
	return get[domain.Record](db, "records", id)
}
func (db *DB) ListRecords(batchID string) ([]domain.Record, error) {
	items, err := all[domain.Record](db, "records")
	if err != nil {
		return nil, err
	}
	filtered := items[:0]
	for _, item := range items {
		if batchID == "" || item.BatchID == batchID {
			filtered = append(filtered, item)
		}
	}
	return domain.SortRecords(filtered), nil
}

func (db *DB) SaveAudit(event domain.AuditEvent) error { return put(db, "audits", event.ID, event) }
func (db *DB) ListAudits(batchID string) ([]domain.AuditEvent, error) {
	items, err := all[domain.AuditEvent](db, "audits")
	if err != nil {
		return nil, err
	}
	filtered := items[:0]
	for _, item := range items {
		if batchID == "" || item.BatchID == batchID {
			filtered = append(filtered, item)
		}
	}
	sort.Slice(filtered, func(i, j int) bool { return filtered[i].Sequence < filtered[j].Sequence })
	return filtered, nil
}

func (db *DB) SaveNote(note domain.CollaborationNote) error { return put(db, "notes", note.ID, note) }
func (db *DB) ListNotes(batchID string) ([]domain.CollaborationNote, error) {
	items, err := all[domain.CollaborationNote](db, "notes")
	if err != nil {
		return nil, err
	}
	filtered := items[:0]
	for _, item := range items {
		if batchID == "" || item.BatchID == batchID {
			filtered = append(filtered, item)
		}
	}
	sort.Slice(filtered, func(i, j int) bool { return filtered[i].Sequence < filtered[j].Sequence })
	return filtered, nil
}

func (db *DB) SaveSnapshot(snapshot domain.ExportSnapshot) error {
	return put(db, "snapshots", snapshot.ID, snapshot)
}
func (db *DB) GetSnapshot(id string) (domain.ExportSnapshot, error) {
	return get[domain.ExportSnapshot](db, "snapshots", id)
}
func (db *DB) ListSnapshots(batchID string) ([]domain.ExportSnapshot, error) {
	items, err := all[domain.ExportSnapshot](db, "snapshots")
	if err != nil {
		return nil, err
	}
	filtered := items[:0]
	for _, item := range items {
		if batchID == "" || item.BatchID == batchID {
			filtered = append(filtered, item)
		}
	}
	sort.Slice(filtered, func(i, j int) bool { return filtered[i].Sequence < filtered[j].Sequence })
	return filtered, nil
}
