package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Environment represents a single Vault environment configuration.
type Environment struct {
	Name      string `yaml:"name"`
	Address   string `yaml:"address"`
	Token     string `yaml:"token"`
	Namespace string `yaml:"namespace,omitempty"`
}

// Config holds the top-level vaultdiff configuration.
type Config struct {
	Environments []Environment `yaml:"environments"`
	Paths        []string      `yaml:"paths"`
	MaskValues   bool          `yaml:"mask_values"`
}

// Load reads and parses a YAML config file from the given path.
// Token values may be overridden by environment variables of the form
// VAULTDIFF_TOKEN_<NAME> (uppercased environment name).
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	cfg.applyEnvOverrides()
	return &cfg, nil
}

// validate checks that required fields are present.
func (c *Config) validate() error {
	if len(c.Environments) < 2 {
		return fmt.Errorf("config must define at least 2 environments, got %d", len(c.Environments))
	}
	for i, env := range c.Environments {
		if env.Name == "" {
			return fmt.Errorf("environment[%d]: name is required", i)
		}
		if env.Address == "" {
			return fmt.Errorf("environment %q: address is required", env.Name)
		}
	}
	if len(c.Paths) == 0 {
		return fmt.Errorf("config must specify at least one secret path")
	}
	return nil
}

// applyEnvOverrides replaces empty tokens with values from the environment.
func (c *Config) applyEnvOverrides() {
	for i, env := range c.Environments {
		key := fmt.Sprintf("VAULTDIFF_TOKEN_%s", uppercaseEnvName(env.Name))
		if val := os.Getenv(key); val != "" {
			c.Environments[i].Token = val
		}
	}
}

func uppercaseEnvName(name string) string {
	out := make([]byte, len(name))
	for i := 0; i < len(name); i++ {
		ch := name[i]
		if ch >= 'a' && ch <= 'z' {
			ch -= 32
		}
		out[i] = ch
	}
	return string(out)
}
