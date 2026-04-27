# vaultdiff

> CLI tool to compare and audit secrets across multiple HashiCorp Vault namespaces or environments

---

## Installation

```bash
go install github.com/yourusername/vaultdiff@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/vaultdiff/releases).

---

## Usage

Set your Vault credentials and compare secrets between two environments:

```bash
export VAULT_TOKEN=<your-token>

# Compare secrets at a path between two Vault addresses
vaultdiff \
  --source https://vault.staging.example.com \
  --target https://vault.prod.example.com \
  --path secret/data/myapp
```

**Example output:**

```
[~] secret/data/myapp/config
    DB_HOST   staging-db.internal  →  prod-db.internal
    LOG_LEVEL debug                →  info

[+] secret/data/myapp/feature-flags  (only in target)
[-] secret/data/myapp/legacy         (only in source)
```

### Flags

| Flag | Description |
|------|-------------|
| `--source` | Source Vault address |
| `--target` | Target Vault address |
| `--path` | Secret path to compare |
| `--namespace` | Vault namespace (optional) |
| `--output` | Output format: `text`, `json`, `yaml` |

---

## License

[MIT](LICENSE)