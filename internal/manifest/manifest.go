package manifest

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"os"
	"path/filepath"

	"github.com/MeguruMacabre/meguru-pack-compact/internal/scanner"
	"github.com/MeguruMacabre/meguru-pack-compact/internal/syncpolicy"
)

type Manifest struct {
	FormatVersion int    `json:"format_version"`
	Files         []File `json:"files"`
}

type File struct {
	Path    string `json:"path"`
	Size    int64  `json:"size"`
	SHA256  string `json:"sha256"`
	Managed bool   `json:"managed"`
}

func Build(formatVersion int, scannedFiles []scanner.File, instanceRoot string, gameRoot string) (Manifest, error) {
	files := make([]File, 0, len(scannedFiles))
	for _, file := range scannedFiles {
		absoluteFilePath := filepath.Join(instanceRoot, file.Path)
		isManaged, err := syncpolicy.IsManaged(gameRoot, absoluteFilePath)
		if err != nil {
			return Manifest{}, err
		}
		files = append(files, File{
			Path:    filepath.ToSlash(file.Path),
			Size:    file.Size,
			SHA256:  file.Hash,
			Managed: isManaged,
		})
	}

	manifest := Manifest{
		FormatVersion: formatVersion,
		Files:         files,
	}
	return manifest, nil
}

func Save(dir string, manifest Manifest) error {
	data, err := json.Marshal(&manifest, jsontext.WithIndent("  "))
	if err != nil {
		return err
	}

	manifestPath := filepath.Join(dir, "manifest.json")

	err = os.WriteFile(manifestPath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func Load(dir string) (Manifest, error) {
	manifestPath := filepath.Join(dir, "manifest.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return Manifest{}, err
	}

	var manifest Manifest
	err = json.Unmarshal(data, &manifest)
	if err != nil {
		return Manifest{}, err
	}
	return manifest, nil
}
