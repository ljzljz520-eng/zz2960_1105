package service

import (
	"context"
	"testing"

	"inventoryseal/internal/domain"
)

// TestBusiness11NoPartialStateOnCancel asserts the deeper guarantee of the
// 2960-11 regression fix: a canceled ConfirmAndPublish must not leave the
// batch or its records in an intermediate confirmed/published state, so that
// any subsequent confirm-and-publish observes the previously published,
// self-contained business result instead of inheriting a corrupt partial write.
func TestBusiness11NoPartialStateOnCancel(t *testing.T) {
	s := newService(t)
	seedBatch(t, s, "2960-11")
	if err := s.ReviewBatch(context.Background(), "2960-11", "operator"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.ConfirmAndPublish(context.Background(), "2960-11", "operator"); err != nil {
		t.Fatal(err)
	}
	first, err := s.DB().GetBatch("2960-11")
	if err != nil {
		t.Fatal(err)
	}
	if first.Status != domain.BatchPublished {
		t.Fatalf("expected published, got %s", first.Status)
	}
	firstRecord, err := s.DB().GetRecord("2960-11-r1")
	if err != nil {
		t.Fatal(err)
	}
	if !firstRecord.Confirmed || !firstRecord.Published {
		t.Fatalf("first publish should confirm and publish record: %+v", firstRecord)
	}
	firstExport, err := s.LatestExport("2960-11")
	if err != nil {
		t.Fatal(err)
	}

	// Reopen for a correction cycle, then cancel the new publication attempt.
	if err := s.ReopenBatch(context.Background(), "2960-11", "operator"); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := s.ConfirmAndPublish(ctx, "2960-11", "operator"); err == nil {
		t.Fatal("expected cancellation error from canceled publish")
	}

	// The canceled attempt must leave no partial state behind.
	batch, err := s.DB().GetBatch("2960-11")
	if err != nil {
		t.Fatal(err)
	}
	if batch.Status != domain.BatchReview {
		t.Fatalf("canceled operation changed status to %s", batch.Status)
	}
	record, err := s.DB().GetRecord("2960-11-r1")
	if err != nil {
		t.Fatal(err)
	}
	// Records retain the true confirmed/published result from the first
	// successful publication; the canceled attempt must not flip them into an
	// intermediate confirmed-without-publish or revert them.
	if record.Confirmed != firstRecord.Confirmed || record.Published != firstRecord.Published {
		t.Fatalf("canceled publish mutated record flags to confirmed=%v published=%v", record.Confirmed, record.Published)
	}
	// No new export snapshot should have been produced by the canceled attempt.
	exports, err := s.DB().ListSnapshots("2960-11")
	if err != nil {
		t.Fatal(err)
	}
	if len(exports) != 1 || exports[0].ID != firstExport.ID {
		t.Fatalf("canceled publish should not add a snapshot: %+v", exports)
	}

	// A subsequent publish must succeed from the independent review state and
	// produce a fresh, self-contained export.
	if _, err := s.ConfirmAndPublish(context.Background(), "2960-11", "operator"); err != nil {
		t.Fatalf("subsequent publish failed: %v", err)
	}
	batch, err = s.DB().GetBatch("2960-11")
	if err != nil {
		t.Fatal(err)
	}
	if batch.Status != domain.BatchPublished {
		t.Fatalf("expected published after retry, got %s", batch.Status)
	}
	exports, err = s.DB().ListSnapshots("2960-11")
	if err != nil {
		t.Fatal(err)
	}
	if len(exports) != 2 {
		t.Fatalf("expected two independent exports, got %d", len(exports))
	}
	if exports[0].ID == exports[1].ID {
		t.Fatalf("exports are not independent: %s", exports[0].ID)
	}
}
