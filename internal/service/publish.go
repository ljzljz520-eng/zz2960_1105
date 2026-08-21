package service

import (
	"context"
	"encoding/json"
	"fmt"

	"inventoryseal/internal/crypto"
	"inventoryseal/internal/domain"
	"inventoryseal/internal/store"
)

func (s *Service) ConfirmAndPublish(ctx context.Context, batchID, actor string) (domain.ExportSnapshot, error) {
	if err := checkContext(ctx); err != nil {
		return domain.ExportSnapshot{}, err
	}
	batch, err := s.db.GetBatch(batchID)
	if err != nil {
		return domain.ExportSnapshot{}, err
	}
	if batch.Status != domain.BatchReview {
		return domain.ExportSnapshot{}, fmt.Errorf("batch %s is not awaiting confirmation", batchID)
	}
	records, err := s.db.ListRecords(batchID)
	if err != nil {
		return domain.ExportSnapshot{}, err
	}
	if len(records) == 0 {
		return domain.ExportSnapshot{}, fmt.Errorf("batch has no records")
	}
	for i := range records {
		if err := checkContext(ctx); err != nil {
			return domain.ExportSnapshot{}, err
		}
		records[i].Confirmed = true
		records[i].Published = true
		records[i].Result = domain.EvaluateRecord(records[i])
	}
	published := batch
	published.Status = domain.BatchConfirmed
	published.Version++
	published.Status = domain.BatchPublished
	published.Version++
	payload, err := json.Marshal(struct {
		Batch   domain.Batch    `json:"batch"`
		Records []domain.Record `json:"records"`
	}{published, domain.SortRecords(records)})
	if err != nil {
		return domain.ExportSnapshot{}, err
	}
	if err := checkContext(ctx); err != nil {
		return domain.ExportSnapshot{}, err
	}
	sequence, err := s.db.NextSnapshotSequence(batchID)
	if err != nil {
		return domain.ExportSnapshot{}, err
	}
	snapshot := domain.ExportSnapshot{
		ID:          fmt.Sprintf("%s-export-%03d", batchID, sequence),
		BatchID:     batchID,
		Payload:     string(payload),
		Digest:      crypto.Digest(payload),
		RecordCount: len(records),
		PublishedBy: actor,
		Sequence:    sequence,
	}
	if err := checkContext(ctx); err != nil {
		return domain.ExportSnapshot{}, err
	}
	auditSequence, err := s.db.NextAuditSequence(batchID)
	if err != nil {
		return domain.ExportSnapshot{}, err
	}
	audit := domain.AuditEvent{
		ID:       fmt.Sprintf("%s-audit-%03d", batchID, auditSequence),
		BatchID:  batchID,
		Action:   "publish",
		Actor:    actor,
		Detail:   snapshot.Digest,
		Sequence: auditSequence,
	}
	if err := checkContext(ctx); err != nil {
		return domain.ExportSnapshot{}, err
	}
	err = s.db.WriteTransaction(func(tx store.Transaction) error {
		if err := checkContext(ctx); err != nil {
			return err
		}
		for _, record := range records {
			if err := tx.PutRecord(record); err != nil {
				return err
			}
		}
		if err := tx.PutBatch(published); err != nil {
			return err
		}
		if err := tx.PutSnapshot(snapshot); err != nil {
			return err
		}
		return tx.PutAudit(audit)
	})
	if err != nil {
		return domain.ExportSnapshot{}, err
	}
	return snapshot, nil
}

func (s *Service) ArchiveBatch(ctx context.Context, batchID, actor string) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	return s.transition(batchID, domain.BatchArchived, actor, "archive completed")
}

func (s *Service) LatestExport(batchID string) (domain.ExportSnapshot, error) {
	items, err := s.db.ListSnapshots(batchID)
	if err != nil {
		return domain.ExportSnapshot{}, err
	}
	if len(items) == 0 {
		return domain.ExportSnapshot{}, fmt.Errorf("no export for %s", batchID)
	}
	return items[len(items)-1], nil
}
