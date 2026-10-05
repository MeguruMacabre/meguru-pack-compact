package host

import (
	"github.com/MeguruMacabre/meguru-pack-compact/internal/manifest"
	"github.com/MeguruMacabre/meguru-pack-compact/internal/scanner"
)

func buildPackManifest(appRoot string, paths InstancePaths) (manifest.Manifest, error) {

	scannedFiles, err := scanner.Scan(paths.InstanceRoot)
	if err != nil {
		return manifest.Manifest{}, err
	}

	packManifest, err := manifest.Build(
		1,
		scannedFiles,
		paths.InstanceRoot,
		paths.GameRoot,
	)
	if err != nil {
		return manifest.Manifest{}, err
	}

	err = manifest.Save(appRoot, packManifest)
	if err != nil {
		return manifest.Manifest{}, err
	}
	return packManifest, nil
}
