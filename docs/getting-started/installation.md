# Installation

## Requirements

- Go 1.24 or later
- An Aha! workspace with API access
- An Aha! API token

## Install from Source

```bash
go install github.com/grokify/aha-mcp-server/cmd/mcp-aha@latest
```

## Build from Source

```bash
git clone https://github.com/grokify/aha-mcp-server.git
cd aha-mcp-server
go build ./cmd/mcp-aha
```

## Verify Installation

```bash
mcp-aha version
```

## Finding the Binary Path

### If installed with `go install`

```bash
# Find GOPATH
go env GOPATH

# Binary is at $GOPATH/bin/mcp-aha
# Typically: ~/go/bin/mcp-aha
```

### If built from source

Use the full path to where you built it:

```bash
/path/to/mcp-aha/mcp-aha
```
