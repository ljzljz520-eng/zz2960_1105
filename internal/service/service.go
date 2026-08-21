package service

import (
	"context"
	"fmt"

	"inventoryseal/internal/domain"
	"inventoryseal/internal/store"
)

type Service struct{ db *store.DB }

func New(db *store.DB) *Service  { return &Service{db: db} }
func (s *Service) DB() *store.DB { return s.db }

func checkContext(ctx context.Context) error {
	if ctx == nil {
		return nil
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (s *Service) CreateBatch(ctx context.Context, batch domain.Batch) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	if batch.Status == "" {
		batch.Status = domain.BatchDraft
	}
	if batch.Version == 0 {
		batch.Version = 1
	}
	if err := domain.ValidateBatch(batch); err != nil {
		return err
	}
	if _, err := s.db.GetBatch(batch.ID); err == nil {
		return fmt.Errorf("batch %s already exists", batch.ID)
	}
	if err := s.db.SaveBatch(batch); err != nil {
		return err
	}
	return s.recordAudit(batch.ID, "create", batch.CreatedBy, "batch registered")
}

func (s *Service) AddRecord(ctx context.Context, record domain.Record) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	if err := domain.ValidateRecord(record); err != nil {
		return err
	}
	batch, err := s.db.GetBatch(record.BatchID)
	if err != nil {
		return err
	}
	if !domain.IsEditable(batch.Status) {
		return fmt.Errorf("batch %s is not editable", batch.ID)
	}
	if record.Version == 0 {
		record.Version = 1
	}
	record.Result = domain.EvaluateRecord(record)
	if err := s.db.SaveRecord(record); err != nil {
		return err
	}
	batch.RecordCount++
	if err := s.db.SaveBatch(batch); err != nil {
		return err
	}
	return s.recordAudit(record.BatchID, "record.add", record.UpdatedBy, record.ID)
}

func (s *Service) UpdateRecord(ctx context.Context, record domain.Record, actor string) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	batch, err := s.db.GetBatch(record.BatchID)
	if err != nil {
		return err
	}
	if !domain.IsEditable(batch.Status) {
		return fmt.Errorf("batch %s is not editable", batch.ID)
	}
	record.Result = domain.EvaluateRecord(record)
	record.UpdatedBy = actor
	if err := s.db.UpdateRecord(record, record.Version); err != nil {
		return err
	}
	return s.recordAudit(record.BatchID, "record.update", actor, record.ID)
}

func (s *Service) ReviewBatch(ctx context.Context, batchID, actor string) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	return s.transition(batchID, domain.BatchReview, actor, "review requested")
}

func (s *Service) ReopenBatch(ctx context.Context, batchID, actor string) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	batch, err := s.db.GetBatch(batchID)
	if err != nil {
		return err
	}
	if batch.Status != domain.BatchPublished {
		return fmt.Errorf("only published batches can reopen")
	}
	batch.Status = domain.BatchReview
	batch.Version++
	if err := s.db.SaveBatch(batch); err != nil {
		return err
	}
	return s.recordAudit(batchID, "reopen", actor, "correction cycle")
}

func (s *Service) transition(batchID string, next domain.BatchStatus, actor, detail string) error {
	batch, err := s.db.GetBatch(batchID)
	if err != nil {
		return err
	}
	if !domain.CanTransition(batch.Status, next) {
		return fmt.Errorf("invalid transition %s to %s", batch.Status, next)
	}
	batch.Status = next
	batch.Version++
	if err := s.db.SaveBatch(batch); err != nil {
		return err
	}
	return s.recordAudit(batchID, string(next), actor, detail)
}

func (s *Service) recordAudit(batchID, action, actor, detail string) error {
	audits, err := s.db.ListAudits(batchID)
	if err != nil {
		return err
	}
	event := domain.AuditEvent{ID: fmt.Sprintf("%s-audit-%03d", batchID, len(audits)+1), BatchID: batchID, Action: action, Actor: actor, Detail: detail, Sequence: len(audits) + 1}
	return s.db.SaveAudit(event)
}

func (s *Service) AddNote(ctx context.Context, note domain.CollaborationNote) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	if _, err := s.db.GetBatch(note.BatchID); err != nil {
		return err
	}
	notes, err := s.db.ListNotes(note.BatchID)
	if err != nil {
		return err
	}
	note.Sequence = len(notes) + 1
	if note.ID == "" {
		note.ID = fmt.Sprintf("%s-note-%03d", note.BatchID, note.Sequence)
	}
	if err := s.db.SaveNote(note); err != nil {
		return err
	}
	return s.recordAudit(note.BatchID, "note.add", note.Author, note.ID)
}
