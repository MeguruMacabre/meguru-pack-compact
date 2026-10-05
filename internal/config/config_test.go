package config

import "testing"

func TestValidateS3ConfigValid(t *testing.T) {
	cfg := HostConfig{
		S3Endpoint: "Endpoint",
		S3Bucket:   "Bucket",
		S3Region:   "Region",
	}

	err := ValidateS3Config(cfg)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestValidateS3ConfigMissingEndpoint(t *testing.T) {
	cfg := HostConfig{
		S3Endpoint: "",
		S3Bucket:   "Bucket",
		S3Region:   "Region",
	}

	err := ValidateS3Config(cfg)
	if err == nil {
		t.Error("expected error for missing S3 endpoint")
	}
}

func TestValidateS3ConfigMissingBucket(t *testing.T) {
	cfg := HostConfig{
		S3Endpoint: "Endpoint",
		S3Bucket:   "",
		S3Region:   "Region",
	}

	err := ValidateS3Config(cfg)
	if err == nil {
		t.Error("expected error for missing S3 bucket")
	}
}

func TestValidateS3ConfigMissingRegion(t *testing.T) {
	cfg := HostConfig{
		S3Endpoint: "Endpoint",
		S3Bucket:   "Bucket",
		S3Region:   "",
	}

	err := ValidateS3Config(cfg)
	if err == nil {
		t.Error("expected error for missing S3 region")
	}
}
