package service

import (
	"context"
	"path/filepath"
	"testing"

	"inventoryseal/internal/domain"
	"inventoryseal/internal/store"
)

func newService(t *testing.T) *Service {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "data.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return New(db)
}

func seedBatch(t *testing.T, s *Service, id string) {
	t.Helper()
	if err := s.CreateBatch(context.Background(), domain.Batch{ID: id, Title: "sealed", Owner: "ops", CreatedBy: "registrar"}); err != nil {
		t.Fatal(err)
	}
	if err := s.AddRecord(context.Background(), domain.Record{ID: id + "-r1", BatchID: id, Label: "A", Expected: 10, Observed: 10, UpdatedBy: "registrar"}); err != nil {
		t.Fatal(err)
	}
}

func TestCreateReview(t *testing.T) {
	s := newService(t)
	seedBatch(t, s, "b1")
	if err := s.ReviewBatch(context.Background(), "b1", "reviewer"); err != nil {
		t.Fatal(err)
	}
	b, err := s.DB().GetBatch("b1")
	if err != nil || b.Status != domain.BatchReview {
		t.Fatalf("batch %#v err %v", b, err)
	}
}
