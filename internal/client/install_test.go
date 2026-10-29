package client

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func createInstallTestRoot(t *testing.T) (*os.Root, string) {
	t.Helper()

	dir := t.TempDir()

	root, err := os.OpenRoot(dir)
	if err != nil {
		t.Fatalf("open test root: %v", err)
	}

	t.Cleanup(func() {
		if err := root.Close(); err != nil {
			t.Errorf("close test root: %v", err)
		}
	})

	return root, dir
}

func createInstallTestSource(t *testing.T, content []byte) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "download.part")

	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatalf("create source file: %v", err)
	}

	return path
}

func TestInstallFileSuccess(t *testing.T) {
	root, dir := createInstallTestRoot(t)

	content := []byte("Meguru Pack Compact")
	source := createInstallTestSource(t, content)

	err := installFile(root, source, ".minecraft/mods/test.jar")
	if err != nil {
		t.Fatalf("installFile failed: %v", err)
	}

	destination := filepath.Join(
		dir, ".minecraft", "mods", "test.jar",
	)

	actual, err := os.ReadFile(destination)
	if err != nil {
		t.Fatalf("read installed file: %v", err)
	}

	if !bytes.Equal(actual, content) {
		t.Fatal("installed file content does not match")
	}

	if _, err := os.Stat(source); err != nil {
		t.Fatalf("source file should remain available: %v", err)
	}
}

func TestInstallFileReplace(t *testing.T) {
	root, dir := createInstallTestRoot(t)

	destination := filepath.Join(dir, "mods", "test.jar")

	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(destination, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}

	source := createInstallTestSource(t, []byte("new"))

	if err := installFile(root, source, "mods/test.jar"); err != nil {
		t.Fatalf("replace file: %v", err)
	}

	actual, err := os.ReadFile(destination)
	if err != nil {
		t.Fatal(err)
	}

	if string(actual) != "new" {
		t.Fatalf("expected new content, got %q", actual)
	}
}

func TestInstallFileCreatesDirectories(t *testing.T) {
	root, dir := createInstallTestRoot(t)

	source := createInstallTestSource(t, []byte("nested file"))

	err := installFile(
		root,
		source,
		".minecraft/mods/nested/example.jar",
	)
	if err != nil {
		t.Fatalf("install file: %v", err)
	}

	destination := filepath.Join(
		dir, ".minecraft", "mods", "nested", "example.jar",
	)

	if _, err := os.Stat(destination); err != nil {
		t.Fatalf("installed file does not exist: %v", err)
	}
}

func TestInstallFileInvalidPath(t *testing.T) {
	root, _ := createInstallTestRoot(t)
	source := createInstallTestSource(t, []byte("test"))

	invalidPaths := []string{
		"",
		"../secret.txt",
		"/etc/passwd",
		"mods//test.jar",
		`mods\test.jar`,
	}

	for _, path := range invalidPaths {
		t.Run(path, func(t *testing.T) {
			err := installFile(root, source, path)
			if err == nil {
				t.Fatalf("expected error for path %q", path)
			}
		})
	}
}

func TestInstallFileCleansUpOnError(t *testing.T) {
	root, dir := createInstallTestRoot(t)

	modsDir := filepath.Join(dir, "mods")
	destination := filepath.Join(modsDir, "test.jar")

	// Создаём директорию на месте будущего файла.
	// Переименование файла поверх директории должно завершиться ошибкой.
	if err := os.MkdirAll(destination, 0o755); err != nil {
		t.Fatal(err)
	}

	source := createInstallTestSource(t, []byte("test"))

	err := installFile(root, source, "mods/test.jar")
	if err == nil {
		t.Fatal("expected installation error")
	}

	entries, err := os.ReadDir(modsDir)
	if err != nil {
		t.Fatal(err)
	}

	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".meguru-") {
			t.Fatalf("temporary file was not removed: %s", entry.Name())
		}
	}
}

func TestInstallFileSymlinkEscape(t *testing.T) {
	root, dir := createInstallTestRoot(t)
	outsideDir := t.TempDir()

	linkPath := filepath.Join(dir, "mods")

	if err := os.Symlink(outsideDir, linkPath); err != nil {
		t.Skipf("cannot create symlink: %v", err)
	}

	source := createInstallTestSource(t, []byte("test"))

	err := installFile(root, source, "mods/escape.jar")
	if err == nil {
		t.Fatal("expected symlink escape to be rejected")
	}

	outsideFile := filepath.Join(outsideDir, "escape.jar")

	if _, err := os.Stat(outsideFile); !os.IsNotExist(err) {
		t.Fatalf("unexpected result for outside file: %v", err)
	}

	entries, err := os.ReadDir(outsideDir)
	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 0 {
		t.Fatalf("files were created outside instance: %v", entries)
	}
}
