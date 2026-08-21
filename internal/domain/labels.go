package domain

import (
	"sort"
	"strings"
)

func NormalizeLabel(value string) string { return strings.ToUpper(strings.TrimSpace(value)) }
func NormalizeLabels(records []Record) []Record {
	result := append([]Record(nil), records...)
	for index := range result {
		result[index].Label = NormalizeLabel(result[index].Label)
	}
	return result
}
func LabelIndex(records []Record) map[string][]string {
	result := map[string][]string{}
	for _, record := range records {
		key := NormalizeLabel(record.Label)
		result[key] = append(result[key], record.ID)
	}
	for key := range result {
		sort.Strings(result[key])
	}
	return result
}
func HasDuplicateLabels(records []Record) bool {
	index := map[string]bool{}
	for _, record := range records {
		key := NormalizeLabel(record.Label)
		if index[key] {
			return true
		}
		index[key] = true
	}
	return false
}
