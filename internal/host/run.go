package host

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/MeguruMacabre/meguru-pack-compact/internal/config"
	"github.com/MeguruMacabre/meguru-pack-compact/internal/storage"
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

	s3Credentials, err := config.LoadS3Credentials()
	if err != nil {
		return err
	}

	client, err := storage.NewS3Client(ctx, cfg, s3Credentials)
	if err != nil {
		return err
	}

	err = storage.CheckBucket(ctx, client, cfg.S3Bucket)
	if err != nil {
		return err
	}

	fmt.Println("S3 bucket connection successful")

	body := strings.NewReader("Meguru Pack Compact S3 test")

	err = storage.Upload(
		ctx,
		client,
		cfg.S3Bucket,
		"pack/test.txt",
		body,
	)
	if err != nil {
		return err
	}

	fmt.Println("Test object uploaded")

	return nil
}
