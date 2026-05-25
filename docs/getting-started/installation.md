# Installation

## Requirements

- Go 1.24 or later
- An Aha! workspace with API access
- An Aha! API token

## Install from Source

```bash
go install github.com/grokify/aha-mcp-server/cmd/aha-mcp-server@latest
```

## Build from Source

```bash
git clone https://github.com/grokify/aha-mcp-server.git
cd aha-mcp-server
go build ./cmd/aha-mcp-server
```

## Verify Installation

```bash
aha-mcp-server version
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
