package instance

import (
	"os"
	"path/filepath"
)

func FindGameRoot(instancePath string) (string, error) {
	entries, err := os.ReadDir(instancePath)
	if err != nil {
		return "", err
	}
	for _, entry := range entries {
		if entry.Name() == ".minecraft" && entry.IsDir() {
			return filepath.Join(instancePath, entry.Name()), nil
		}
	}
	return instancePath, nil
}
