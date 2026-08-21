package ledger

import "inventoryseal/internal/domain"

type Forecast struct {
	Expected   int
	Observed   int
	Delta      int
	Confidence string
}

func ForecastRecord(record domain.Record) Forecast {
	delta := record.Observed - record.Expected
	confidence := "stable"
	if delta != 0 {
		confidence = "review"
	}
	return Forecast{Expected: record.Expected, Observed: record.Observed, Delta: delta, Confidence: confidence}
}
func ForecastBatch(records []domain.Record) Forecast {
	result := Forecast{Confidence: "stable"}
	for _, record := range records {
		value := ForecastRecord(record)
		result.Expected += value.Expected
		result.Observed += value.Observed
	}
	result.Delta = result.Observed - result.Expected
	if result.Delta != 0 {
		result.Confidence = "review"
	}
	return result
}
func Balance(records []domain.Record) int {
	total := 0
	for _, record := range records {
		total += record.Observed
	}
	return total
}
func ExpectedBalance(records []domain.Record) int {
	total := 0
	for _, record := range records {
		total += record.Expected
	}
	return total
}
