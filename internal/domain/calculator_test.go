package domain

import "testing"

func TestMeasurements(t *testing.T) {
	value := MeasurementOf(Record{ID: "r", Expected: 4, Observed: 2})
	if value.Direction != "under" || value.Absolute != 2 {
		t.Fatalf("measurement %#v", value)
	}
}
