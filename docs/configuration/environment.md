# Environment Variables

All command-line flags can be set via environment variables.

## Available Variables

### Credential Configuration

| Variable | Flag | Description |
|----------|------|-------------|
| `AHA_DOMAIN` | `--subdomain` | Aha! subdomain |
| `AHA_API_TOKEN` | `--api-key` | Aha! API key |
| `OMNITOKEN_VAULT_URI` | `--vault` | Vault URI for credentials |
| `OMNITOKEN_CREDENTIALS_NAME` | `--credentials-name` | Name of credentials in vault |

### Vault Provider Authentication

| Variable | Description |
|----------|-------------|
| `OP_SERVICE_ACCOUNT_TOKEN` | 1Password service account token |
| `BW_ACCESS_TOKEN` | Bitwarden access token |
| `BW_ORGANIZATION_ID` | Bitwarden organization ID |
| `KSM_TOKEN` | Keeper token (format: `REGION:TOKEN`) |
| `KSM_CONFIG` | Keeper config (base64-encoded JSON) |

## Precedence

Command-line flags take precedence over environment variables.

```bash
# Environment variable is used
export AHA_DOMAIN=mycompany
aha-mcp-server
# Uses: mycompany

# Flag overrides environment
export AHA_DOMAIN=mycompany
aha-mcp-server --subdomain othercompany
# Uses: othercompany
```

## Examples

### Direct Credentials

```bash
export AHA_DOMAIN=mycompany
export AHA_API_TOKEN=your-api-key
aha-mcp-server
```

### 1Password

```bash
export OP_SERVICE_ACCOUNT_TOKEN="ops_..."
export OMNITOKEN_VAULT_URI=op://MyVault
export OMNITOKEN_CREDENTIALS_NAME=aha
aha-mcp-server
```

### Bitwarden

```bash
export BW_ACCESS_TOKEN="..."
export BW_ORGANIZATION_ID="..."
export OMNITOKEN_VAULT_URI=bw://org-id
export OMNITOKEN_CREDENTIALS_NAME=aha
aha-mcp-server
```

### Keeper

```bash
export KSM_TOKEN="US:..."
export OMNITOKEN_VAULT_URI=keeper://
export OMNITOKEN_CREDENTIALS_NAME=aha
aha-mcp-server
```

## Shell Configuration

### Bash/Zsh

Add to `~/.bashrc` or `~/.zshrc`:

```bash
# Aha! MCP Server credentials
export AHA_DOMAIN="mycompany"
export AHA_API_TOKEN="your-api-key"
```

### Fish

Add to `~/.config/fish/config.fish`:

```fish
set -gx AHA_DOMAIN "mycompany"
set -gx AHA_API_TOKEN "your-api-key"
```

## Dotenv Files

For project-specific configuration, use a `.env` file:

```bash
# .env
AHA_DOMAIN=mycompany
AHA_API_TOKEN=your-api-key
```

Load with `source` or a dotenv tool:

```bash
source .env
aha-mcp-server
```

Or use a tool like `direnv` for automatic loading.
