package diff

import (
	"testing"

	"github.com/MeguruMacabre/meguru-pack-compact/internal/manifest"
)

func TestCompareKeep(t *testing.T) {
	oldManifest := manifest.Manifest{
		Files: []manifest.File{
			{
				Path:   "keep.txt",
				SHA256: "AAA",
			},
		},
	}

	newManifest := manifest.Manifest{
		Files: []manifest.File{
			{
				Path:   "keep.txt",
				SHA256: "AAA",
			},
		},
	}

	changes := Compare(oldManifest, newManifest)

	if len(changes) != 1 {
		t.Fatalf("Expected 1 change, got %d", len(changes))
	}

	if changes[0].Action != KEEP {
		t.Fatalf("Expected 'KEEP', got %s", changes[0].Action)
	}
}

func TestCompareAdd(t *testing.T) {
	oldManifest := manifest.Manifest{}

	newManifest := manifest.Manifest{
		Files: []manifest.File{
			{
				Path:   "add.txt",
				SHA256: "BBB",
			},
		},
	}

	changes := Compare(oldManifest, newManifest)

	if len(changes) != 1 {
		t.Fatalf("Expected 1 change, got %d", len(changes))
	}

	if changes[0].Action != ADD {
		t.Errorf("expected ADD, got %s", changes[0].Action)
	}

	if changes[0].NewFile.Path != "add.txt" {
		t.Errorf("expected new path add.txt, got %s", changes[0].NewFile.Path)
	}

	if changes[0].OldFile != nil {
		t.Error("expected OldFile to be nil")
	}

	if changes[0].NewFile == nil {
		t.Error("expected NewFile to be non-nil")
	}
}

func TestCompareUpdate(t *testing.T) {
	oldManifest := manifest.Manifest{
		Files: []manifest.File{
			{
				Path:   "update.txt",
				SHA256: "AAA",
			},
		},
	}

	newManifest := manifest.Manifest{
		Files: []manifest.File{
			{
				Path:   "update.txt",
				SHA256: "BBB",
			},
		},
	}

	changes := Compare(oldManifest, newManifest)

	if len(changes) != 1 {
		t.Fatalf("Expected 1 change, got %d", len(changes))
	}

	if changes[0].Action != UPDATE {
		t.Fatalf("Expected 'UPDATE', got %s", changes[0].Action)
	}

	if changes[0].OldFile == nil {
		t.Fatal("expected OldFile to be non-nil")
	}

	if changes[0].NewFile == nil {
		t.Fatal("expected NewFile to be non-nil")
	}

	if changes[0].OldFile.SHA256 != "AAA" {
		t.Errorf(
			"expected old SHA256 AAA, got %s",
			changes[0].OldFile.SHA256,
		)
	}

	if changes[0].NewFile.SHA256 != "BBB" {
		t.Errorf(
			"expected new SHA256 BBB, got %s",
			changes[0].NewFile.SHA256,
		)
	}
}

func TestCompareDelete(t *testing.T) {
	oldManifest := manifest.Manifest{
		Files: []manifest.File{
			{
				Path:   "delete.txt",
				SHA256: "AAA",
			},
		},
	}

	newManifest := manifest.Manifest{}

	changes := Compare(oldManifest, newManifest)

	if len(changes) != 1 {
		t.Fatalf("Expected 1 change, got %d", len(changes))
	}

	if changes[0].Action != DELETE {
		t.Errorf("expected DELETE, got %s", changes[0].Action)
	}

	if changes[0].OldFile.Path != "delete.txt" {
		t.Errorf("expected old path delete.txt, got %s", changes[0].OldFile.Path)
	}

	if changes[0].OldFile == nil {
		t.Error("expected OldFile to be non-nil")
	}

	if changes[0].NewFile != nil {
		t.Error("expected NewFile to be nil")
	}
}
