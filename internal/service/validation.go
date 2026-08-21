package service

import (
	"context"
	"fmt"

	"inventoryseal/internal/domain"
)

type ValidationIssue struct {
	RecordID string
	Message  string
}

func (s *Service) ValidateBatch(ctx context.Context, batchID string) ([]ValidationIssue, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	if _, err := s.db.GetBatch(batchID); err != nil {
		return nil, err
	}
	records, err := s.db.ListRecords(batchID)
	if err != nil {
		return nil, err
	}
	issues := make([]ValidationIssue, 0)
	for _, record := range records {
		if err := domain.ValidateRecord(record); err != nil {
			issues = append(issues, ValidationIssue{record.ID, err.Error()})
			continue
		}
		if record.Result != domain.EvaluateRecord(record) {
			issues = append(issues, ValidationIssue{record.ID, "stored result is stale"})
		}
		if record.Confirmed && !record.Published {
			issues = append(issues, ValidationIssue{record.ID, fmt.Sprintf("%s confirmed without publication", record.ID)})
		}
	}
	return issues, nil
}

func (s *Service) CountByResult(ctx context.Context, batchID string) (map[string]int, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	records, err := s.db.ListRecords(batchID)
	if err != nil {
		return nil, err
	}
	result := map[string]int{"match": 0, "overage": 0, "shortage": 0}
	for _, record := range records {
		result[record.Result]++
	}
	return result, nil
}
