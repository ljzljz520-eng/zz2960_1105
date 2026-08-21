package query

import (
	"inventoryseal/internal/domain"
	"strings"
)

func MatchRecord(record domain.Record, term string) bool {
	needle := strings.ToLower(strings.TrimSpace(term))
	if needle == "" {
		return true
	}
	return strings.Contains(strings.ToLower(record.ID+" "+record.Label+" "+record.Result), needle)
}

func FilterRecords(records []domain.Record, term string, confirmed *bool) []domain.Record {
	result := make([]domain.Record, 0)
	for _, record := range records {
		if confirmed != nil && record.Confirmed != *confirmed {
			continue
		}
		if MatchRecord(record, term) {
			result = append(result, record)
		}
	}
	return domain.SortRecords(result)
}

func GroupByResult(records []domain.Record) map[string][]domain.Record {
	grouped := map[string][]domain.Record{}
	for _, record := range records {
		grouped[record.Result] = append(grouped[record.Result], record)
	}
	return grouped
}
