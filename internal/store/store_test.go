package store

import (
	"path/filepath"
	"testing"

	"inventoryseal/internal/domain"
)

func TestOpenSaveGet(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "data.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	b := domain.Batch{ID: "b1", Title: "first", Owner: "ops", Version: 1}
	if err := db.SaveBatch(b); err != nil {
		t.Fatal(err)
	}
	got, err := db.GetBatch("b1")
	if err != nil || got.Title != b.Title {
		t.Fatalf("got %#v err %v", got, err)
	}
}
