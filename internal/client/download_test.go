package client

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func testSHA256(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func assertTempDirEmpty(t *testing.T, dir string) {
	t.Helper()

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read temporary directory: %v", err)
	}

	if len(entries) != 0 {
		t.Fatalf("temporary files were not cleaned up: %v", entries)
	}
}

func TestDownloadObjectSuccess(t *testing.T) {
	content := []byte("Meguru Pack Compact test file")

	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}

			_, _ = w.Write(content)
		},
	))
	defer server.Close()

	path, err := downloadObject(
		context.Background(),
		server.Client(),
		server.URL,
		testSHA256(content),
		int64(len(content)),
	)
	if err != nil {
		t.Fatalf("downloadObject failed: %v", err)
	}
	defer os.Remove(path)

	downloaded, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read downloaded file: %v", err)
	}

	if !bytes.Equal(downloaded, content) {
		t.Fatal("downloaded content does not match expected content")
	}
}

func TestDownloadObjectWrongSize(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("TMPDIR", tempDir)

	content := []byte("hello world")

	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write(content)
		},
	))
	defer server.Close()

	_, err := downloadObject(
		context.Background(),
		server.Client(),
		server.URL,
		testSHA256(content),
		5,
	)

	if err == nil || !strings.Contains(err.Error(), "unexpected file size") {
		t.Fatalf("expected file size error, got: %v", err)
	}

	assertTempDirEmpty(t, tempDir)
}

func TestDownloadObjectWrongHash(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("TMPDIR", tempDir)

	content := []byte("hello world")

	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write(content)
		},
	))
	defer server.Close()

	_, err := downloadObject(
		context.Background(),
		server.Client(),
		server.URL,
		strings.Repeat("0", 64),
		int64(len(content)),
	)

	if err == nil || !strings.Contains(err.Error(), "SHA-256 mismatch") {
		t.Fatalf("expected SHA-256 mismatch, got: %v", err)
	}

	assertTempDirEmpty(t, tempDir)
}

func TestDownloadObjectNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		},
	))
	defer server.Close()

	_, err := downloadObject(
		context.Background(),
		server.Client(),
		server.URL,
		strings.Repeat("0", 64),
		10,
	)

	if err == nil || !strings.Contains(err.Error(), "404") {
		t.Fatalf("expected HTTP 404 error, got: %v", err)
	}
}

func TestDownloadObjectInterrupted(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("TMPDIR", tempDir)

	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "100")
			w.WriteHeader(http.StatusOK)

			_, _ = w.Write([]byte("partial"))
		},
	))
	defer server.Close()

	_, err := downloadObject(
		context.Background(),
		server.Client(),
		server.URL,
		strings.Repeat("0", 64),
		100,
	)

	if err == nil {
		t.Fatal("expected interrupted download error")
	}

	assertTempDirEmpty(t, tempDir)
}
