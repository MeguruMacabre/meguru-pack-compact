package host

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/MeguruMacabre/meguru-pack-compact/internal/manifest"
	"github.com/MeguruMacabre/meguru-pack-compact/internal/storage"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

// publishPackFiles проверяет наличие всех объектов сборки в S3 и загружает недостающие.
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

// uploadPackFile создаёт проверенную временную копию файла и загружает её в S3, если объект отсутствует.
func uploadPackFile(
	ctx context.Context,
	client *s3.Client,
	bucket string,
	instanceRoot string,
	file manifest.File,
) error {
	key := storage.ObjectKey(file.SHA256)

	exists, err := storage.ObjectExists(ctx, client, bucket, key)
	if err != nil {
		return fmt.Errorf("check pack object %q: %w", key, err)
	}

	if exists {
		return nil
	}

	relativePath := filepath.FromSlash(file.Path)
	absoluteFilePath := filepath.Join(instanceRoot, relativePath)

	openedFile, err := os.Open(absoluteFilePath)
	if err != nil {
		return fmt.Errorf("open pack file %q: %w", file.Path, err)
	}
	defer openedFile.Close()

	snapshot, err := os.CreateTemp("", "meguru-pack-*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(snapshot.Name())
	defer snapshot.Close()

	_, err = io.Copy(snapshot, openedFile)
	if err != nil {
		return err
	}
	_, err = snapshot.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	err = verifyFileHash(snapshot, file.SHA256)
	if err != nil {
		return fmt.Errorf("verify pack file %q: %w", file.Path, err)
	}

	err = storage.UploadIfAbsent(
		ctx,
		client,
		bucket,
		key,
		snapshot,
	)
	if err != nil {
		var apiErr smithy.APIError

		if errors.As(err, &apiErr) &&
			apiErr.ErrorCode() == "PreconditionFailed" {

			exists, checkErr := storage.ObjectExists(
				ctx,
				client,
				bucket,
				key,
			)
			if checkErr != nil {
				return fmt.Errorf("recheck pack object %q: %w", key, checkErr)
			}

			if exists {
				return nil
			}
		}

		return fmt.Errorf("upload pack file %q: %w", file.Path, err)
	}

	return nil
}

// publishManifest сериализует манифест в JSON и публикует его в S3.
func publishManifest(
	ctx context.Context,
	client *s3.Client,
	bucket string,
	packManifest manifest.Manifest,
) error {
	data, err := json.Marshal(
		&packManifest,
		jsontext.WithIndent("  "),
	)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(data)

	err = storage.Upload(
		ctx,
		client,
		bucket,
		storage.ManifestKey,
		reader,
	)
	if err != nil {
		return fmt.Errorf("upload manifest: %w", err)
	}

	return nil
}

// verifyFileHash проверяет SHA-256 открытого файла и возвращает позицию чтения в начало.
func verifyFileHash(
	openedFile *os.File,
	expectedHash string,
) error {
	hasher := sha256.New()

	_, err := io.Copy(hasher, openedFile)
	if err != nil {
		return fmt.Errorf("hash file: %w", err)
	}

	actualHash := hex.EncodeToString(hasher.Sum(nil))
	if expectedHash != actualHash {
		return fmt.Errorf(
			"SHA-256 mismatch: expected %s, got %s",
			expectedHash,
			actualHash,
		)
	}

	_, err = openedFile.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("reset file position: %w", err)
	}

	return nil
}
