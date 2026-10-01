package scanner

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type File struct {
	Path string
	Size int64
	Hash string
}

var ignoredDirs = map[string]struct{}{
	"saves":                 {},
	"logs":                  {},
	"screenshots":           {},
	"crash-reports":         {},
	"server-resource-packs": {},
}

var ignoredFiles = map[string]struct{}{
	".DS_Store":   {},
	"Thumbs.db":   {},
	"desktop.ini": {},
}

func Scan(root string) ([]File, error) {
	var gameFiles []File
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if shouldSkipDir(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		if shouldSkipFile(d.Name()) {
			return nil
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		hash, err := hashFile(path)
		if err != nil {
			return err
		}

		file := File{
			Path: relPath,
			Size: info.Size(),
			Hash: hash,
		}

		gameFiles = append(gameFiles, file)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return gameFiles, nil
}

func shouldSkipDir(name string) bool {
	_, exists := ignoredDirs[name]
	return exists
}

func shouldSkipFile(name string) bool {
	_, exists := ignoredFiles[name]
	return exists
}

func hashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	_, err = io.Copy(hasher, file)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
