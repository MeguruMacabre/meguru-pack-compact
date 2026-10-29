package diff

import "github.com/MeguruMacabre/meguru-pack-compact/internal/manifest"

type Action string

const (
	ADD    Action = "ADD"
	UPDATE Action = "UPDATE"
	DELETE Action = "DELETE"
	KEEP   Action = "KEEP"
)

type Change struct {
	Action  Action
	Path    string
	OldFile *manifest.File
	NewFile *manifest.File
}

// indexFiles создаёт индекс файлов по их путям для быстрого поиска.
func indexFiles(files []manifest.File) map[string]manifest.File {
	filesByPath := make(map[string]manifest.File)
	for _, f := range files {
		filesByPath[f.Path] = f
	}
	return filesByPath
}

// Compare сравнивает два манифеста и определяет действия ADD, UPDATE, DELETE и KEEP.
func Compare(oldManifest, newManifest manifest.Manifest) []Change {
	var changes []Change
	oldFilesByPath := indexFiles(oldManifest.Files)
	for _, newFile := range newManifest.Files {
		oldFile, exists := oldFilesByPath[newFile.Path]
		if !exists {
			changes = append(changes, Change{
				Action:  ADD,
				Path:    newFile.Path,
				OldFile: nil,
				NewFile: &newFile,
			})
			continue
		}
		if oldFile.SHA256 != newFile.SHA256 {
			changes = append(changes, Change{
				Action:  UPDATE,
				Path:    newFile.Path,
				OldFile: &oldFile,
				NewFile: &newFile,
			})
		} else {
			changes = append(changes, Change{
				Action:  KEEP,
				Path:    newFile.Path,
				OldFile: &oldFile,
				NewFile: &newFile,
			})
		}
	}

	newFilesByPath := indexFiles(newManifest.Files)
	for _, oldFile := range oldManifest.Files {
		_, exists := newFilesByPath[oldFile.Path]
		if !exists {
			changes = append(changes, Change{
				Action:  DELETE,
				Path:    oldFile.Path,
				OldFile: &oldFile,
				NewFile: nil,
			})
		}
	}
	return changes
}
