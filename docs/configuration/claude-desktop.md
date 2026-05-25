# Claude Desktop Configuration

Configure Claude Desktop to use the Aha! MCP Server.

## Configuration File Location

| OS | Path |
|----|------|
| macOS | `~/Library/Application Support/Claude/claude_desktop_config.json` |
| Windows | `%APPDATA%\Claude\claude_desktop_config.json` |
| Linux | `~/.config/Claude/claude_desktop_config.json` |

## Basic Configuration

### With Direct Credentials

```json
{
  "mcpServers": {
    "aha": {
      "command": "/path/to/aha-mcp-server",
      "env": {
        "AHA_DOMAIN": "your_subdomain",
        "AHA_API_TOKEN": "your_api_token"
      }
    }
  }
}
```

### With 1Password

```json
{
  "mcpServers": {
    "aha": {
      "command": "/path/to/aha-mcp-server",
      "env": {
        "OP_SERVICE_ACCOUNT_TOKEN": "ops_...",
        "OMNITOKEN_VAULT_URI": "op://MyVault",
        "OMNITOKEN_CREDENTIALS_NAME": "aha"
      }
    }
  }
}
```

### With Bitwarden

```json
{
  "mcpServers": {
    "aha": {
      "command": "/path/to/aha-mcp-server",
      "env": {
        "BW_ACCESS_TOKEN": "...",
        "BW_ORGANIZATION_ID": "...",
        "OMNITOKEN_VAULT_URI": "bw://org-id",
        "OMNITOKEN_CREDENTIALS_NAME": "aha"
      }
    }
  }
}
```

### With Keeper

```json
{
  "mcpServers": {
    "aha": {
      "command": "/path/to/aha-mcp-server",
      "env": {
        "KSM_TOKEN": "US:...",
        "OMNITOKEN_VAULT_URI": "keeper://",
        "OMNITOKEN_CREDENTIALS_NAME": "aha"
      }
    }
  }
}
```

### With File Vault

```json
{
  "mcpServers": {
    "aha": {
      "command": "/path/to/aha-mcp-server",
      "args": [
        "--vault", "file:///path/to/secrets",
        "--credentials-name", "aha"
      ]
    }
  }
}
```

## Environment Variables Reference

| Variable | Description |
|----------|-------------|
| `AHA_DOMAIN` | Aha! subdomain |
| `AHA_API_TOKEN` | Aha! API key |
| `OMNITOKEN_VAULT_URI` | Vault URI (e.g., `op://MyVault`) |
| `OMNITOKEN_CREDENTIALS_NAME` | Credential name in vault (default: `aha`) |
| `OP_SERVICE_ACCOUNT_TOKEN` | 1Password service account token |
| `BW_ACCESS_TOKEN` | Bitwarden access token |
| `BW_ORGANIZATION_ID` | Bitwarden organization ID |
| `KSM_TOKEN` | Keeper token (format: `REGION:TOKEN`) |

## Multiple Servers

You can run multiple MCP servers alongside Aha!:

```json
{
  "mcpServers": {
    "aha": {
      "command": "/path/to/aha-mcp-server",
      "env": {
        "AHA_DOMAIN": "mycompany",
        "AHA_API_TOKEN": "aha-token"
      }
    },
    "google": {
      "command": "/path/to/mcp-google",
      "env": {
        "GOOGLE_CREDENTIALS_FILE": "/path/to/service-account.json"
      }
    },
    "confluence": {
      "command": "/path/to/mcp-confluence",
      "env": {
        "CONFLUENCE_BASE_URL": "https://example.atlassian.net/wiki",
        "CONFLUENCE_USERNAME": "user@example.com",
        "CONFLUENCE_API_TOKEN": "token"
      }
    }
  }
}
```

## Finding the Binary Path

### If installed with `go install`

```bash
# Find GOPATH
go env GOPATH

# Binary is at $GOPATH/bin/aha-mcp-server
# Typically: ~/go/bin/aha-mcp-server
```

### If built from source

Use the full path to where you built it:

```bash
/path/to/aha-mcp-server/aha-mcp-server
```

## Troubleshooting

### Server Not Starting

Check the Claude Desktop logs:

- macOS: `~/Library/Logs/Claude/`
- Windows: `%APPDATA%\Claude\logs\`

Common issues:

1. **Binary not found**: Verify the path is correct
2. **Credentials not found**: Check environment variables
3. **Permission denied**: Ensure the binary is executable (`chmod +x`)

### Verifying Configuration

Test the server manually:

```bash
# Should start and wait for input (Ctrl+C to exit)
/path/to/aha-mcp-server --subdomain mycompany --api-key your-key
```

### JSON Syntax Errors

Validate your JSON:

```bash
# On macOS/Linux
cat ~/Library/Application\ Support/Claude/claude_desktop_config.json | python3 -m json.tool
```

## Available Tools in Claude

Once configured, you can ask Claude to:

- "Search Aha! for documents about [topic]"
- "Get the feature FEAT-123"
- "List ideas tagged with 'mobile'"
- "Show me the epic EPIC-456"
- "Get user details for USER-789"
