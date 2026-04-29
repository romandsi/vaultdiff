package config

import (
	"strings"
	"testing"
)

func validConfig() *Config {
	return &Config{
		Environments: []Environment{
			{Name: "staging", Address: "https://vault-staging.example.com", Token: "tok-a"},
			{Name: "production", Address: "https://vault-prod.example.com", Token: "tok-b"},
		},
		Paths: []string{"secret/data/app"},
	}
}

func TestValidate_Success(t *testing.T) {
	if err := Validate(validConfig()); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidate_TooFewEnvironments(t *testing.T) {
	cfg := validConfig()
	cfg.Environments = cfg.Environments[:1]
	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for single environment")
	}
	if !strings.Contains(err.Error(), "at least two") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidate_DuplicateEnvName(t *testing.T) {
	cfg := validConfig()
	cfg.Environments[1].Name = "staging"
	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for duplicate environment name")
	}
	if !strings.Contains(err.Error(), "duplicate") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidate_InvalidAddress(t *testing.T) {
	cfg := validConfig()
	cfg.Environments[0].Address = "not-a-url"
	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for invalid address")
	}
	if !strings.Contains(err.Error(), "valid URL") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidate_MissingToken(t *testing.T) {
	cfg := validConfig()
	cfg.Environments[1].Token = ""
	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for missing token")
	}
	if !strings.Contains(err.Error(), "token must not be empty") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidate_NoPaths(t *testing.T) {
	cfg := validConfig()
	cfg.Paths = nil
	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
	if !strings.Contains(err.Error(), "at least one secret path") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidate_MultipleIssues(t *testing.T) {
	cfg := &Config{}
	verr, ok := Validate(cfg).(*ValidationError)
	if !ok {
		t.Fatal("expected *ValidationError")
	}
	if len(verr.Issues) < 2 {
		t.Errorf("expected multiple issues, got %d", len(verr.Issues))
	}
}
