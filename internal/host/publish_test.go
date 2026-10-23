package host

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVerifyFileHashValid(t *testing.T) {
	content := []byte("Hello Meguru")
	filePath := filepath.Join(t.TempDir(), "test.txt")

	err := os.WriteFile(filePath, content, 0644)
	if err != nil {
		t.Fatal(err)
	}

	openedFile, err := os.Open(filePath)
	if err != nil {
		t.Fatal(err)
	}
	defer openedFile.Close()

	hash := sha256.Sum256(content)
	expectedHash := hex.EncodeToString(hash[:])

	err = verifyFileHash(openedFile, expectedHash)
	if err != nil {
		t.Fatalf("verifyFileHash() unexpected error: %v", err)
	}

	position, err := openedFile.Seek(0, io.SeekCurrent)
	if err != nil {
		t.Fatal(err)
	}

	if position != 0 {
		t.Errorf("file position = %d, want 0", position)
	}
}

func TestVerifyFileHashMismatch(t *testing.T) {
	content := []byte("Hello Meguru")
	filePath := filepath.Join(t.TempDir(), "test.txt")

	err := os.WriteFile(filePath, content, 0644)
	if err != nil {
		t.Fatal(err)
	}

	openedFile, err := os.Open(filePath)
	if err != nil {
		t.Fatal(err)
	}
	defer openedFile.Close()

	expectedHash := strings.Repeat("0", 64)

	err = verifyFileHash(openedFile, expectedHash)
	if err == nil {
		t.Fatal("verifyFileHash() expected mismatch error, got nil")
	}

	if !strings.Contains(err.Error(), "SHA-256 mismatch") {
		t.Errorf("unexpected error: %v", err)
	}
}
