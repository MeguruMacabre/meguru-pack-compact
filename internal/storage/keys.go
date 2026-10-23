package storage

import "path"

const ManifestKey = "pack/manifest.json"

func FileKey(manifestPath string) string {
	return path.Join("pack", "files", manifestPath)
}

func ObjectKey(sha256 string) string {
	return path.Join("pack", "objects", sha256)
}
