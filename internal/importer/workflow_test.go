package importer

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"inventoryseal/internal/domain"
	"inventoryseal/internal/service"
	"inventoryseal/internal/store"
)

func TestWorkflowImportReport(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "data.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	svc := service.New(db)
	if err := svc.CreateBatch(context.Background(), domain.Batch{ID: "wf-import", Title: "incoming", Owner: "ops", CreatedBy: "loader"}); err != nil {
		t.Fatal(err)
	}
	result, err := Import(context.Background(), svc, "wf-import", "loader", strings.NewReader("r1,alpha,4,4\nr2,beta,6,5\n"))
	if err != nil || result.Imported != 2 {
		t.Fatalf("result %#v err %v", result, err)
	}
	issues, err := svc.ValidateBatch(context.Background(), "wf-import")
	if err != nil || len(issues) != 0 {
		t.Fatalf("issues %#v err %v", issues, err)
	}
	counts, err := svc.CountByResult(context.Background(), "wf-import")
	if err != nil || counts["match"] != 1 || counts["shortage"] != 1 {
		t.Fatalf("counts %#v err %v", counts, err)
	}
}
