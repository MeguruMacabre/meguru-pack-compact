package host

import (
	"context"
	"encoding/json/v2"
	"fmt"
	"io"

	"github.com/MeguruMacabre/meguru-pack-compact/internal/manifest"
	"github.com/MeguruMacabre/meguru-pack-compact/internal/storage"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// readRemoteManifest скачивает опубликованный манифест из S3 и преобразует JSON в структуру Go.
func readRemoteManifest(
	ctx context.Context,
	client *s3.Client,
	bucket string,
) (manifest.Manifest, error) {
	body, err := storage.OpenObject(
		ctx,
		client,
		bucket,
		storage.ManifestKey,
	)
	if err != nil {
		return manifest.Manifest{}, err
	}

	defer body.Close()

	data, err := io.ReadAll(body)
	if err != nil {
		return manifest.Manifest{}, err
	}

	var result manifest.Manifest

	err = json.Unmarshal(data, &result)
	if err != nil {
		return manifest.Manifest{}, fmt.Errorf("read remote manifest: %w", err)
	}
	return result, nil
}
