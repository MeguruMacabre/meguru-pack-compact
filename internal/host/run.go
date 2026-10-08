package host

import (
	"context"
	"fmt"
	"os"

	"github.com/MeguruMacabre/meguru-pack-compact/internal/config"
)

func Run(ctx context.Context) error {
	appRoot, err := os.Getwd()
	if err != nil {
		return err
	}

	cfg, err := loadOrCreateHostConfig(appRoot)
	if err != nil {
		return err
	}

	err = config.ValidateS3Config(cfg)
	if err != nil {
		return err
	}

	paths, err := prepareInstance(appRoot, cfg)
	if err != nil {
		return err
	}

	fmt.Println("Host Config:", cfg.InstanceDirectory)
	fmt.Println("Game directory:", paths.GameRoot)

	packManifest, err := buildPackManifest(appRoot, paths)
	if err != nil {
		return err
	}

	printManifest(packManifest)

	client, err := connectS3(ctx, cfg)
	if err != nil {
		return err
	}

	err = publishPackFiles(
		ctx,
		client,
		cfg.S3Bucket,
		paths.InstanceRoot,
		packManifest,
	)
	if err != nil {
		return err
	}

	fmt.Println("Pack files uploaded")

	err = publishManifest(
		ctx,
		client,
		cfg.S3Bucket,
		appRoot,
	)
	if err != nil {
		return err
	}

	fmt.Println("Manifest published")

	return nil
}
