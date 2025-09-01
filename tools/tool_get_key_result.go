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

type GetKeyResultParams struct {
	KeyResultID string `json:"key_result_id" description:"Key Result ID to get"`
}

func (tc *ToolsClient) GetKeyResult(ctx context.Context, req *mcp.CallToolRequest, params GetKeyResultParams) (*mcp.CallToolResult, any, error) {
	if resp, err := tc.simpleClient.Do(ctx, httpsimple.Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("/api/v1/key_results/%s", params.KeyResultID),
	}); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("error getting Key Result: %v", err), true), nil, err
	} else if keyResultJSON, err := io.ReadAll(resp.Body); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error unmarshaling API response: %v", err), true), nil, err
	} else if jsonData, err := json.MarshalIndent(map[string]any{
		"key_result":  keyResultJSON,
		"status_code": resp.StatusCode,
	}, "", "  "); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error marshaling response: %v", err), true), nil, err
	} else {
		return mcputil.NewCallToolResultForAny(string(jsonData), false), string(jsonData), nil
	}
}

func GetKeyResultTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_key_result",
		Description: "Get Key Result from Aha",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"key_result_id": {
					Type:        "string",
					Description: "Key Result ID to get",
				},
			},
			Required: []string{"key_result_id"},
		},
	}
}
