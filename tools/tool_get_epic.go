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

type GetEpicParams struct {
	EpicID string `json:"epic_id" description:"Epic ID to get"`
}

// GetEpicResponse represents the structured response for getting an epic
type GetEpicResponse struct {
	Epic       interface{} `json:"epic"`
	StatusCode int         `json:"status_code"`
}

func (tc *ToolsClient) GetEpic(ctx context.Context, req *mcp.CallToolRequest, params GetEpicParams) (*mcp.CallToolResult, any, error) {
	resp, err := tc.simpleClient.Do(ctx, httpsimple.Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("/api/v1/epics/%s", params.EpicID),
	})
	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("error getting Epic: %v", err), true), nil, err
	}
	defer resp.Body.Close()

	epicJSON, err := io.ReadAll(resp.Body)
	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error reading API response: %v", err), true), nil, err
	}

	var epic interface{}
	if err := json.Unmarshal(epicJSON, &epic); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error unmarshaling API response: %v", err), true), nil, err
	}

	response := GetEpicResponse{
		Epic:       epic,
		StatusCode: resp.StatusCode,
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error marshaling response: %v", err), true), nil, err
	}

	return mcputil.NewCallToolResultForAny(string(jsonData), false), string(jsonData), nil
}

func GetEpicTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_epic",
		Description: "Get Epic from Aha",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"epic_id": {
					Type:        "string",
					Description: "Epic ID to get",
				},
			},
			Required: []string{"epic_id"},
		},
	}
}
