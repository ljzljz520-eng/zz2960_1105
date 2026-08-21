package report

import (
	"fmt"
	"strings"

	"inventoryseal/internal/domain"
	"inventoryseal/internal/ledger"
)

func Markdown(batch domain.Batch, totals ledger.Totals) string {
	var lines []string
	lines = append(lines, "# "+batch.Title, "", "- Batch: "+batch.ID, "- Status: "+string(batch.Status), fmt.Sprintf("- Records: %d", totals.Total), fmt.Sprintf("- Match: %d", totals.Match), fmt.Sprintf("- Overage: %d", totals.Overage), fmt.Sprintf("- Shortage: %d", totals.Shortage))
	return strings.Join(lines, "\n")
}

func Plain(batch domain.Batch, records []domain.Record) string {
	lines := []string{batch.ID + " [" + string(batch.Status) + "]"}
	for _, record := range records {
		lines = append(lines, fmt.Sprintf("%s %s %d/%d %s", record.ID, record.Label, record.Expected, record.Observed, record.Result))
	}
	return strings.Join(lines, "\n")
}

func StableDigestInput(batch domain.Batch, records []domain.Record) string {
	items := domain.SortRecords(records)
	parts := []string{batch.ID, string(batch.Status)}
	for _, record := range items {
		parts = append(parts, record.ID, record.Result, fmt.Sprint(record.Expected), fmt.Sprint(record.Observed))
	}
	return strings.Join(parts, "|")
}
