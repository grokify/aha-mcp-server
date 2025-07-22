package mcputil

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func NewCallToolResultForAny(msg string, isError bool) *mcp.CallToolResultFor[any] {
	return &mcp.CallToolResultFor[any]{
		IsError: isError,
		Content: []mcp.Content{&mcp.TextContent{Text: msg}},
	}
}
