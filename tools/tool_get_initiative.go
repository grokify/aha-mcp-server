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

type GetInitiativeParams struct {
	InitiativeID string `json:"initiative_id" description:"Initiative ID to get"`
}

func (tc *ToolsClient) GetInitiative(ctx context.Context, req *mcp.CallToolRequest, params GetInitiativeParams) (*mcp.CallToolResult, any, error) {
	if resp, err := tc.simpleClient.Do(ctx, httpsimple.Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("/api/v1/initiatives/%s", params.InitiativeID),
	}); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("error getting Initiative: %v", err), true), nil, err
	} else if initiativeJSON, err := io.ReadAll(resp.Body); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error unmarshaling API response: %v", err), true), nil, err
	} else if jsonData, err := json.MarshalIndent(map[string]any{
		"initiative":  initiativeJSON,
		"status_code": resp.StatusCode,
	}, "", "  "); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error marshaling response: %v", err), true), nil, err
	} else {
		return mcputil.NewCallToolResultForAny(string(jsonData), false), string(jsonData), nil
	}
}

func GetInitiativeTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_initiative",
		Description: "Get Initiative from Aha",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"initiative_id": {
					Type:        "string",
					Description: "Initiative ID to get",
				},
			},
			Required: []string{"initiative_id"},
		},
	}
}
