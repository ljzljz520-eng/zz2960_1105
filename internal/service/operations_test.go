package service

import (
	"context"
	"testing"
)

func TestPrepareReview(t *testing.T) {
	svc := newService(t)
	seedBatch(t, svc, "review-summary")
	summary, err := svc.PrepareReview(context.Background(), "review-summary")
	if err != nil || summary.Totals.Total != 1 || !summary.Policy.Allowed {
		t.Fatalf("summary %#v err %v", summary, err)
	}
}
