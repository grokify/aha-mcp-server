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

// GetKeyResultResponse represents the structured response for getting a key result
type GetKeyResultResponse struct {
	KeyResult  interface{} `json:"key_result"`
	StatusCode int         `json:"status_code"`
}

func (tc *ToolsClient) GetKeyResult(ctx context.Context, req *mcp.CallToolRequest, params GetKeyResultParams) (*mcp.CallToolResult, any, error) {
	resp, err := tc.simpleClient.Do(ctx, httpsimple.Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("/api/v1/key_results/%s", params.KeyResultID),
	})
	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("error getting Key Result: %v", err), true), nil, err
	}
	defer resp.Body.Close()

	keyResultJSON, err := io.ReadAll(resp.Body)
	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error reading API response: %v", err), true), nil, err
	}

	var keyResult interface{}
	if err := json.Unmarshal(keyResultJSON, &keyResult); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error unmarshaling API response: %v", err), true), nil, err
	}

	response := GetKeyResultResponse{
		KeyResult:  keyResult,
		StatusCode: resp.StatusCode,
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error marshaling response: %v", err), true), nil, err
	}

	return mcputil.NewCallToolResultForAny(string(jsonData), false), string(jsonData), nil
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
