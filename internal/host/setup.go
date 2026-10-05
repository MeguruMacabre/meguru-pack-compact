package host

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/MeguruMacabre/meguru-pack-compact/internal/config"
	"github.com/MeguruMacabre/meguru-pack-compact/internal/instance"
)

func loadOrCreateHostConfig(appRoot string) (config.HostConfig, error) {
	configPath := filepath.Join(appRoot, "host-config.json")

	configExists, err := fileExists(configPath)
	if err != nil {
		return config.HostConfig{}, err
	}

	if !configExists {
		dirs, err := instance.FindCandidates(appRoot)
		if err != nil {
			return config.HostConfig{}, err
		}
		if len(dirs) == 0 {
			return config.HostConfig{}, fmt.Errorf("no candidates found")
		}

		for i, dir := range dirs {
			fmt.Printf("%d: %s\n", i+1, dir)
		}

		fmt.Println()
		fmt.Print("Select instance: ")

		selectedCandidate, err := chooseCandidate(dirs)
		if err != nil {
			return config.HostConfig{}, err
		}

		fmt.Println("Selected:", selectedCandidate)
		cfg := config.HostConfig{
			InstanceDirectory: selectedCandidate,
		}

		err = config.SaveHostConfig(appRoot, cfg)
		if err != nil {
			return config.HostConfig{}, err
		}

	}

	cfg, err := config.LoadHostConfig(appRoot)
	if err != nil {
		return config.HostConfig{}, err
	}
	return cfg, nil
}

func chooseCandidate(dirs []string) (string, error) {
	var number int
	_, err := fmt.Scan(&number)
	if err != nil {
		return "", err
	}

	if number < 1 || number > len(dirs) {
		return "", fmt.Errorf("invalid candidate number: %d", number)
	}
	return dirs[number-1], nil
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
