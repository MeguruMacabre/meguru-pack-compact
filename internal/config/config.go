package config

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"os"
	"path/filepath"
)

type HostConfig struct {
	InstanceDirectory string `json:"instance_directory"`
	S3Endpoint        string `json:"s3_endpoint"`
	S3Region          string `json:"s3_region"`
	S3Bucket          string `json:"s3_bucket"`
}

type S3Credentials struct {
	AccessKey string
	SecretKey string
}

func LoadS3Credentials() (S3Credentials, error) {
	var credentials S3Credentials

	accessKey, exists := os.LookupEnv("S3_ACCESS_KEY")
	if !exists || accessKey == "" {
		return S3Credentials{}, fmt.Errorf("missing environment variable: S3_ACCESS_KEY")
	}

	secretKey, exists := os.LookupEnv("S3_SECRET_KEY")
	if !exists || secretKey == "" {
		return S3Credentials{}, fmt.Errorf("missing environment variable: S3_SECRET_KEY")
	}

	credentials.AccessKey = accessKey
	credentials.SecretKey = secretKey

	return credentials, nil
}

func SaveHostConfig(root string, cfg HostConfig) error {
	data, err := json.Marshal(&cfg, jsontext.WithIndent("  "))
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

func ValidateS3Config(cfg HostConfig) error {
	if cfg.S3Endpoint == "" {
		return fmt.Errorf("s3_endpoint is required")
	}

	if cfg.S3Region == "" {
		return fmt.Errorf("s3_region is required")
	}

	if cfg.S3Bucket == "" {
		return fmt.Errorf("s3_bucket is required")
	}

	return nil
}
