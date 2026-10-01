package syncpolicy

import (
	"path/filepath"
	"strings"
)

var managedDirs = map[string]struct{}{
	"mods":          {},
	"kubejs":        {},
	"resourcepacks": {},
	"shaderpacks":   {},
}

func IsManaged(gameRoot string, filePath string) (bool, error) {
	relPath, err := filepath.Rel(gameRoot, filePath)
	if err != nil {
		return false, err
	}

	names := strings.Split(relPath, string(filepath.Separator))
	_, exist := managedDirs[names[0]]
	return exist, nil
}
