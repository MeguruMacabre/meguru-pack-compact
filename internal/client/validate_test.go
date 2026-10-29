package client

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/MeguruMacabre/meguru-pack-compact/internal/manifest"
)

func TestValidateManifest(t *testing.T) {
	sum := sha256.Sum256(nil)
	validHash := hex.EncodeToString(sum[:])

	tests := []struct {
		name    string
		modify  func(*manifest.Manifest)
		wantErr string
	}{
		{
			name: "valid manifest",
		},
		{
			name: "unsupported version",
			modify: func(m *manifest.Manifest) {
				m.FormatVersion = 2
			},
			wantErr: "unsupported manifest version",
		},
		{
			name: "invalid path",
			modify: func(m *manifest.Manifest) {
				m.Files[0].Path = "../secret.txt"
			},
			wantErr: "invalid file path",
		},
		{
			name: "negative size",
			modify: func(m *manifest.Manifest) {
				m.Files[0].Size = -1
			},
			wantErr: "negative file size",
		},
		{
			name: "invalid SHA-256 length",
			modify: func(m *manifest.Manifest) {
				m.Files[0].SHA256 = "abcd"
			},
			wantErr: "invalid SHA-256",
		},
		{
			name: "invalid SHA-256 characters",
			modify: func(m *manifest.Manifest) {
				m.Files[0].SHA256 = strings.Repeat("z", 64)
			},
			wantErr: "invalid SHA-256",
		},
		{
			name: "uppercase SHA-256",
			modify: func(m *manifest.Manifest) {
				m.Files[0].SHA256 = strings.ToUpper(validHash)
			},
			wantErr: "non-canonical SHA-256",
		},
		{
			name: "duplicate path",
			modify: func(m *manifest.Manifest) {
				m.Files = append(m.Files, m.Files[0])
			},
			wantErr: "duplicate file path",
		},
		{
			name: "same hash different paths",
			modify: func(m *manifest.Manifest) {
				m.Files = append(m.Files, manifest.File{
					Path:   ".minecraft/mods/another.jar",
					Size:   0,
					SHA256: validHash,
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packManifest := manifest.Manifest{
				FormatVersion: 1,
				Files: []manifest.File{
					{
						Path:    ".minecraft/mods/test.jar",
						Size:    0,
						SHA256:  validHash,
						Managed: true,
					},
				},
			}

			if tt.modify != nil {
				tt.modify(&packManifest)
			}

			err := validateManifest(packManifest)

			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("expected error containing %q", tt.wantErr)
			}

			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf(
					"expected error containing %q, got %v",
					tt.wantErr,
					err,
				)
			}
		})
	}
}
