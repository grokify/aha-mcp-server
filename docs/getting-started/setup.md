# Setup

## Get Your Aha! Credentials

### API Token

1. Log in to your Aha! workspace
2. Go to **Settings** > **Personal** > **Developer** > **API**
3. Generate a new API token
4. Copy the token (it won't be shown again)

### Domain

Your Aha! subdomain is the first part of your workspace URL:

- If your workspace is at `mycompany.aha.io`, your domain is `mycompany`

## Configure Credentials

### Option 1: Environment Variables

```bash
export AHA_DOMAIN="your_subdomain"
export AHA_API_TOKEN="your_api_token"
```

### Option 2: Command-Line Flags

```bash
mcp-aha --subdomain mycompany --api-key your-api-key
```

### Option 3: Vault-Backed Credentials

For production use, store credentials in a vault:

```bash
# 1Password
export OP_SERVICE_ACCOUNT_TOKEN="ops_..."
mcp-aha --vault op://MyVault --credentials-name aha

# Bitwarden
export BW_ACCESS_TOKEN="..."
export BW_ORGANIZATION_ID="..."
mcp-aha --vault bw://org-id --credentials-name aha

# Keeper
export KSM_TOKEN="US:..."
mcp-aha --vault keeper:// --credentials-name aha
```

See [Credentials](../configuration/credentials.md) for detailed configuration options.

## Test Your Setup

```bash
# Test with CLI mode
mcp-aha search-documents "product roadmap" --subdomain mycompany --api-key your-key
```
