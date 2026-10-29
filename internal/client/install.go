package client

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// installFile устанавливает скачанный файл внутри указанного корня через временную копию.
// При успешной установке заменяет целевой файл, не удаляя исходный скачанный файл.
func installFile(
	root *os.Root,
	downloadedPath string,
	manifestPath string,
) error {
	destinationPath, err := validateManifestPath(manifestPath)
	if err != nil {
		return fmt.Errorf("validate installation path: %w", err)
	}

	source, err := os.Open(downloadedPath)
	if err != nil {
		return fmt.Errorf("open downloaded file: %w", err)
	}
	defer source.Close()

	parentDir := filepath.Dir(destinationPath)

	if parentDir != "." {
		if err := root.MkdirAll(parentDir, 0o755); err != nil {
			return fmt.Errorf("create installation directories: %w", err)
		}
	}

	tempPath := filepath.Join(
		parentDir,
		".meguru-"+rand.Text()+".part",
	)

	tempFile, err := root.OpenFile(
		tempPath,
		os.O_WRONLY|os.O_CREATE|os.O_EXCL,
		0o600,
	)
	if err != nil {
		return fmt.Errorf("create installation temp file: %w", err)
	}

	installed := false

	defer func() {
		_ = tempFile.Close()

		if !installed {
			_ = root.Remove(tempPath)
		}
	}()

	if _, err := io.Copy(tempFile, source); err != nil {
		return fmt.Errorf("copy downloaded file: %w", err)
	}

	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("close installation temp file: %w", err)
	}

	if err := root.Rename(tempPath, destinationPath); err != nil {
		return fmt.Errorf("install file %q: %w", manifestPath, err)
	}

	installed = true

	return nil
}
