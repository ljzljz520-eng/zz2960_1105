package reconcile

import (
	"strings"
	"testing"

	"inventoryseal/internal/domain"
)

func TestAnalyzeFindsStaleAndDuplicateRecords(t *testing.T) {
	batch := domain.Batch{ID: "r-1", Status: domain.BatchReview, Version: 2, RecordCount: 2}
	records := []domain.Record{{ID: "b", BatchID: "r-1", Label: "same", Expected: 2, Observed: 1, Result: "match"}, {ID: "a", BatchID: "r-1", Label: "same", Expected: 2, Observed: 2, Result: "match"}}
	report := Analyze(batch, records)
	if report.Ready || !HasCode(report, "stale_result") || !HasCode(report, "duplicate_label") {
		t.Fatalf("unexpected report: %+v", report)
	}
	if len(Diff(report, Compact(report))) != 0 {
		t.Fatal("compacting findings changed report state")
	}
}

func TestFingerprintIsStableAndFormatIsUseful(t *testing.T) {
	batch := domain.Batch{ID: "r-2", Status: domain.BatchReview, Version: 1, RecordCount: 1}
	record := domain.Record{ID: "a", BatchID: "r-2", Label: "A", Expected: 4, Observed: 4, Result: "match"}
	left := Analyze(batch, []domain.Record{record})
	right := Analyze(batch, []domain.Record{record})
	if !SameState(left, right) || Fingerprint(batch, []domain.Record{record}) == "" {
		t.Fatal("fingerprint was not stable")
	}
	if !strings.Contains(Format(left), "ready=true") || !strings.Contains(Describe(left), "ready") {
		t.Fatalf("unexpected report text: %s", Format(left))
	}
}
