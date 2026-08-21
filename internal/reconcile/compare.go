package reconcile

import (
	"fmt"
	"sort"
	"strings"
)

type Difference struct {
	Field  string `json:"field"`
	Before string `json:"before"`
	After  string `json:"after"`
}

func Diff(left, right Report) []Difference {
	differences := make([]Difference, 0)
	if left.BatchID != right.BatchID {
		differences = append(differences, Difference{"batch_id", left.BatchID, right.BatchID})
	}
	if left.Status != right.Status {
		differences = append(differences, Difference{"status", left.Status, right.Status})
	}
	if left.Fingerprint != right.Fingerprint {
		differences = append(differences, Difference{"fingerprint", left.Fingerprint, right.Fingerprint})
	}
	if left.Ready != right.Ready {
		differences = append(differences, Difference{"ready", fmt.Sprint(left.Ready), fmt.Sprint(right.Ready)})
	}
	if left.Counts != right.Counts {
		differences = append(differences, Difference{"counts", fmt.Sprintf("%+v", left.Counts), fmt.Sprintf("%+v", right.Counts)})
	}
	return differences
}

func SortFindings(findings []Finding) []Finding {
	ordered := append([]Finding(nil), findings...)
	sort.SliceStable(ordered, func(i, j int) bool {
		if ordered[i].Code == ordered[j].Code {
			return ordered[i].RecordID < ordered[j].RecordID
		}
		return ordered[i].Code < ordered[j].Code
	})
	return ordered
}

func Compact(report Report) Report {
	report.Findings = SortFindings(report.Findings)
	return report
}

func HasCode(report Report, code string) bool {
	for _, finding := range report.Findings {
		if finding.Code == code {
			return true
		}
	}
	return false
}

func Describe(report Report) string {
	if report.Ready {
		return fmt.Sprintf("%s is ready with %d records", report.BatchID, report.Counts.Total)
	}
	severity := HighestSeverity(report.Findings)
	if len(report.Findings) == 0 {
		return fmt.Sprintf("%s is not ready because it has no records", report.BatchID)
	}
	return fmt.Sprintf("%s is blocked by %s findings: %s", report.BatchID, severity, strings.Join(Codes(report.Findings), ","))
}
