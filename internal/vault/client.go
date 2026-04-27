package vault

import (
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

// Config holds configuration for connecting to a Vault instance.
type Config struct {
	Address   string
	Token     string
	Namespace string
}

// Client wraps the Vault API client with additional context.
type Client struct {
	api       *vaultapi.Client
	Namespace string
}

// NewClient creates and configures a new Vault client from the given Config.
// If Token or Address are empty, it falls back to environment variables.
func NewClient(cfg Config) (*Client, error) {
	vcfg := vaultapi.DefaultConfig()

	address := cfg.Address
	if address == "" {
		address = os.Getenv("VAULT_ADDR")
	}
	if address == "" {
		return nil, fmt.Errorf("vault address is required (set Address or VAULT_ADDR)")
	}
	vcfg.Address = address

	client, err := vaultapi.NewClient(vcfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	token := cfg.Token
	if token == "" {
		token = os.Getenv("VAULT_TOKEN")
	}
	if token == "" {
		return nil, fmt.Errorf("vault token is required (set Token or VAULT_TOKEN)")
	}
	client.SetToken(token)

	if cfg.Namespace != "" {
		client.SetNamespace(cfg.Namespace)
	}

	return &Client{
		api:       client,
		Namespace: cfg.Namespace,
	}, nil
}

// ReadSecret reads a KV v2 secret at the given path and returns its data map.
func (c *Client) ReadSecret(path string) (map[string]interface{}, error) {
	secret, err := c.api.Logical().Read(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret at %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("secret not found at path %q", path)
	}

	// KV v2 wraps data under secret.Data["data"]
	if data, ok := secret.Data["data"]; ok {
		if m, ok := data.(map[string]interface{}); ok {
			return m, nil
		}
	}

	return secret.Data, nil
}

// ListSecrets lists secret keys at the given path prefix.
func (c *Client) ListSecrets(path string) ([]string, error) {
	secret, err := c.api.Logical().List(path)
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets at %q: %w", path, err)
	}
	if secret == nil {
		return nil, nil
	}

	keys, ok := secret.Data["keys"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format listing %q", path)
	}

	result := make([]string, 0, len(keys))
	for _, k := range keys {
		if s, ok := k.(string); ok {
			result = append(result, s)
		}
	}
	return result, nil
}
