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

type GetTeamParams struct {
	TeamID string `json:"team_id" description:"Team ID to get"`
}

// GetTeamResponse represents the structured response for getting a team
type GetTeamResponse struct {
	Team       interface{} `json:"team"`
	StatusCode int         `json:"status_code"`
}

func (tc *ToolsClient) GetTeam(ctx context.Context, req *mcp.CallToolRequest, params GetTeamParams) (*mcp.CallToolResult, any, error) {
	resp, err := tc.simpleClient.Do(ctx, httpsimple.Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("/api/v1/teams/%s", params.TeamID),
	})
	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("error getting Team: %v", err), true), nil, err
	}
	defer resp.Body.Close()

	teamJSON, err := io.ReadAll(resp.Body)
	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error reading API response: %v", err), true), nil, err
	}

	var team interface{}
	if err := json.Unmarshal(teamJSON, &team); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error unmarshaling API response: %v", err), true), nil, err
	}

	response := GetTeamResponse{
		Team:       team,
		StatusCode: resp.StatusCode,
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error marshaling response: %v", err), true), nil, err
	}

	return mcputil.NewCallToolResultForAny(string(jsonData), false), string(jsonData), nil
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
