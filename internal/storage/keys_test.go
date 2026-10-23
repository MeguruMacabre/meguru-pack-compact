package storage

import "testing"

func TestObjectKey(t *testing.T) {
	hash := "abc123"
	want := "pack/objects/abc123"
	got := ObjectKey(hash)
	if got != want {
		t.Errorf("ObjectKey(%q) = %q, want %q", hash, got, want)
	}
}
