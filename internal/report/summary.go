package report

import (
	"fmt"
	"inventoryseal/internal/domain"
	"sort"
	"strings"
)

type Line struct {
	Key   string
	Value string
}

func Lines(batch domain.Batch, records []domain.Record) []Line {
	totals := Build(batch, records, "")
	return []Line{{"batch", batch.ID}, {"status", string(batch.Status)}, {"total", fmt.Sprint(totals.Total)}, {"matches", fmt.Sprint(totals.Matches)}, {"overages", fmt.Sprint(totals.Overages)}, {"shortages", fmt.Sprint(totals.Shortages)}}
}
func KeyValue(batch domain.Batch, records []domain.Record) string {
	lines := Lines(batch, records)
	sort.Slice(lines, func(i, j int) bool { return lines[i].Key < lines[j].Key })
	values := make([]string, 0, len(lines))
	for _, line := range lines {
		values = append(values, line.Key+"="+line.Value)
	}
	return strings.Join(values, "\n")
}
func ResultText(record domain.Record) string {
	switch record.Result {
	case "match":
		return "数量一致"
	case "overage":
		return "盘盈"
	case "shortage":
		return "盘亏"
	default:
		return "待评估"
	}
}
func StatusText(status domain.BatchStatus) string {
	switch status {
	case domain.BatchDraft:
		return "草稿"
	case domain.BatchReview:
		return "审核中"
	case domain.BatchConfirmed:
		return "已确认"
	case domain.BatchPublished:
		return "已发布"
	case domain.BatchArchived:
		return "已归档"
	default:
		return "未知"
	}
}
