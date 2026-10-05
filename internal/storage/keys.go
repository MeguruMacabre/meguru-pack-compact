package storage

import "path"

const ManifestKey = "pack/manifest.json"

func FileKey(manifestPath string) string {
	return path.Join("pack", "files", manifestPath)
}
