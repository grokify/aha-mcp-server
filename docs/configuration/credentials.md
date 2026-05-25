# Credentials

The Aha! MCP Server supports multiple credential sources for authentication.

## Option 1: Direct Credentials

The simplest option - provide your Aha! subdomain and API token directly.

### Setup

1. Log in to your Aha! workspace
2. Go to **Settings** > **Personal** > **Developer** > **API**
3. Generate a new API token

### Usage

```bash
mcp-aha --subdomain mycompany --api-key your-api-key
```

Or with environment variables:

```bash
export AHA_DOMAIN=mycompany
export AHA_API_TOKEN=your-api-key
mcp-aha
```

## Option 2: Vault-Backed Credentials

Use [omnitoken](https://github.com/plexusone/omnitoken) with vault backends for secure credential storage.

### Supported Vault URIs

| Provider | URI Pattern | Requirements |
|----------|-------------|--------------|
| 1Password | `op://vault` | `OP_SERVICE_ACCOUNT_TOKEN` env var |
| Bitwarden | `bw://org-id` | `BW_ACCESS_TOKEN` and `BW_ORGANIZATION_ID` env vars |
| Keeper | `keeper://` | `KSM_TOKEN` or `KSM_CONFIG` env var |
| File | `file:///path/to/dir` | None |

### 1Password

Store your credentials in 1Password and access them securely:

```bash
export OP_SERVICE_ACCOUNT_TOKEN="ops_..."
mcp-aha --vault op://MyVault --credentials-name aha
```

The credential item should contain:

- `subdomain` or `domain` field
- `api_key` or `api_token` field with Authorization header

### Bitwarden

Store credentials in Bitwarden Secrets Manager:

```bash
export BW_ACCESS_TOKEN="..."
export BW_ORGANIZATION_ID="..."
mcp-aha --vault bw://org-id --credentials-name aha
```

### Keeper

Store credentials in Keeper Secrets Manager:

```bash
export KSM_TOKEN="US:..."
mcp-aha --vault keeper:// --credentials-name aha
```

### File Vault

For local development:

```bash
mcp-aha --vault file:///path/to/secrets --credentials-name aha
```

Create a file at `/path/to/secrets/aha.json`:

```json
{
  "type": "headerquery",
  "subdomain": "mycompany",
  "headerQuery": {
    "header": {
      "Authorization": ["Bearer your-api-token"]
    }
  }
}
```

## Credential Format

When using vault storage, credentials should be in goauth format:

```json
{
  "type": "headerquery",
  "subdomain": "mycompany",
  "headerQuery": {
    "serverURL": "https://mycompany.aha.io",
    "header": {
      "Authorization": ["Bearer your-api-token"]
    }
  }
}
```

Or with separate fields:

| Field | Description |
|-------|-------------|
| `subdomain` | Aha! subdomain |
| `api_key` / `api_token` | API token |
| `Authorization` header | Bearer token |

## Security Best Practices

1. **Never commit credentials** - Add credentials files to `.gitignore`
2. **Use vault backends** - For production, use proper secrets management
3. **Rotate keys** - Periodically rotate API tokens
4. **Use file permissions** - `chmod 600` for credential files
