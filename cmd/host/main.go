package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/MeguruMacabre/meguru-pack-compact/internal/config"
	"github.com/MeguruMacabre/meguru-pack-compact/internal/instance"
	"github.com/MeguruMacabre/meguru-pack-compact/internal/scanner"
	"github.com/MeguruMacabre/meguru-pack-compact/internal/syncpolicy"
)

func main() {
	appRoot, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	configPath := filepath.Join(appRoot, "host-config.json")

	configExists, err := fileExists(configPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if !configExists {
		dirs, err := instance.FindCandidates(appRoot)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if len(dirs) == 0 {
			fmt.Println("Error: no candidates found")
			return
		}

		for i, dir := range dirs {
			fmt.Printf("%d: %s\n", i+1, dir)
		}

		fmt.Println()
		fmt.Print("Select instance: ")

		selectedCandidate, err := chooseCandidate(dirs)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Selected:", selectedCandidate)
		cfg := config.HostConfig{
			InstanceDirectory: selectedCandidate,
		}
		err = config.SaveHostConfig(appRoot, cfg)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

	}

	cfg, err := config.LoadHostConfig(appRoot)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	instanceExists, err := dirExists(appRoot, cfg.InstanceDirectory)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if !instanceExists {
		fmt.Println("Error: no such directory")
		return
	}
	fmt.Println("Host Config:", cfg.InstanceDirectory)

	instanceRoot := filepath.Join(appRoot, cfg.InstanceDirectory)

	gameRoot, err := instance.FindGameRoot(instanceRoot)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Game directory:", gameRoot)

	scannedFiles, err := scanner.Scan(instanceRoot)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, file := range scannedFiles {
		absoluteFilePath := filepath.Join(instanceRoot, file.Path)
		isManaged, err := syncpolicy.IsManaged(gameRoot, absoluteFilePath)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Printf("%s - %d bytes | %t | %s\n", file.Path, file.Size, isManaged, file.Hash)
	}
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func dirExists(root string, name string) (bool, error) {
	fullPath := filepath.Join(root, name)

	info, err := os.Stat(fullPath)
	if err == nil {
		return info.IsDir(), nil
	}

	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}

	return false, err
}

func chooseCandidate(dirs []string) (string, error) {
	var number int
	_, err := fmt.Scan(&number)
	if err != nil {
		return "", err
	}

	length := len(dirs)
	if number < 1 || number > length {
		return "", fmt.Errorf("invalid candidate number: %d", number)
	}
	return dirs[number-1], nil
}
