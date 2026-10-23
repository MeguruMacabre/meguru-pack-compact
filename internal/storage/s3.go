package storage

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/MeguruMacabre/meguru-pack-compact/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
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

func OpenObject(
	ctx context.Context,
	client *s3.Client,
	bucket string,
	key string,
) (io.ReadCloser, error) {

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	output, err := client.GetObject(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("get S3 object %q: %w", key, err)
	}
	return output.Body, nil
}

func ObjectExists(
	ctx context.Context,
	client *s3.Client,
	bucket string,
	key string,
) (bool, error) {
	_, err := client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		var apiErr smithy.APIError

		if errors.As(err, &apiErr) {
			switch apiErr.ErrorCode() {
			case "NotFound", "404", "NoSuchKey":
				return false, nil
			}
		}

		return false, fmt.Errorf("head S3 object %q: %w", key, err)
	}

	return true, nil
}

func UploadIfAbsent(
	ctx context.Context,
	client *s3.Client,
	bucket string,
	key string,
	body io.Reader,
) error {
	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		IfNoneMatch: aws.String("*"),
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        body,
	})
	if err != nil {
		return err
	}

	return nil
}
