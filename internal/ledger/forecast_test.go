package ledger

import (
	"inventoryseal/internal/domain"
	"testing"
)

func TestForecast(t *testing.T) {
	value := ForecastBatch([]domain.Record{{Expected: 2, Observed: 3}})
	if value.Delta != 1 || value.Confidence != "review" {
		t.Fatalf("forecast %#v", value)
	}
}
