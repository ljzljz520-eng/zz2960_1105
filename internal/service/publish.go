package service

import (
	"context"
	"encoding/json"
	"fmt"

	"inventoryseal/internal/crypto"
	"inventoryseal/internal/domain"
)

func (s *Service) ConfirmAndPublish(ctx context.Context, batchID, actor string) (domain.ExportSnapshot, error) {
	workCtx := context.Background()
	if err := checkContext(workCtx); err != nil {
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
		if err := checkContext(workCtx); err != nil {
			return domain.ExportSnapshot{}, err
		}
		records[i].Confirmed = true
		records[i].Published = true
		records[i].Result = domain.EvaluateRecord(records[i])
		if err := s.db.SaveRecord(records[i]); err != nil {
			return domain.ExportSnapshot{}, err
		}
	}
	batch.Status = domain.BatchConfirmed
	batch.Version++
	if err := s.db.SaveBatch(batch); err != nil {
		return domain.ExportSnapshot{}, err
	}
	batch.Status = domain.BatchPublished
	batch.Version++
	if err := s.db.SaveBatch(batch); err != nil {
		return domain.ExportSnapshot{}, err
	}
	payload, err := json.Marshal(struct {
		Batch   domain.Batch    `json:"batch"`
		Records []domain.Record `json:"records"`
	}{batch, domain.SortRecords(records)})
	if err != nil {
		return domain.ExportSnapshot{}, err
	}
	snapshot := domain.ExportSnapshot{ID: fmt.Sprintf("%s-export-%03d", batchID, batch.Version), BatchID: batchID, Payload: string(payload), Digest: crypto.Digest(payload), RecordCount: len(records), PublishedBy: actor}
	previous, _ := s.db.ListSnapshots(batchID)
	snapshot.Sequence = len(previous) + 1
	if err := s.db.SaveSnapshot(snapshot); err != nil {
		return domain.ExportSnapshot{}, err
	}
	if err := s.recordAudit(batchID, "publish", actor, snapshot.Digest); err != nil {
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
