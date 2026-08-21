package reconcile

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"inventoryseal/internal/domain"
)

type Severity string

const (
	SeverityInfo  Severity = "info"
	SeverityWarn  Severity = "warning"
	SeverityError Severity = "error"
)

type Finding struct {
	Code     string   `json:"code"`
	Severity Severity `json:"severity"`
	RecordID string   `json:"record_id,omitempty"`
	Detail   string   `json:"detail"`
}

type Counts struct {
	Total     int `json:"total"`
	Matches   int `json:"matches"`
	Overages  int `json:"overages"`
	Shortages int `json:"shortages"`
}

type Report struct {
	BatchID     string    `json:"batch_id"`
	Status      string    `json:"status"`
	Counts      Counts    `json:"counts"`
	Findings    []Finding `json:"findings"`
	Fingerprint string    `json:"fingerprint"`
	Ready       bool      `json:"ready"`
}

func Analyze(batch domain.Batch, records []domain.Record) Report {
	ordered := CanonicalRecords(records)
	report := Report{BatchID: batch.ID, Status: string(batch.Status), Counts: Count(ordered), Findings: make([]Finding, 0), Fingerprint: Fingerprint(batch, ordered)}
	seenIDs := make(map[string]bool, len(ordered))
	seenLabels := make(map[string]string, len(ordered))
	for _, record := range ordered {
		if seenIDs[record.ID] {
			report.Findings = append(report.Findings, Finding{"duplicate_id", SeverityError, record.ID, "record id occurs more than once"})
		}
		seenIDs[record.ID] = true
		label := strings.ToLower(strings.TrimSpace(record.Label))
		if previous, ok := seenLabels[label]; ok && previous != record.ID {
			report.Findings = append(report.Findings, Finding{"duplicate_label", SeverityError, record.ID, fmt.Sprintf("label also belongs to %s", previous)})
		}
		seenLabels[label] = record.ID
		if err := domain.ValidateRecord(record); err != nil {
			report.Findings = append(report.Findings, Finding{"invalid_record", SeverityError, record.ID, err.Error()})
			continue
		}
		if record.Result != domain.EvaluateRecord(record) {
			report.Findings = append(report.Findings, Finding{"stale_result", SeverityError, record.ID, "stored evaluation differs from observed counts"})
		}
		if record.Confirmed && !record.Published {
			report.Findings = append(report.Findings, Finding{"unpublished_confirmation", SeverityWarn, record.ID, "record was confirmed but not published"})
		}
	}
	if batch.RecordCount != len(ordered) {
		report.Findings = append(report.Findings, Finding{"record_count_mismatch", SeverityError, "", fmt.Sprintf("batch declares %d records but stores %d", batch.RecordCount, len(ordered))})
	}
	if batch.Status == domain.BatchPublished && len(ordered) == 0 {
		report.Findings = append(report.Findings, Finding{"empty_publication", SeverityError, "", "published batch has no records"})
	}
	report.Ready = Ready(report)
	return report
}

func Count(records []domain.Record) Counts {
	counts := Counts{Total: len(records)}
	for _, record := range records {
		switch domain.EvaluateRecord(record) {
		case "match":
			counts.Matches++
		case "overage":
			counts.Overages++
		case "shortage":
			counts.Shortages++
		}
	}
	return counts
}

func Ready(report Report) bool {
	if report.Counts.Total == 0 {
		return false
	}
	for _, finding := range report.Findings {
		if finding.Severity == SeverityError {
			return false
		}
	}
	return true
}

func CanonicalRecords(records []domain.Record) []domain.Record {
	ordered := append([]domain.Record(nil), records...)
	sort.SliceStable(ordered, func(i, j int) bool {
		if ordered[i].ID == ordered[j].ID {
			return ordered[i].Version < ordered[j].Version
		}
		return ordered[i].ID < ordered[j].ID
	})
	return ordered
}

func Fingerprint(batch domain.Batch, records []domain.Record) string {
	ordered := CanonicalRecords(records)
	var b strings.Builder
	fmt.Fprintf(&b, "%s|%s|%d|%d|", batch.ID, batch.Status, batch.Version, batch.RecordCount)
	for _, record := range ordered {
		fmt.Fprintf(&b, "%s:%s:%d:%d:%s:%t:%t:%d|", record.ID, record.Label, record.Expected, record.Observed, record.Result, record.Confirmed, record.Published, record.Version)
	}
	digest := sha256.Sum256([]byte(b.String()))
	return hex.EncodeToString(digest[:])
}

func SameState(left, right Report) bool {
	return left.BatchID == right.BatchID && left.Fingerprint == right.Fingerprint && left.Ready == right.Ready
}

func SeverityRank(value Severity) int {
	switch value {
	case SeverityError:
		return 3
	case SeverityWarn:
		return 2
	case SeverityInfo:
		return 1
	default:
		return 0
	}
}

func HighestSeverity(findings []Finding) Severity {
	highest := SeverityInfo
	for _, finding := range findings {
		if SeverityRank(finding.Severity) > SeverityRank(highest) {
			highest = finding.Severity
		}
	}
	return highest
}

func Codes(findings []Finding) []string {
	codes := make([]string, 0, len(findings))
	for _, finding := range findings {
		codes = append(codes, finding.Code)
	}
	sort.Strings(codes)
	return codes
}

func Format(report Report) string {
	var b strings.Builder
	fmt.Fprintf(&b, "batch=%s status=%s total=%d matches=%d overages=%d shortages=%d ready=%t fingerprint=%s", report.BatchID, report.Status, report.Counts.Total, report.Counts.Matches, report.Counts.Overages, report.Counts.Shortages, report.Ready, report.Fingerprint)
	for _, finding := range report.Findings {
		fmt.Fprintf(&b, "\n%s [%s] %s", finding.Code, finding.Severity, finding.Detail)
	}
	return b.String()
}
