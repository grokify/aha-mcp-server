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

type GetFeatureParams struct {
	FeatureID string `json:"feature_id" description:"Feature ID to get"`
}

func (tc *ToolsClient) GetFeature(ctx context.Context, req *mcp.CallToolRequest, params GetFeatureParams) (*mcp.CallToolResult, any, error) {
	if resp, err := tc.simpleClient.Do(ctx, httpsimple.Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("/api/v1/features/%s", params.FeatureID),
	}); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("error getting Feature: %v", err), true), nil, err
	} else if featureJSON, err := io.ReadAll(resp.Body); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error unmarshaling API response: %v", err), true), nil, err
	} else if jsonData, err := json.MarshalIndent(map[string]any{
		"feature":     featureJSON,
		"status_code": resp.StatusCode,
	}, "", "  "); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error marshaling response: %v", err), true), nil, err
	} else {
		return mcputil.NewCallToolResultForAny(string(jsonData), false), string(jsonData), nil
	}
}

func GetFeatureTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_feature",
		Description: "Get Feature from Aha",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"feature_id": {
					Type:        "string",
					Description: "Feature ID to get",
				},
			},
			Required: []string{"feature_id"},
		},
	}
}
