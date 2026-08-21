package service

import (
	"context"
	"fmt"
	"strings"

	"inventoryseal/internal/domain"
	"inventoryseal/internal/ledger"
	"inventoryseal/internal/policy"
)

type ReviewSummary struct {
	Batch  domain.Batch
	Totals ledger.Totals
	Policy policy.BatchDecision
	Issues []ValidationIssue
}

func (s *Service) PrepareReview(ctx context.Context, batchID string) (ReviewSummary, error) {
	if err := checkContext(ctx); err != nil {
		return ReviewSummary{}, err
	}
	batch, err := s.db.GetBatch(batchID)
	if err != nil {
		return ReviewSummary{}, err
	}
	records, err := s.db.ListRecords(batchID)
	if err != nil {
		return ReviewSummary{}, err
	}
	issues, err := s.ValidateBatch(ctx, batchID)
	if err != nil {
		return ReviewSummary{}, err
	}
	return ReviewSummary{Batch: batch, Totals: ledger.Count(records), Policy: policy.CanReview(batch, records), Issues: issues}, nil
}

func (s *Service) ResolveIssue(ctx context.Context, batchID, recordID, actor string) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	record, err := s.db.GetRecord(recordID)
	if err != nil {
		return err
	}
	if record.BatchID != batchID {
		return fmt.Errorf("record belongs to another batch")
	}
	record.Result = domain.EvaluateRecord(record)
	record.UpdatedBy = actor
	if err := s.db.SaveRecord(record); err != nil {
		return err
	}
	return s.recordAudit(batchID, "issue.resolve", actor, recordID)
}

func (s *Service) ConfirmRecord(ctx context.Context, batchID, recordID, actor string) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	record, err := s.db.GetRecord(recordID)
	if err != nil {
		return err
	}
	if record.BatchID != batchID {
		return fmt.Errorf("record belongs to another batch")
	}
	if record.Result != domain.EvaluateRecord(record) {
		return fmt.Errorf("record result is stale")
	}
	record.Confirmed = true
	record.UpdatedBy = actor
	return s.db.SaveRecord(record)
}

func (s *Service) PublishStatus(ctx context.Context, batchID string) (string, error) {
	if err := checkContext(ctx); err != nil {
		return "", err
	}
	batch, err := s.db.GetBatch(batchID)
	if err != nil {
		return "", err
	}
	records, err := s.db.ListRecords(batchID)
	if err != nil {
		return "", err
	}
	decision := policy.CanPublish(batch, records)
	if !decision.Allowed {
		return "", fmt.Errorf("publication denied: %s", decision.Reason)
	}
	return strings.ToLower(string(batch.Status)), nil
}

func (s *Service) ArchiveSummary(ctx context.Context, batchID string) (ledger.Totals, error) {
	if err := checkContext(ctx); err != nil {
		return ledger.Totals{}, err
	}
	records, err := s.db.ListRecords(batchID)
	if err != nil {
		return ledger.Totals{}, err
	}
	return ledger.Count(records), nil
}

func (s *Service) ActorSummary(ctx context.Context, batchID string) ([]ledger.ActorTotals, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	events, err := s.db.ListAudits(batchID)
	if err != nil {
		return nil, err
	}
	return ledger.Actors(events), nil
}

func (s *Service) Timeline(ctx context.Context, batchID string) ([]domain.AuditEvent, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	events, err := s.db.ListAudits(batchID)
	if err != nil {
		return nil, err
	}
	return ledger.Timeline(events), nil
}
