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

// NextAuditSequence returns the next audit sequence number for a batch, counting
// existing audit events for that batch. It is safe to call before opening a
// write transaction so the sequence can be reserved up front; it stays
// consistent with ListAudits because both filter by the event's BatchID.
func (db *DB) NextAuditSequence(batchID string) (int, error) {
	var seq int
	err := db.bolt.View(func(tx *bbolt.Tx) error {
		var err error
		seq, err = Transaction{tx: tx}.NextAuditSequence(batchID)
		return err
	})
	return seq, err
}

// NextSnapshotSequence returns the next export snapshot sequence number for a
// batch.
func (db *DB) NextSnapshotSequence(batchID string) (int, error) {
	var seq int
	err := db.bolt.View(func(tx *bbolt.Tx) error {
		var err error
		seq, err = Transaction{tx: tx}.NextSnapshotSequence(batchID)
		return err
	})
	return seq, err
}

func (t Transaction) PutBatch(batch domain.Batch) error {
	data, err := json.Marshal(batch)
	if err != nil {
		return err
	}
	return t.tx.Bucket([]byte("batches")).Put([]byte(batch.ID), data)
}
func (t Transaction) GetBatch(id string) (domain.Batch, error) {
	data := t.tx.Bucket([]byte("batches")).Get([]byte(id))
	if data == nil {
		return domain.Batch{}, ErrNotFound
	}
	var batch domain.Batch
	if err := json.Unmarshal(append([]byte(nil), data...), &batch); err != nil {
		return domain.Batch{}, err
	}
	return batch, nil
}
func (t Transaction) PutRecord(record domain.Record) error {
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}
	return t.tx.Bucket([]byte("records")).Put([]byte(record.ID), data)
}
func (t Transaction) GetRecord(id string) (domain.Record, error) {
	data := t.tx.Bucket([]byte("records")).Get([]byte(id))
	if data == nil {
		return domain.Record{}, ErrNotFound
	}
	var record domain.Record
	if err := json.Unmarshal(append([]byte(nil), data...), &record); err != nil {
		return domain.Record{}, err
	}
	return record, nil
}
func (t Transaction) ListRecords(batchID string) ([]domain.Record, error) {
	items := make([]domain.Record, 0)
	err := t.tx.Bucket([]byte("records")).ForEach(func(_, data []byte) error {
		if data == nil {
			return nil
		}
		var record domain.Record
		if err := json.Unmarshal(data, &record); err != nil {
			return err
		}
		if batchID == "" || record.BatchID == batchID {
			items = append(items, record)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return domain.SortRecords(items), nil
}
func (t Transaction) PutAudit(event domain.AuditEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return t.tx.Bucket([]byte("audits")).Put([]byte(event.ID), data)
}
func (t Transaction) NextAuditSequence(batchID string) (int, error) {
	count := 0
	err := t.tx.Bucket([]byte("audits")).ForEach(func(_, data []byte) error {
		if data == nil {
			return nil
		}
		var event domain.AuditEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return err
		}
		if event.BatchID == batchID {
			count++
		}
		return nil
	})
	return count + 1, err
}
func (t Transaction) PutSnapshot(snapshot domain.ExportSnapshot) error {
	data, err := json.Marshal(snapshot)
	if err != nil {
		return err
	}
	return t.tx.Bucket([]byte("snapshots")).Put([]byte(snapshot.ID), data)
}
func (t Transaction) NextSnapshotSequence(batchID string) (int, error) {
	count := 0
	err := t.tx.Bucket([]byte("snapshots")).ForEach(func(_, data []byte) error {
		if data == nil {
			return nil
		}
		var snapshot domain.ExportSnapshot
		if err := json.Unmarshal(data, &snapshot); err != nil {
			return err
		}
		if snapshot.BatchID == batchID {
			count++
		}
		return nil
	})
	return count + 1, err
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
