package crypto

import "testing"

func TestDigestDeterministic(t *testing.T) {
	first := Digest([]byte("batch"))
	if first != Digest([]byte("batch")) || !Verify([]byte("batch"), first) {
		t.Fatal("digest is not deterministic")
	}
}
