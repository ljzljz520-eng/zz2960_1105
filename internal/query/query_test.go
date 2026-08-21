package query

import (
	"inventoryseal/internal/domain"
	"testing"
)

func TestFilterRecords(t *testing.T) {
	records := []domain.Record{{ID: "1", Label: "cold", Result: "match"}, {ID: "2", Label: "warm", Result: "shortage"}}
	got := FilterRecords(records, "cold", nil)
	if len(got) != 1 || got[0].ID != "1" {
		t.Fatalf("got %#v", got)
	}
}
