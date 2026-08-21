package report

import (
	"encoding/json"
	"fmt"
	"sort"

	"inventoryseal/internal/domain"
)

type BatchReport struct {
	BatchID   string             `json:"batch_id"`
	Status    domain.BatchStatus `json:"status"`
	Total     int                `json:"total"`
	Matches   int                `json:"matches"`
	Overages  int                `json:"overages"`
	Shortages int                `json:"shortages"`
	Digest    string             `json:"digest"`
}

func Build(batch domain.Batch, records []domain.Record, digest string) BatchReport {
	matches, overages, shortages := domain.Summary(records)
	return BatchReport{BatchID: batch.ID, Status: batch.Status, Total: len(records), Matches: matches, Overages: overages, Shortages: shortages, Digest: digest}
}

func JSON(value BatchReport) ([]byte, error) { return json.MarshalIndent(value, "", "  ") }

func CSV(records []domain.Record) string {
	items := domain.SortRecords(records)
	lines := []string{"id,label,expected,observed,result,confirmed,published"}
	for _, item := range items {
		lines = append(lines, fmt.Sprintf("%s,%s,%d,%d,%s,%t,%t", item.ID, item.Label, item.Expected, item.Observed, item.Result, item.Confirmed, item.Published))
	}
	return fmt.Sprintln(sort.StringSlice(lines))
}
