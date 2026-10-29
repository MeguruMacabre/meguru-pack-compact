package client

import (
	"testing"

	"github.com/MeguruMacabre/meguru-pack-compact/internal/manifest"
)

func TestValidatePathConflicts(t *testing.T) {
	tests := []struct {
		name    string
		paths   []string
		wantErr bool
	}{
		{
			name:  "different files",
			paths: []string{"mods/a.jar", "mods/b.jar"},
		},
		{
			name:  "different directories",
			paths: []string{"mods/a/file.jar", "mods/b/file.jar"},
		},
		{
			name:    "case conflict",
			paths:   []string{"mods/Test.jar", "mods/test.jar"},
			wantErr: true,
		},
		{
			name:    "duplicate path",
			paths:   []string{"mods/a.jar", "mods/a.jar"},
			wantErr: true,
		},
		{
			name:    "parent file conflict",
			paths:   []string{"mods/a", "mods/a/b.jar"},
			wantErr: true,
		},
		{
			name:    "nested parent conflict",
			paths:   []string{"mods/a", "mods/a/b/c.jar"},
			wantErr: true,
		},
		{
			name:    "reverse order conflict",
			paths:   []string{"mods/a/b.jar", "mods/a"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files := make([]manifest.File, 0, len(tt.paths))

			for _, filePath := range tt.paths {
				files = append(files, manifest.File{
					Path: filePath,
				})
			}

			err := validatePathConflicts(files)

			if tt.wantErr && err == nil {
				t.Fatal("expected conflict error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected conflict error: %v", err)
			}
		})
	}
}
