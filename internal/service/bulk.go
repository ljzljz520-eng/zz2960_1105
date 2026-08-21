package service

import (
	"context"
	"fmt"

	"inventoryseal/internal/domain"
)

type BulkResult struct {
	Accepted int
	Rejected int
	Errors   []string
}

func (s *Service) AddRecords(ctx context.Context, records []domain.Record) BulkResult {
	result := BulkResult{Errors: make([]string, 0)}
	for _, record := range records {
		if err := s.AddRecord(ctx, record); err != nil {
			result.Rejected++
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", record.ID, err))
			continue
		}
		result.Accepted++
	}
	return result
}

func (s *Service) ConfirmRecords(ctx context.Context, batchID, actor string, ids []string) BulkResult {
	result := BulkResult{Errors: make([]string, 0)}
	for _, id := range ids {
		if err := s.ConfirmRecord(ctx, batchID, id, actor); err != nil {
			result.Rejected++
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", id, err))
			continue
		}
		result.Accepted++
	}
	return result
}

func (s *Service) AddRecordsFromMap(ctx context.Context, batchID, actor string, values map[string][2]int) BulkResult {
	items := make([]domain.Record, 0, len(values))
	for id, counts := range values {
		items = append(items, domain.Record{ID: id, BatchID: batchID, Label: id, Expected: counts[0], Observed: counts[1], UpdatedBy: actor})
	}
	return s.AddRecords(ctx, items)
}

func (s *Service) EnsureBatch(ctx context.Context, id, title, owner, actor string) (domain.Batch, error) {
	if batch, err := s.db.GetBatch(id); err == nil {
		return batch, nil
	}
	batch := domain.Batch{ID: id, Title: title, Owner: owner, CreatedBy: actor}
	if err := s.CreateBatch(ctx, batch); err != nil {
		return domain.Batch{}, err
	}
	return s.db.GetBatch(id)
}
