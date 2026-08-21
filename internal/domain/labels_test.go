package domain

import "testing"

func TestNormalizeLabel(t *testing.T) {
	if NormalizeLabel(" alpha ") != "ALPHA" {
		t.Fatal("label not normalized")
	}
}
