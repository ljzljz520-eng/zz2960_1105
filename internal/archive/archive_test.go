package archive

import "testing"

func TestCompressionRoundTrip(t *testing.T) {
	original := []byte("inventory payload")
	packed, err := Compress(original)
	if err != nil {
		t.Fatal(err)
	}
	unpacked, err := Decompress(packed)
	if err != nil || string(unpacked) != string(original) {
		t.Fatalf("unpacked %q err %v", unpacked, err)
	}
}
