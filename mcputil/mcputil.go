package mcputil

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func NewCallToolResultForAny(msg string, isError bool) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		IsError: isError,
		Content: []mcp.Content{&mcp.TextContent{Text: msg}},
	}
}
