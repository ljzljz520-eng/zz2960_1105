package archive

import "testing"

func TestIndex(t *testing.T) {
	index := NewIndex()
	index.Add("b", "2")
	index.Add("a", "1")
	if len(index.Search("A")) != 1 || !index.Has("b") {
		t.Fatalf("index %#v", index)
	}
}
