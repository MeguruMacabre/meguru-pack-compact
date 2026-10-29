package client

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
)

// downloadObject скачивает объект во временный файл и проверяет его размер и SHA-256.
// Возвращает путь к проверенному файлу, который вызывающая функция должна удалить после использования.
func downloadObject(
	ctx context.Context,
	httpClient *http.Client,
	objectURL string,
	expectedSHA256 string,
	expectedSize int64,
) (string, error) {
	if expectedSize < 0 || expectedSize == (1<<63)-1 {
		return "", fmt.Errorf("invalid expected file size: %d", expectedSize)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		objectURL,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("create download request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("download object: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"unexpected HTTP status: %s",
			resp.Status,
		)
	}

	file, err := os.CreateTemp("", "meguru-pack-*.part")
	if err != nil {
		return "", fmt.Errorf("create temporary file: %w", err)
	}

	success := false

	defer func() {
		if !success {
			_ = file.Close()
			_ = os.Remove(file.Name())
		}
	}()

	hasher := sha256.New()
	writer := io.MultiWriter(file, hasher)

	limited := io.LimitReader(resp.Body, expectedSize+1)

	written, err := io.Copy(writer, limited)
	if err != nil {
		return "", fmt.Errorf("write downloaded object: %w", err)
	}

	if written != expectedSize {
		return "", fmt.Errorf(
			"unexpected file size: expected %d, got %d",
			expectedSize,
			written,
		)
	}

	actualSHA256 := hex.EncodeToString(hasher.Sum(nil))

	if actualSHA256 != expectedSHA256 {
		return "", fmt.Errorf(
			"SHA-256 mismatch: expected %s, got %s",
			expectedSHA256,
			actualSHA256,
		)
	}

	if err := file.Close(); err != nil {
		return "", fmt.Errorf("close downloaded file: %w", err)
	}

	success = true

	return file.Name(), nil
}
