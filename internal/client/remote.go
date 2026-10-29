package client

import (
	"context"
	"encoding/json/v2"
	"fmt"
	"io"
	"net/http"

	"github.com/MeguruMacabre/meguru-pack-compact/internal/manifest"
)

// fetchManifest загружает манифест по HTTP, проверяет статус и размер ответа и декодирует JSON.
func fetchManifest(
	ctx context.Context,
	httpClient *http.Client,
	url string,
) (manifest.Manifest, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return manifest.Manifest{}, fmt.Errorf("create manifest request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return manifest.Manifest{}, fmt.Errorf("send manifest request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return manifest.Manifest{}, fmt.Errorf("unexpected HTTP status: %d", resp.StatusCode)
	}

	const maxSize = 2 << 20

	limited := io.LimitReader(resp.Body, maxSize+1)

	data, err := io.ReadAll(limited)
	if err != nil {
		return manifest.Manifest{}, err
	}

	if len(data) > maxSize {
		return manifest.Manifest{}, fmt.Errorf("manifest exceeds maximum size")
	}

	var result manifest.Manifest

	err = json.Unmarshal(data, &result)
	if err != nil {
		return manifest.Manifest{}, fmt.Errorf("decode manifest: %w", err)
	}

	return result, nil
}
