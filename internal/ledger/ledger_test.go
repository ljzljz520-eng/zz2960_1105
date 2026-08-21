package ledger

import (
	"inventoryseal/internal/domain"
	"testing"
)

func TestCount(t *testing.T) {
	value := Count([]domain.Record{{Result: "match", Confirmed: true}, {Result: "shortage"}})
	if value.Total != 2 || value.Match != 1 || value.Confirmed != 1 {
		t.Fatalf("totals %#v", value)
	}
}
