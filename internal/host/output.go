package host

import (
	"fmt"

	"github.com/MeguruMacabre/meguru-pack-compact/internal/manifest"
)

func printManifest(packManifest manifest.Manifest) {
	fmt.Printf("Manifest v%d\n", packManifest.FormatVersion)
	fmt.Printf("Files: %d\n\n", len(packManifest.Files))

	for _, file := range packManifest.Files {
		status := "Install-only"
		if file.Managed {
			status = "Managed"
		}

		fmt.Printf(
			"%-45s | %10d bytes | %-12s | %s\n",
			file.Path,
			file.Size,
			status,
			file.SHA256,
		)
	}
}
