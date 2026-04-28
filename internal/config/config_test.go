package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.yaml")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_Success(t *testing.T) {
	content := `
environments:
  - name: staging
    address: https://vault-staging.example.com
    token: s.staging
  - name: prod
    address: https://vault-prod.example.com
    token: s.prod
paths:
  - secret/app/config
mask_values: true
`
	path := writeTemp(t, content)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Environments) != 2 {
		t.Errorf("expected 2 environments, got %d", len(cfg.Environments))
	}
	if !cfg.MaskValues {
		t.Error("expected mask_values to be true")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_TooFewEnvironments(t *testing.T) {
	content := `
environments:
  - name: staging
    address: https://vault-staging.example.com
    token: s.staging
paths:
  - secret/app/config
`
	path := writeTemp(t, content)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for fewer than 2 environments")
	}
}

func TestLoad_NoPaths(t *testing.T) {
	content := `
environments:
  - name: staging
    address: https://vault-staging.example.com
    token: s.staging
  - name: prod
    address: https://vault-prod.example.com
    token: s.prod
paths: []
`
	path := writeTemp(t, content)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for empty paths")
	}
}

func TestLoad_EnvTokenOverride(t *testing.T) {
	content := `
environments:
  - name: staging
    address: https://vault-staging.example.com
  - name: prod
    address: https://vault-prod.example.com
paths:
  - secret/app/config
`
	t.Setenv("VAULTDIFF_TOKEN_STAGING", "s.from-env")
	path := writeTemp(t, content)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Environments[0].Token != "s.from-env" {
		t.Errorf("expected token override, got %q", cfg.Environments[0].Token)
	}
}
