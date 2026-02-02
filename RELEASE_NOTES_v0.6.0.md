# Release Notes - v0.6.0

**Release Date:** 2026-02-02

## Highlights

- New `list_ideas` tool for searching and filtering Aha! ideas

## What's New

### New Tool: `list_ideas`

The `list_ideas` tool enables searching and filtering Aha! ideas with comprehensive options:

| Parameter | Description |
|-----------|-------------|
| `q` | Search term to match against the idea name |
| `workflow_status` | Filter by workflow status ID or name |
| `tag` | Filter by tag value |
| `user_id` | Filter by creator user ID |
| `idea_user_id` | Filter by idea user ID |
| `created_before` | UTC timestamp (ISO8601) - ideas created before this time |
| `created_since` | UTC timestamp (ISO8601) - ideas created after this time |
| `updated_since` | UTC timestamp (ISO8601) - ideas updated after this time |
| `sort` | Sort by: `recent`, `trending`, or `popular` |
| `spam` | When true, shows ideas marked as spam |
| `page` | Page number for pagination |
| `per_page` | Results per page |

### Example Usage

```json
{
  "tool": "list_ideas",
  "arguments": {
    "q": "authentication",
    "workflow_status": "Under consideration",
    "sort": "popular",
    "per_page": 20
  }
}
```

## Dependencies

- Upgraded `github.com/grokify/go-aha/v3` to v3.3.0
  - Adds `ListIdeas` API support
  - Fixes `CustomField.Value` type to support arrays
- Bumped `github.com/modelcontextprotocol/go-sdk` from 1.1.0 to 1.2.0
- Bumped `github.com/google/jsonschema-go` from 0.3.0 to 0.4.1
- Bumped `github.com/grokify/mogo` from 0.71.10 to 0.72.3

## Installation

```bash
go install github.com/grokify/aha-mcp-server/cmd/aha-mcp-server@v0.6.0
```

## Full Changelog

See [CHANGELOG.md](CHANGELOG.md) for complete version history.
