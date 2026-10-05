package storage

import (
	"context"
	"io"

	"github.com/MeguruMacabre/meguru-pack-compact/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewS3Client(
	ctx context.Context,
	hostConfig config.HostConfig,
	creds config.S3Credentials,
) (*s3.Client, error) {
	credentialsProvider := credentials.NewStaticCredentialsProvider(
		creds.AccessKey,
		creds.SecretKey,
		"",
	)

	sdkConfig, err := awsconfig.LoadDefaultConfig(
		ctx,
		awsconfig.WithRegion(hostConfig.S3Region),
		awsconfig.WithCredentialsProvider(credentialsProvider),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(sdkConfig, func(options *s3.Options) {
		options.BaseEndpoint = aws.String(hostConfig.S3Endpoint)
	})

	return client, nil
}

func CheckBucket(
	ctx context.Context,
	client *s3.Client,
	bucket string,
) error {
	input := &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	}

	_, err := client.HeadBucket(ctx, input)
	if err != nil {
		return err
	}
	return nil
}

func Upload(
	ctx context.Context,
	client *s3.Client,
	bucket string,
	key string,
	body io.Reader,
) error {
	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   body,
	})
	if err != nil {
		return err
	}

	return nil
}
