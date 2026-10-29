package client

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path"
	"strings"

	"github.com/MeguruMacabre/meguru-pack-compact/internal/manifest"
)

// validateManifest проверяет версию манифеста, пути файлов, размеры, SHA-256 и дубликаты.
func validateManifest(packManifest manifest.Manifest) error {
	if packManifest.FormatVersion != 1 {
		return fmt.Errorf(
			"unsupported manifest version: %d",
			packManifest.FormatVersion,
		)
	}

	seen := make(map[string]bool, len(packManifest.Files))

	for _, file := range packManifest.Files {
		if _, err := validateManifestPath(file.Path); err != nil {
			return fmt.Errorf("invalid file path %q: %w", file.Path, err)
		}

		if file.Size < 0 {
			return fmt.Errorf(
				"negative file size for %q: %d",
				file.Path,
				file.Size,
			)
		}

		decoded, err := hex.DecodeString(file.SHA256)
		if err != nil || len(decoded) != sha256.Size {
			return fmt.Errorf("invalid SHA-256 for %q", file.Path)
		}

		if hex.EncodeToString(decoded) != file.SHA256 {
			return fmt.Errorf("non-canonical SHA-256 for %q", file.Path)
		}

		if seen[file.Path] {
			return fmt.Errorf("duplicate file path: %q", file.Path)
		}

		seen[file.Path] = true
	}

	if err := validatePathConflicts(packManifest.Files); err != nil {
		return fmt.Errorf("validate path conflicts: %w", err)
	}

	return nil
}

// validatePathConflicts проверяет конфлиты имён файлов и родительских путей.
func validatePathConflicts(files []manifest.File) error {
	paths := make(map[string]string, len(files))

	for _, file := range files {
		normalized := strings.ToLower(file.Path)

		if existing, ok := paths[normalized]; ok {
			return fmt.Errorf(
				"conflicting file paths %q and %q",
				existing,
				file.Path,
			)
		}

		paths[normalized] = file.Path
	}

	for _, file := range files {
		parent := path.Dir(strings.ToLower(file.Path))

		for parent != "." {
			if existing, ok := paths[parent]; ok {
				return fmt.Errorf(
					"file %q conflicts with parent file %q",
					file.Path,
					existing,
				)
			}

			parent = path.Dir(parent)
		}
	}
	return nil
}
