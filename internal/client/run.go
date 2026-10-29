package client

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Run создаёт HTTP-клиент, загружает удалённый манифест и выводит информацию о сборке.
func Run(ctx context.Context, manifestURL string) error {
	httpClient := &http.Client{
		Timeout: 15 * time.Second,
	}
	manifest, err := fetchManifest(ctx, httpClient, manifestURL)
	if err != nil {
		return fmt.Errorf("fetch manifest: %w", err)
	}
	fmt.Println(manifest.FormatVersion, len(manifest.Files))
	return nil
}
