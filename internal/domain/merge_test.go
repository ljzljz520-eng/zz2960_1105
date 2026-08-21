package domain

import "testing"

func TestMergeRecords(t *testing.T) {
	got := MergeRecords([]Record{{ID: "r", Version: 1, Observed: 1}}, []Record{{ID: "r", Version: 2, Observed: 2}})
	if len(got) != 1 || got[0].Observed != 2 {
		t.Fatalf("records %#v", got)
	}
}
