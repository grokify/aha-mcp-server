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

// GetFeatureResponse represents the structured response for getting a feature
type GetFeatureResponse struct {
	Feature    interface{} `json:"feature"`
	StatusCode int         `json:"status_code"`
}

func (tc *ToolsClient) GetFeature(ctx context.Context, req *mcp.CallToolRequest, params GetFeatureParams) (*mcp.CallToolResult, any, error) {
	resp, err := tc.simpleClient.Do(ctx, httpsimple.Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("/api/v1/features/%s", params.FeatureID),
	})
	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("error getting Feature: %v", err), true), nil, err
	}
	defer resp.Body.Close()

	featureJSON, err := io.ReadAll(resp.Body)
	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error reading API response: %v", err), true), nil, err
	}

	var feature interface{}
	if err := json.Unmarshal(featureJSON, &feature); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error unmarshaling API response: %v", err), true), nil, err
	}

	response := GetFeatureResponse{
		Feature:    feature,
		StatusCode: resp.StatusCode,
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error marshaling response: %v", err), true), nil, err
	}

	return mcputil.NewCallToolResultForAny(string(jsonData), false), string(jsonData), nil
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
