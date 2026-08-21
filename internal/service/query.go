package service

import (
	"context"
	"strings"

	"inventoryseal/internal/domain"
	"inventoryseal/internal/reconcile"
)

func (s *Service) FindBatches(ctx context.Context, term string, status domain.BatchStatus) ([]domain.Batch, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	items, err := s.db.ListBatches()
	if err != nil {
		return nil, err
	}
	result := make([]domain.Batch, 0, len(items))
	needle := strings.ToLower(term)
	for _, item := range items {
		if status != "" && item.Status != status {
			continue
		}
		if needle != "" && !strings.Contains(strings.ToLower(item.ID+" "+item.Title+" "+item.Owner), needle) {
			continue
		}
		result = append(result, item)
	}
	return result, nil
}

func (s *Service) ReconcileBatch(ctx context.Context, batchID string) (reconcile.Report, error) {
	if err := checkContext(ctx); err != nil {
		return reconcile.Report{}, err
	}
	batch, err := s.db.GetBatch(batchID)
	if err != nil {
		return reconcile.Report{}, err
	}
	records, err := s.db.ListRecords(batchID)
	if err != nil {
		return reconcile.Report{}, err
	}
	return reconcile.Compact(reconcile.Analyze(batch, records)), nil
}

func (s *Service) BatchDetails(ctx context.Context, batchID string) (domain.Batch, []domain.Record, []domain.CollaborationNote, error) {
	if err := checkContext(ctx); err != nil {
		return domain.Batch{}, nil, nil, err
	}
	batch, err := s.db.GetBatch(batchID)
	if err != nil {
		return domain.Batch{}, nil, nil, err
	}
	records, err := s.db.ListRecords(batchID)
	if err != nil {
		return domain.Batch{}, nil, nil, err
	}
	notes, err := s.db.ListNotes(batchID)
	if err != nil {
		return domain.Batch{}, nil, nil, err
	}
	return batch, records, notes, nil
}

func (s *Service) AuditTrail(ctx context.Context, batchID string) ([]domain.AuditEvent, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	return s.db.ListAudits(batchID)
}
