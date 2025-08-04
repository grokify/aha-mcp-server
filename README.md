# Aha! MCP Server

[![Build Status][build-status-svg]][build-status-url]
[![Lint Status][lint-status-svg]][lint-status-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![License][license-svg]][license-url]

A Model Context Protocol server for [Aha!](https://www.aha.io/).

## Tools

1. [ ] Comments
    - [x] `get_comment`: Get a specific comment
1. [ ] Epics
    - [x] `get_epic`: Get a specific epic
1. [ ] Features
    - [x] `get_feature`: Get a specific feature
1. [ ] Goals
    - [x] `get_goal`: Get a specific goal
1. [ ] Ideas
    - [x] `get_idea`: Get a specific idea
1. [ ] Initiatives
    - [x] `get_initiative`: Get a specific initiative
1. [ ] Key Results
    - [x] `get_key_result`: Get a specific key result
1. [ ] Personas
    - [x] `get_persona`: Get a specific persona
1. [ ] Releases
    - [x] `get_release`: Get a specific release
1. [ ] Requirements
    - [x] `get_requirement`: Get a specific requirement
1. [ ] Teams
    - [x] `get_team`: Get a specific team
1. [ ] Users
    - [x] `get_user`: Get a specific user
1. [ ] Workflows
    - [x] `get_workflow`: Get a specific workflow

## Installation

```
% go install github.com/grokify/aha-mcp-server/cmd/aha-mcp-server@v0.3.0
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

## Other Aha! MCP Servers

1. FOSS
    1. [Official Aha! MCP Server](https://support.aha.io/aha-develop/integrations/mcp-server/mcp-server-connection~7493691606168806509) (3 tools)
    1. [`github.com/popand/aha-mcp`](https://github.com/popand/aha-mcp) (4 tools)
1. SaaS
    1. [Zapier](https://zapier.com/mcp/aha) (2 tools)


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