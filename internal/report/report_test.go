package report

import (
	"inventoryseal/internal/domain"
	"testing"
)

func TestBuildReport(t *testing.T) {
	value := Build(domain.Batch{ID: "b", Status: domain.BatchPublished}, []domain.Record{{Result: "match"}, {Result: "shortage"}}, "abc")
	if value.Total != 2 || value.Matches != 1 || value.Shortages != 1 {
		t.Fatalf("report %#v", value)
	}
}
