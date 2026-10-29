package client

import (
	"context"
	"encoding/json/v2"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MeguruMacabre/meguru-pack-compact/internal/manifest"
)

func TestFetchManifestSuccess(t *testing.T) {
	expected := manifest.Manifest{
		FormatVersion: 1,
		Files:         []manifest.File{},
	}

	data, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("marshal test manifest: %v", err)
	}

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			_, _ = w.Write(data)
		}),
	)
	defer server.Close()

	result, err := fetchManifest(
		context.Background(),
		server.Client(),
		server.URL,
	)
	if err != nil {
		t.Fatalf("fetchManifest returned error: %v", err)
	}

	if result.FormatVersion != expected.FormatVersion {
		t.Errorf(
			"expected format version %d, got %d",
			expected.FormatVersion,
			result.FormatVersion,
		)
	}

	if len(result.Files) != 0 {
		t.Errorf("expected 0 files, got %d", len(result.Files))
	}
}

func TestFetchManifestNotFound(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		}),
	)
	defer server.Close()

	_, err := fetchManifest(
		context.Background(),
		server.Client(),
		server.URL,
	)

	if err == nil {
		t.Fatal("expected an error for HTTP 404, got nil")
	}

	if !strings.Contains(err.Error(), "404") {
		t.Errorf("expected HTTP 404 error, got: %v", err)
	}
}

func TestFetchManifestInvalidJSON(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"format_version":`))
		}),
	)
	defer server.Close()

	_, err := fetchManifest(
		context.Background(),
		server.Client(),
		server.URL,
	)

	if err == nil {
		t.Fatal("expected JSON decoding error")
	}

	if !strings.Contains(err.Error(), "decode manifest") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestFetchManifestTooLarge(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(strings.Repeat("x", (2<<20)+1)))
		}),
	)
	defer server.Close()

	_, err := fetchManifest(
		context.Background(),
		server.Client(),
		server.URL,
	)

	if err == nil {
		t.Fatal("expected manifest size limit error")
	}

	if !strings.Contains(err.Error(), "exceeds maximum size") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestFetchManifestCanceledContext(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := fetchManifest(ctx, server.Client(), server.URL)

	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got: %v", err)
	}
}
