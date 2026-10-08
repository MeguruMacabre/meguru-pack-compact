package host

import (
	"context"
	"fmt"

	"github.com/MeguruMacabre/meguru-pack-compact/internal/config"
	"github.com/MeguruMacabre/meguru-pack-compact/internal/storage"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func connectS3(
	ctx context.Context,
	cfg config.HostConfig,
) (*s3.Client, error) {
	s3Credentials, err := config.LoadS3Credentials()
	if err != nil {
		return nil, fmt.Errorf("load S3 credentials: %w", err)
	}

	client, err := storage.NewS3Client(ctx, cfg, s3Credentials)
	if err != nil {
		return nil, fmt.Errorf("create S3 client: %w", err)
	}

	err = storage.CheckBucket(ctx, client, cfg.S3Bucket)
	if err != nil {
		return nil, fmt.Errorf("check S3 bucket: %w", err)
	}

	return client, nil
}
