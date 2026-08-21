package service

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"inventoryseal/internal/domain"
	"inventoryseal/internal/store"
)

func TestWorkflowCreateReviewArchive(t *testing.T) {
	s := newService(t)
	seedBatch(t, s, "wf-archive")
	if err := s.ReviewBatch(context.Background(), "wf-archive", "auditor"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.ConfirmAndPublish(context.Background(), "wf-archive", "auditor"); err != nil {
		t.Fatal(err)
	}
	if err := s.ArchiveBatch(context.Background(), "wf-archive", "archivist"); err != nil {
		t.Fatal(err)
	}
	batch, err := s.DB().GetBatch("wf-archive")
	if err != nil || batch.Status != domain.BatchArchived {
		t.Fatalf("batch %#v err %v", batch, err)
	}
}

func TestWorkflowSearchUpdatePublish(t *testing.T) {
	s := newService(t)
	seedBatch(t, s, "wf-search")
	items, err := s.FindBatches(context.Background(), "search", "")
	if err != nil || len(items) != 1 {
		t.Fatalf("items %#v err %v", items, err)
	}
	record, err := s.DB().GetRecord("wf-search-r1")
	if err != nil {
		t.Fatal(err)
	}
	record.Observed = 12
	if err := s.UpdateRecord(context.Background(), record, "operator"); err != nil {
		t.Fatal(err)
	}
	if err := s.ReviewBatch(context.Background(), "wf-search", "reviewer"); err != nil {
		t.Fatal(err)
	}
	export, err := s.ConfirmAndPublish(context.Background(), "wf-search", "publisher")
	if err != nil || export.RecordCount != 1 {
		t.Fatalf("export %#v err %v", export, err)
	}
}

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "reopen.db")
	db, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	s := New(db)
	seedBatch(t, s, "reopen")
	if err := s.AddNote(context.Background(), domain.CollaborationNote{BatchID: "reopen", Author: "a", Message: "checked"}); err != nil {
		t.Fatal(err)
	}
	if err := db.SaveSnapshot(domain.ExportSnapshot{ID: "snap", BatchID: "reopen", Digest: "d", Payload: "p"}); err != nil {
		t.Fatal(err)
	}
	if err := db.SaveAudit(domain.AuditEvent{ID: "audit", BatchID: "reopen", Action: "create"}); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	db, err = store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.GetBatch("reopen"); err != nil {
		t.Fatal(err)
	}
	if _, err := db.GetRecord("reopen-r1"); err != nil {
		t.Fatal(err)
	}
	if _, err := db.GetSnapshot("snap"); err != nil {
		t.Fatal(err)
	}
	if notes, err := db.ListNotes("reopen"); err != nil || len(notes) != 1 {
		t.Fatalf("notes %#v err %v", notes, err)
	}
	if audits, err := db.ListAudits("reopen"); err != nil || len(audits) < 3 {
		t.Fatalf("audits %#v err %v", audits, err)
	}
}

func TestBusiness11Regression(t *testing.T) {
	s := newService(t)
	seedBatch(t, s, "2960-11")
	if err := s.ReviewBatch(context.Background(), "2960-11", "operator"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.ConfirmAndPublish(context.Background(), "2960-11", "operator"); err != nil {
		t.Fatal(err)
	}
	if err := s.ReopenBatch(context.Background(), "2960-11", "operator"); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := s.ConfirmAndPublish(ctx, "2960-11", "operator"); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected cancellation, got %v", err)
	}
	batch, err := s.DB().GetBatch("2960-11")
	if err != nil {
		t.Fatal(err)
	}
	if batch.Status != domain.BatchReview {
		t.Fatalf("canceled operation changed status to %s", batch.Status)
	}
}
