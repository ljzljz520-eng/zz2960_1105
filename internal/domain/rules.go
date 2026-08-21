package domain

import "fmt"

func ValidateBatch(batch Batch) error {
	if batch.ID == "" {
		return fmt.Errorf("batch id is required")
	}
	if batch.Title == "" {
		return fmt.Errorf("batch title is required")
	}
	if batch.Owner == "" {
		return fmt.Errorf("batch owner is required")
	}
	if batch.Version < 1 {
		return fmt.Errorf("batch version must be positive")
	}
	return nil
}

func ValidateRecord(record Record) error {
	if record.ID == "" || record.BatchID == "" {
		return fmt.Errorf("record identity is required")
	}
	if record.Label == "" {
		return fmt.Errorf("record label is required")
	}
	if record.Expected < 0 || record.Observed < 0 {
		return fmt.Errorf("record counts cannot be negative")
	}
	return nil
}

func EvaluateRecord(record Record) string {
	if record.Observed == record.Expected {
		return "match"
	}
	if record.Observed > record.Expected {
		return "overage"
	}
	return "shortage"
}

func CanTransition(from, to BatchStatus) bool {
	if from == BatchDraft && to == BatchReview {
		return true
	}
	if from == BatchReview && to == BatchConfirmed {
		return true
	}
	if from == BatchConfirmed && to == BatchPublished {
		return true
	}
	if from == BatchPublished && to == BatchArchived {
		return true
	}
	return false
}

func NormalizeStatus(status string) BatchStatus {
	if status == string(BatchReview) {
		return BatchReview
	}
	if status == string(BatchConfirmed) {
		return BatchConfirmed
	}
	if status == string(BatchPublished) {
		return BatchPublished
	}
	if status == string(BatchArchived) {
		return BatchArchived
	}
	return BatchDraft
}

func Summary(records []Record) (matches, overages, shortages int) {
	for _, record := range records {
		switch record.Result {
		case "match":
			matches++
		case "overage":
			overages++
		case "shortage":
			shortages++
		}
	}
	return
}
