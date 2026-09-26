package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/MeguruMacabre/meguru-pack-compact/internal/config"
	"github.com/MeguruMacabre/meguru-pack-compact/internal/instance"
)

func main() {
	root, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fullPath := filepath.Join(root, "host-config.json")
	fileIsExist, err := fileExists(fullPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if !fileIsExist {
		dirs, err := instance.FindCandidates(root)
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
		err = config.SaveHostConfig(root, cfg)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

	}

	cfg, err := config.LoadHostConfig(root)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	dirIsExist, err := dirExists(root, cfg.InstanceDirectory)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if !dirIsExist {
		fmt.Println("Error: no such directory")
		return
	}
	fmt.Println("Host Config:", cfg.InstanceDirectory)

	instancePath := filepath.Join(root, cfg.InstanceDirectory)

	gameRoot, err := instance.FindGameRoot(instancePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Game directory:", gameRoot)
}

func fileExists(filename string) (bool, error) {
	_, err := os.Stat(filename)
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
