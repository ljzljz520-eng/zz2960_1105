package httpapi

import (
	"net/http/httptest"
	"path/filepath"
	"testing"

	"inventoryseal/internal/service"
	"inventoryseal/internal/store"
)

func TestHealth(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "d.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	response := httptest.NewRecorder()
	New(service.New(db)).Handler().ServeHTTP(response, httptest.NewRequest("GET", "/healthz", nil))
	if response.Code != 200 {
		t.Fatalf("status %d", response.Code)
	}
}
