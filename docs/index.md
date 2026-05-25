# Aha! MCP Server

A comprehensive Model Context Protocol (MCP) server for [Aha!](https://www.aha.io/) that enables AI assistants to interact with your Aha! workspace data.

## Features

- **14 comprehensive tools** for accessing and searching Aha! objects
- **Secure authentication** using Aha! API tokens or vault-backed credentials
- **CLI mode** for testing and scripting
- **Easy configuration** with environment variables or command-line flags
- **Built with Go** for performance and reliability

## What is MCP?

The [Model Context Protocol](https://modelcontextprotocol.io/) is an open standard that enables AI assistants to securely connect to external data sources and tools. This Aha! MCP server acts as a bridge between AI assistants (like Claude) and your Aha! workspace.

## Available Tools

| Category | Tool | Description |
|----------|------|-------------|
| **Search** | `search_documents` | Search for documents across your Aha! workspace |
| **Ideas** | `get_idea` | Retrieve a specific idea by ID |
| **Ideas** | `list_ideas` | List ideas with filtering and pagination |
| **Features** | `get_feature` | Retrieve a specific feature by ID |
| **Epics** | `get_epic` | Retrieve a specific epic by ID |
| **Releases** | `get_release` | Retrieve a specific release by ID |
| **Goals** | `get_goal` | Retrieve a specific goal by ID |
| **Initiatives** | `get_initiative` | Retrieve a specific initiative by ID |
| **Key Results** | `get_key_result` | Retrieve a specific key result by ID |
| **Personas** | `get_persona` | Retrieve a specific persona by ID |
| **Requirements** | `get_requirement` | Retrieve a specific requirement by ID |
| **Teams** | `get_team` | Retrieve a specific team by ID |
| **Users** | `get_user` | Retrieve a specific user by ID |
| **Workflows** | `get_workflow` | Retrieve a specific workflow by ID |
| **Comments** | `get_comment` | Retrieve a specific comment by ID |

## Quick Start

```bash
# Install
go install github.com/grokify/aha-mcp-server/cmd/aha-mcp-server@latest

# Configure credentials
export AHA_DOMAIN="your_subdomain"
export AHA_API_TOKEN="your_api_token"

# Run as MCP server
aha-mcp-server

# Or use CLI mode
aha-mcp-server get-idea IDEA-123
```

## Next Steps

- [Installation](getting-started/installation.md) - Install the server
- [Setup](getting-started/setup.md) - Configure your credentials
- [Quick Start](getting-started/quickstart.md) - Start using the tools
- [Tools Reference](tools/overview.md) - Detailed tool documentation
