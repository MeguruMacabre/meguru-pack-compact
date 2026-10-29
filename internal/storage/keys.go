package storage

import "path"

const ManifestKey = "pack/manifest.json"

// FileKey формирует ключ файла в прежней схеме хранения pack/files.
func FileKey(manifestPath string) string {
	return path.Join("pack", "files", manifestPath)
}

// ObjectKey формирует ключ объекта в S3 на основе его SHA-256.
func ObjectKey(sha256 string) string {
	return path.Join("pack", "objects", sha256)
}
