package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/grokify/mogo/net/http/httpsimple"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/grokify/aha-mcp-server/mcputil"
)

type GetReleaseParams struct {
	ReleaseID string `json:"release_id" description:"Release ID to get"`
}

func (tc *ToolsClient) GetRelease(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[GetReleaseParams]) (*mcp.CallToolResultFor[any], error) {
	if resp, err := tc.simpleClient.Do(ctx, httpsimple.Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("/api/v1/releases/%s", params.Arguments.ReleaseID),
	}); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("error getting Release: %v", err), true), nil
	} else if releaseJSON, err := io.ReadAll(resp.Body); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error unmarshaling API response: %v", err), true), nil
	} else if jsonData, err := json.MarshalIndent(map[string]any{
		"release":     releaseJSON,
		"status_code": resp.StatusCode,
	}, "", "  "); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error marshaling response: %v", err), true), nil
	} else {
		return mcputil.NewCallToolResultForAny(string(jsonData), false), nil
	}
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
