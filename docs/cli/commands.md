# CLI Commands

The Aha! MCP server can also be used as a command-line tool for testing and scripting.

## Global Flags

| Flag | Environment Variable | Description |
|------|---------------------|-------------|
| `--subdomain` | `AHA_DOMAIN` | Aha! subdomain |
| `--api-key` | `AHA_API_TOKEN` | Aha! API key |
| `--vault` | `OMNITOKEN_VAULT_URI` | Vault URI for credentials |
| `--credentials-name` | `OMNITOKEN_CREDENTIALS_NAME` | Name of credentials in vault |
| `-o, --output` | - | Output format: json (default) or pretty |

## Server Commands

### serve

Start the MCP server (default when no command specified).

```bash
mcp-aha serve
mcp-aha  # Same as above
```

### version

Print version information.

```bash
mcp-aha version
```

## Aha! Commands

### search-documents

Search for documents.

```bash
mcp-aha search-documents "product roadmap"
mcp-aha search-documents "auth" --type Page --limit 5
```

| Flag | Description |
|------|-------------|
| `--type` | Document type (e.g., Page) |
| `--limit` | Maximum results |

### get-idea

Get an idea by ID.

```bash
mcp-aha get-idea IDEA-123
mcp-aha get-idea IDEA-123 --output pretty
```

### list-ideas

List ideas with filtering.

```bash
mcp-aha list-ideas
mcp-aha list-ideas --query "mobile" --workflow-status "Under consideration"
mcp-aha list-ideas --tag "priority" --sort recent --per-page 20
```

| Flag | Description |
|------|-------------|
| `-q, --query` | Search term |
| `--spam` | Show spam ideas |
| `--workflow-status` | Filter by status |
| `--sort` | Sort by: recent, trending, popular |
| `--created-before` | Created before date (ISO8601) |
| `--created-since` | Created after date (ISO8601) |
| `--updated-since` | Updated after date (ISO8601) |
| `--tag` | Filter by tag |
| `--user-id` | Filter by creator |
| `--page` | Page number |
| `--per-page` | Results per page |

### get-feature

Get a feature by ID.

```bash
mcp-aha get-feature FEAT-123
```

### get-epic

Get an epic by ID.

```bash
mcp-aha get-epic EPIC-456
```

### get-release

Get a release by ID.

```bash
mcp-aha get-release REL-789
```

### get-goal

Get a goal by ID.

```bash
mcp-aha get-goal GOAL-123
```

### get-initiative

Get an initiative by ID.

```bash
mcp-aha get-initiative INIT-456
```

### get-key-result

Get a key result by ID.

```bash
mcp-aha get-key-result KR-789
```

### get-persona

Get a persona by ID.

```bash
mcp-aha get-persona PERS-123
```

### get-requirement

Get a requirement by ID.

```bash
mcp-aha get-requirement REQ-456
```

### get-team

Get a team by ID.

```bash
mcp-aha get-team TEAM-789
```

### get-user

Get a user by ID.

```bash
mcp-aha get-user USER-123
```

### get-workflow

Get a workflow by ID.

```bash
mcp-aha get-workflow WF-456
```

### get-comment

Get a comment by ID.

```bash
mcp-aha get-comment COMMENT-789
```

## Examples

### Scripting with JSON

```bash
# Get feature and extract name with jq
mcp-aha get-feature FEAT-123 | jq '.feature.name'

# List ideas and count them
mcp-aha list-ideas --query "mobile" | jq '.ideas | length'

# Search and format results
mcp-aha search-documents "roadmap" | jq -r '.results[].title'
```

### Using with Vault

```bash
# 1Password
export OP_SERVICE_ACCOUNT_TOKEN="ops_..."
mcp-aha get-feature FEAT-123 --vault op://MyVault --credentials-name aha

# Bitwarden
export BW_ACCESS_TOKEN="..."
export BW_ORGANIZATION_ID="..."
mcp-aha list-ideas --vault bw://org-id --credentials-name aha
```
