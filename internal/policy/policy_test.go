package policy

import (
	"inventoryseal/internal/domain"
	"testing"
)

func TestDefaultRules(t *testing.T) {
	result := Evaluate(domain.Record{ID: "r", Label: "a", Expected: 2, Observed: 2, Result: "match", Version: 1}, DefaultRules())
	if !result.Passed {
		t.Fatal(result.Violations)
	}
}
