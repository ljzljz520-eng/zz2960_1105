package telemetry

import "testing"

func TestSummary(t *testing.T) {
	events := []Event{New("review", "b", "a", "", 1), New("publish", "b", "a", "", 2)}
	result := Summarize(events)
	if result.Count != 2 || result.TotalValue != 3 || CountActor(events, "a") != 2 {
		t.Fatalf("summary %#v", result)
	}
}
