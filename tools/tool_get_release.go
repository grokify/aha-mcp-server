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

type GetReleaseParams struct {
	ReleaseID string `json:"release_id" description:"Release ID to get"`
}

// GetReleaseResponse represents the structured response for getting a release
type GetReleaseResponse struct {
	Release    interface{} `json:"release"`
	StatusCode int         `json:"status_code"`
}

func (tc *ToolsClient) GetRelease(ctx context.Context, req *mcp.CallToolRequest, params GetReleaseParams) (*mcp.CallToolResult, any, error) {
	resp, err := tc.simpleClient.Do(ctx, httpsimple.Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("/api/v1/releases/%s", params.ReleaseID),
	})
	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("error getting Release: %v", err), true), nil, err
	}
	defer resp.Body.Close()

	releaseJSON, err := io.ReadAll(resp.Body)
	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error reading API response: %v", err), true), nil, err
	}

	var release interface{}
	if err := json.Unmarshal(releaseJSON, &release); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error unmarshaling API response: %v", err), true), nil, err
	}

	response := GetReleaseResponse{
		Release:    release,
		StatusCode: resp.StatusCode,
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error marshaling response: %v", err), true), nil, err
	}

	return mcputil.NewCallToolResultForAny(string(jsonData), false), string(jsonData), nil
}

func GetReleaseTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_release",
		Description: "Get Release from Aha",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"release_id": {
					Type:        "string",
					Description: "Release ID to get",
				},
			},
			Required: []string{"release_id"},
		},
	}
}
