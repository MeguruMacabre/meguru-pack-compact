package client

import (
	"fmt"
	"path/filepath"
	"strings"
)

// validateManifestPath проверяет безопасность относительного пути и преобразует его в формат текущей ОС.
func validateManifestPath(manifestPath string) (string, error) {
	if manifestPath == "" || manifestPath == "." {
		return "", fmt.Errorf(
			"invalid manifest path %q: empty or current directory",
			manifestPath,
		)
	}

	if strings.ContainsRune(manifestPath, '\\') {
		return "", fmt.Errorf(
			"invalid manifest path %q: backslashes are not allowed",
			manifestPath,
		)
	}

	localPath, err := filepath.Localize(manifestPath)
	if err != nil {
		return "", fmt.Errorf(
			"invalid manifest path %q: %w",
			manifestPath,
			err,
		)
	}

	return localPath, nil
}
