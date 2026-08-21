package service

import (
	"context"
	"fmt"
	"inventoryseal/internal/domain"
	"inventoryseal/internal/store"
)

func (s *Service) RetentionCheck(ctx context.Context, batchID string) (bool, error) {
	if err := checkContext(ctx); err != nil {
		return false, err
	}
	batch, err := s.db.GetBatch(batchID)
	if err != nil {
		return false, err
	}
	if !domain.IsFinal(batch.Status) {
		return false, fmt.Errorf("batch %s is not final", batchID)
	}
	snapshots, err := s.db.ListSnapshots(batchID)
	if err != nil {
		return false, err
	}
	return len(snapshots) > 0, nil
}

func (s *Service) DeleteDraft(ctx context.Context, batchID string) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	batch, err := s.db.GetBatch(batchID)
	if err != nil {
		return err
	}
	if batch.Status != domain.BatchDraft {
		return fmt.Errorf("only draft batches can delete")
	}
	records, err := s.db.ListRecords(batchID)
	if err != nil {
		return err
	}
	return s.db.WriteTransaction(func(tx store.Transaction) error {
		for _, record := range records {
			if err := tx.DeleteRecord(record.ID); err != nil {
				return err
			}
		}
		return nil
	})
}
