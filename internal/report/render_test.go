package report

import (
	"inventoryseal/internal/domain"
	"inventoryseal/internal/ledger"
	"testing"
)

func TestMarkdown(t *testing.T) {
	text := Markdown(domain.Batch{ID: "b", Title: "batch", Status: domain.BatchPublished}, ledger.Totals{Total: 1})
	if text == "" {
		t.Fatal("empty report")
	}
}
