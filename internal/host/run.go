package host

import (
	"context"
	"fmt"
	"os"
	"slices"

	"github.com/MeguruMacabre/meguru-pack-compact/internal/config"
	"github.com/MeguruMacabre/meguru-pack-compact/internal/diff"
	"github.com/MeguruMacabre/meguru-pack-compact/internal/manifest"
	"github.com/MeguruMacabre/meguru-pack-compact/internal/storage"
)

// Run запускает Host: сканирует сборку, сравнивает изменения и публикует обновление в S3.
func Run(ctx context.Context) error {
	appRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get application directory: %w", err)
	}

	cfg, err := loadOrCreateHostConfig(appRoot)
	if err != nil {
		return fmt.Errorf("load host config: %w", err)
	}

	if err := config.ValidateS3Config(cfg); err != nil {
		return fmt.Errorf("validate S3 config: %w", err)
	}

	paths, err := prepareInstance(appRoot, cfg)
	if err != nil {
		return fmt.Errorf("prepare instance: %w", err)
	}

	fmt.Println("Host Config:", cfg.InstanceDirectory)
	fmt.Println("Game directory:", paths.GameRoot)

	packManifest, err := buildPackManifest(appRoot, paths)
	if err != nil {
		return fmt.Errorf("build pack manifest: %w", err)
	}

	printManifest(packManifest)

	client, err := connectS3(ctx, cfg)
	if err != nil {
		return fmt.Errorf("connect to S3: %w", err)
	}

	exists, err := storage.ObjectExists(
		ctx,
		client,
		cfg.S3Bucket,
		storage.ManifestKey,
	)
	if err != nil {
		return fmt.Errorf("check published manifest: %w", err)
	}

	remoteManifest := manifest.Manifest{
		FormatVersion: packManifest.FormatVersion,
	}

	if exists {
		remoteManifest, err = readRemoteManifest(
			ctx,
			client,
			cfg.S3Bucket,
		)
		if err != nil {
			return fmt.Errorf("read published manifest: %w", err)
		}

		fmt.Println("Existing pack found in S3")
	} else {
		fmt.Println("No published pack found. Starting initial publication")
	}

	changes := diff.Compare(remoteManifest, packManifest)

	for _, change := range changes {
		fmt.Printf("%-8s %s\n", change.Action, change.Path)
	}

	err = publishPackFiles(
		ctx,
		client,
		cfg.S3Bucket,
		paths.InstanceRoot,
		packManifest,
	)
	if err != nil {
		return fmt.Errorf("publish pack objects: %w", err)
	}

	fmt.Println("Pack objects ready")

	if exists &&
		remoteManifest.FormatVersion == packManifest.FormatVersion &&
		slices.Equal(remoteManifest.Files, packManifest.Files) {
		fmt.Println("Pack is already up to date")
		return nil
	}

	err = publishManifest(
		ctx,
		client,
		cfg.S3Bucket,
		packManifest,
	)
	if err != nil {
		return fmt.Errorf("publish pack manifest: %w", err)
	}

	fmt.Println("Pack published successfully")

	return nil
}
