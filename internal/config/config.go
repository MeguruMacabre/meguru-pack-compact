package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type HostConfig struct {
	InstanceDirectory string `json:"instance_directory"`
}

func SaveHostConfig(root string, cfg HostConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	fullPath := filepath.Join(root, "host-config.json")
	err = os.WriteFile(fullPath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func LoadHostConfig(root string) (HostConfig, error) {
	fullPath := filepath.Join(root, "host-config.json")
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return HostConfig{}, err
	}

	var cfg HostConfig
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return HostConfig{}, err
	}
	return cfg, nil
}
