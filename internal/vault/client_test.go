package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newMockVaultServer(t *testing.T, path string, response map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("failed to encode mock response: %v", err)
		}
	}))
}

func TestNewClient_MissingAddress(t *testing.T) {
	t.Setenv("VAULT_ADDR", "")
	t.Setenv("VAULT_TOKEN", "")

	_, err := NewClient(Config{})
	if err == nil {
		t.Fatal("expected error when address is missing, got nil")
	}
}

func TestNewClient_MissingToken(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "")

	_, err := NewClient(Config{Address: "http://127.0.0.1:8200"})
	if err == nil {
		t.Fatal("expected error when token is missing, got nil")
	}
}

func TestNewClient_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClient(Config{
		Address:   server.URL,
		Token:     "test-token",
		Namespace: "dev",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
	if client.Namespace != "dev" {
		t.Errorf("expected namespace %q, got %q", "dev", client.Namespace)
	}
}

func TestReadSecret_KVv2(t *testing.T) {
	mockResponse := map[string]interface{}{
		"data": map[string]interface{}{
			"data": map[string]interface{}{
				"username": "admin",
				"password": "s3cr3t",
			},
		},
	}
	server := newMockVaultServer(t, "/v1/secret/data/myapp", mockResponse)
	defer server.Close()

	client, err := NewClient(Config{Address: server.URL, Token: "tok"})
	if err != nil {
		t.Fatalf("client creation failed: %v", err)
	}

	data, err := client.ReadSecret("secret/data/myapp")
	if err != nil {
		t.Fatalf("ReadSecret failed: %v", err)
	}
	if data["username"] != "admin" {
		t.Errorf("expected username=admin, got %v", data["username"])
	}
}
