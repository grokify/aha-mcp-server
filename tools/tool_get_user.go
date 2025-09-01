package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/grokify/mogo/net/http/httpsimple"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/grokify/aha-mcp-server/mcputil"
)

type GetUserParams struct {
	UserID string `json:"user_id" description:"User ID to get"`
}

func (tc *ToolsClient) GetUser(ctx context.Context, req *mcp.CallToolRequest, params GetUserParams) (*mcp.CallToolResult, any, error) {
	if resp, err := tc.simpleClient.Do(ctx, httpsimple.Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("/api/v1/users/%s", params.UserID),
	}); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("error getting User: %v", err), true), nil, err
	} else if userJSON, err := io.ReadAll(resp.Body); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error unmarshaling API response: %v", err), true), nil, err
	} else if jsonData, err := json.MarshalIndent(map[string]any{
		"user":        userJSON,
		"status_code": resp.StatusCode,
	}, "", "  "); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error marshaling response: %v", err), true), nil, err
	} else {
		return mcputil.NewCallToolResultForAny(string(jsonData), false), string(jsonData), nil
	}
}

func GetUserTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_user",
		Description: "Get User from Aha",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"user_id": {
					Type:        "string",
					Description: "User ID to get",
				},
			},
			Required: []string{"user_id"},
		},
	}
}
