package policy

import "inventoryseal/internal/domain"

type BatchDecision struct {
	Allowed bool
	Reason  string
}

func CanReview(batch domain.Batch, records []domain.Record) BatchDecision {
	if batch.Status != domain.BatchDraft {
		return BatchDecision{Reason: "batch must be draft"}
	}
	if len(records) == 0 {
		return BatchDecision{Reason: "at least one record required"}
	}
	for _, record := range records {
		if !Evaluate(record, DefaultRules()).Passed {
			return BatchDecision{Reason: "record policy failed"}
		}
	}
	return BatchDecision{Allowed: true, Reason: "ready for review"}
}

func CanPublish(batch domain.Batch, records []domain.Record) BatchDecision {
	if batch.Status != domain.BatchReview {
		return BatchDecision{Reason: "batch must be in review"}
	}
	for _, record := range records {
		if record.Confirmed || record.Published {
			continue
		}
		return BatchDecision{Reason: "all records require confirmation"}
	}
	return BatchDecision{Allowed: len(records) > 0, Reason: "ready for publication"}
}

func CanArchive(batch domain.Batch) BatchDecision {
	if batch.Status != domain.BatchPublished {
		return BatchDecision{Reason: "batch must be published"}
	}
	return BatchDecision{Allowed: true, Reason: "retention ready"}
}
