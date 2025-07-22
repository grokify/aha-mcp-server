# Aha! MCP Server

[![Build Status][build-status-svg]][build-status-url]
[![Lint Status][lint-status-svg]][lint-status-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![License][license-svg]][license-url]

A Model Context Protocol server for [Aha!](https://www.aha.io/).

## Installation

```
% go install github.com/grokify/aha-mcp-server/cmd/aha-mcp-server@v0.1.0
```

## Configuration

Configure with the following:

```json
{
	"mcpServers": {
		"aha": {
			"command": "aha-mcp-server",
			"env": {
				"AHA_API_TOKEN": "<your_aha_token>",
				"AHA_DOMAIN": "<your_aha_subdomain>"
			}
		}
	}
}
```

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