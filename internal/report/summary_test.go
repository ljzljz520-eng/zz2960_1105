package report

import (
	"inventoryseal/internal/domain"
	"testing"
)

func TestResultText(t *testing.T) {
	if ResultText(domain.Record{Result: "match"}) == "" || StatusText(domain.BatchArchived) == "" {
		t.Fatal("missing text")
	}
}
