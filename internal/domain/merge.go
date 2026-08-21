package domain

import "strings"

func MergeRecords(primary, secondary []Record) []Record {
	result := make(map[string]Record, len(primary)+len(secondary))
	for _, record := range primary {
		result[record.ID] = record
	}
	for _, record := range secondary {
		existing, ok := result[record.ID]
		if !ok || record.Version >= existing.Version {
			result[record.ID] = record
		}
	}
	merged := make([]Record, 0, len(result))
	for _, record := range result {
		merged = append(merged, record)
	}
	return SortRecords(merged)
}

func Labels(records []Record) []string {
	labels := make([]string, 0, len(records))
	seen := map[string]bool{}
	for _, record := range records {
		key := strings.TrimSpace(record.Label)
		if key == "" || seen[key] {
			continue
		}
		labels = append(labels, key)
		seen[key] = true
	}
	return labels
}

func ConfirmedOnly(records []Record) []Record {
	result := make([]Record, 0)
	for _, record := range records {
		if record.Confirmed {
			result = append(result, record)
		}
	}
	return SortRecords(result)
}

func PublishedOnly(records []Record) []Record {
	result := make([]Record, 0)
	for _, record := range records {
		if record.Published {
			result = append(result, record)
		}
	}
	return SortRecords(result)
}
