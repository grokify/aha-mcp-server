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

type GetTeamParams struct {
	TeamID string `json:"team_id" description:"Team ID to get"`
}

func (tc *ToolsClient) GetTeam(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[GetTeamParams]) (*mcp.CallToolResultFor[any], error) {
	if resp, err := tc.simpleClient.Do(ctx, httpsimple.Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("/api/v1/teams/%s", params.Arguments.TeamID),
	}); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("error getting Team: %v", err), true), nil
	} else if teamJSON, err := io.ReadAll(resp.Body); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error unmarshaling API response: %v", err), true), nil
	} else if jsonData, err := json.MarshalIndent(map[string]any{
		"team":        teamJSON,
		"status_code": resp.StatusCode,
	}, "", "  "); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error marshaling response: %v", err), true), nil
	} else {
		return mcputil.NewCallToolResultForAny(string(jsonData), false), nil
	}
}

func GetTeamTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_team",
		Description: "Get Team from Aha",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"team_id": {
					Type:        "string",
					Description: "Team ID to get",
				},
			},
			Required: []string{"team_id"},
		},
	}
}
