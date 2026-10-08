package host

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MeguruMacabre/meguru-pack-compact/internal/manifest"
	"github.com/MeguruMacabre/meguru-pack-compact/internal/storage"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func publishPackFiles(
	ctx context.Context,
	client *s3.Client,
	bucket string,
	instanceRoot string,
	packManifest manifest.Manifest,
) error {
	for _, file := range packManifest.Files {
		err := uploadPackFile(ctx, client, bucket, instanceRoot, file)
		if err != nil {
			return fmt.Errorf("publish pack file %q: %w", file.Path, err)
		}
	}
	return nil
}

func uploadPackFile(
	ctx context.Context,
	client *s3.Client,
	bucket string,
	instanceRoot string,
	file manifest.File,
) error {
	relativePath := filepath.FromSlash(file.Path)

	absoluteFilePath := filepath.Join(instanceRoot, relativePath)

	openedFile, err := os.Open(absoluteFilePath)
	if err != nil {
		return fmt.Errorf("open pack file %q: %w", file.Path, err)
	}
	defer openedFile.Close()

	key := storage.FileKey(file.Path)

	err = storage.Upload(
		ctx,
		client,
		bucket,
		key,
		openedFile,
	)
	if err != nil {
		return fmt.Errorf("upload pack file %q: %w", file.Path, err)
	}
	return nil
}

func publishManifest(
	ctx context.Context,
	client *s3.Client,
	bucket string,
	appRoot string,
) error {
	path := filepath.Join(appRoot, "manifest.json")

	manifestFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open manifest file %q: %w", path, err)
	}
	defer manifestFile.Close()

	key := storage.ManifestKey

	err = storage.Upload(ctx, client, bucket, key, manifestFile)
	if err != nil {
		return fmt.Errorf("upload manifest file %q: %w", manifestFile.Name(), err)
	}
	return nil
}
