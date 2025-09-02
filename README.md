# Aha! MCP Server

[![Build Status][build-status-svg]][build-status-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![License][license-svg]][license-url]

A comprehensive Model Context Protocol (MCP) server for [Aha!](https://www.aha.io/) that enables AI assistants to interact with your Aha! workspace data. This server provides 13 tools to retrieve and search various Aha! objects, making it easy to integrate Aha! data into AI workflows.

## What is MCP?

The [Model Context Protocol](https://modelcontextprotocol.io/) is an open standard that enables AI assistants to securely connect to external data sources and tools. This Aha! MCP server acts as a bridge between AI assistants (like Claude) and your Aha! workspace.

## Features

- **13 comprehensive tools** for accessing and searching Aha! objects
- **Secure authentication** using Aha! API tokens
- **Easy configuration** with environment variables
- **Multiple deployment options** (stdio or HTTP)
- **Built with Go** for performance and reliability
- **MIT licensed** and open source

## Available Tools

This server provides the following tools to retrieve and search Aha! data:

| Category | Tool | Description |
|----------|------|-------------|
| **Search** | `search_documents` | Search for documents across your Aha! workspace using GraphQL |
| **Comments** | `get_comment` | Retrieve a specific comment by ID |
| **Epics** | `get_epic` | Retrieve a specific epic by ID |
| **Features** | `get_feature` | Retrieve a specific feature by ID |
| **Goals** | `get_goal` | Retrieve a specific goal by ID |
| **Ideas** | `get_idea` | Retrieve a specific idea by ID |
| **Initiatives** | `get_initiative` | Retrieve a specific initiative by ID |
| **Key Results** | `get_key_result` | Retrieve a specific key result by ID |
| **Personas** | `get_persona` | Retrieve a specific persona by ID |
| **Releases** | `get_release` | Retrieve a specific release by ID |
| **Requirements** | `get_requirement` | Retrieve a specific requirement by ID |
| **Teams** | `get_team` | Retrieve a specific team by ID |
| **Users** | `get_user` | Retrieve a specific user by ID |
| **Workflows** | `get_workflow` | Retrieve a specific workflow by ID |

All tools return JSON data including the requested object and HTTP status code.

## Prerequisites

- Go 1.24.1 or later
- An Aha! workspace with API access
- An Aha! API token (see [Aha! API documentation](https://www.aha.io/api))

## Installation

### Install from Source

```bash
go install github.com/grokify/aha-mcp-server/cmd/aha-mcp-server@v0.4.2
```

### Build from Source

```bash
git clone https://github.com/grokify/aha-mcp-server.git
cd aha-mcp-server
go build ./cmd/aha-mcp-server
```

## Configuration

### Get Your Aha! Credentials

1. **API Token**: Generate an API token from your Aha! account settings
2. **Domain**: Your Aha! subdomain (e.g., if your workspace is at `mycompany.aha.io`, your domain is `mycompany`)

### Environment Variables

Set the following environment variables:

```bash
export AHA_API_TOKEN="your_api_token_here"
export AHA_DOMAIN="your_aha_subdomain"
```

### Claude Desktop Configuration

Add this to your Claude Desktop configuration file:

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`  
**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "aha": {
      "command": "aha-mcp-server",
      "env": {
        "AHA_API_TOKEN": "your_api_token_here",
        "AHA_DOMAIN": "your_aha_subdomain"
      }
    }
  }
}
```

### Other MCP Clients

For other MCP clients, configure them to run the `aha-mcp-server` command with the required environment variables.

## Usage

### Basic Usage

Once configured, you can use natural language with your AI assistant to interact with Aha! data:

- "Search for documents about product roadmap"
- "Find all pages related to user authentication"
- "Show me feature AHA-123"
- "Get details for epic EPIC-456" 
- "What's in release REL-789?"
- "Tell me about user john.doe"

### Tool Parameters

**Search Tool:**
- `search_documents` requires:
  - `query` (required): Search query string
  - `searchable_type` (optional): Document type to search (defaults to "Page")

**Get Tools:**
Each get tool requires a specific ID parameter:
- `get_feature` requires `feature_id`
- `get_epic` requires `epic_id`
- `get_release` requires `release_id`
- And so on...

Example tool calls:
```json
{
  "tool": "search_documents",
  "parameters": {
    "query": "product roadmap",
    "searchable_type": "Page"
  }
}
```

```json
{
  "tool": "get_feature",
  "parameters": {
    "feature_id": "AHA-123"
  }
}
```

## Advanced Configuration

### HTTP Mode

You can run the server in HTTP mode for debugging or integration with other tools:

```bash
aha-mcp-server --http :8080
```

This will start an HTTP server on port 8080 instead of using stdio.

### Command Line Options

```bash
aha-mcp-server [OPTIONS]

Options:
  -h, --http string    HTTP address (e.g., :8080) - if set, uses HTTP instead of stdio
```

## Troubleshooting

### Common Issues

1. **"AHA_DOMAIN environment variable is required"**
   - Make sure you've set the `AHA_DOMAIN` environment variable
   - Verify it contains only your subdomain (e.g., `mycompany`, not `mycompany.aha.io`)

2. **"AHA_API_TOKEN environment variable is required"**
   - Ensure you've set a valid Aha! API token
   - Check that the token has the necessary permissions

3. **Connection errors**
   - Verify your Aha! subdomain is correct
   - Check that your API token is valid and not expired
   - Ensure your network allows connections to `*.aha.io`

### Debug Mode

Run with debug logging by setting the environment variable:

```bash
export MCP_DEBUG=1
```

## Development

### Project Structure

```
aha-mcp-server/
├── cmd/aha-mcp-server/     # Main application entry point
├── tools/                  # Tool implementations
├── mcputil/               # MCP utility functions
├── codegen/               # Code generation templates
├── server.go              # Core server implementation
└── go.mod                 # Go module definition
```

### Building

```bash
go build ./cmd/aha-mcp-server
```

### Testing

```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## Updating

### Version Number

When updating the version, update it in both [`README.md`](README.md) and [`server.go`](server.go).

## Comparison with Other Aha! MCP Servers

| Server | Tools | License | Language |
|--------|-------|---------|-----------|
| **This Server** | 13 | MIT | Go |
| [Official Aha! MCP](https://support.aha.io/aha-develop/integrations/mcp-server/mcp-server-connection~7493691606168806509) | 3 | ISC | TypeScript |
| [popand/aha-mcp](https://github.com/popand/aha-mcp) | 4 | ISC | TypeScript |
| [Zapier MCP](https://zapier.com/mcp/aha) | 2 | SaaS | - |

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.


 [build-status-svg]: https://github.com/grokify/aha-mcp-server/actions/workflows/ci.yaml/badge.svg?branch=main
 [build-status-url]: https://github.com/grokify/aha-mcp-server/actions/workflows/ci.yaml
 [lint-status-svg]: https://github.com/grokify/aha-mcp-server/actions/workflows/lint.yaml/badge.svg?branch=main
 [lint-status-url]: https://github.com/grokify/aha-mcp-server/actions/workflows/lint.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/aha-mcp-server
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/aha-mcp-server
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/aha-mcp-server
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/aha-mcp-server
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/aha-mcp-server/blob/main/LICENSE
 [used-by-svg]: https://sourcegraph.com/github.com/grokify/aha-mcp-server/-/badge.svg
 [used-by-url]: https://sourcegraph.com/github.com/grokify/aha-mcp-server?badge
 [loc-svg]: https://tokei.rs/b1/github/grokify/aha-mcp-server
 [repo-url]: https://github.com/grokify/aha-mcp-server