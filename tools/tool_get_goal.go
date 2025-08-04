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

type GetGoalParams struct {
	GoalID string `json:"goal_id" description:"Goal ID to get"`
}

func (s *ToolsClient) GetGoal(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[GetGoalParams]) (*mcp.CallToolResultFor[any], error) {
	if resp, err := s.simpleClient.Do(ctx, httpsimple.Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("/api/v1/goals/%s", params.Arguments.GoalID),
	}); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("error getting Goal: %v", err), true), nil
	} else if goalJSON, err := io.ReadAll(resp.Body); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error unmarshaling API response: %v", err), true), nil
	} else if jsonData, err := json.MarshalIndent(map[string]any{
		"goal":        goalJSON,
		"status_code": resp.StatusCode,
	}, "", "  "); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error marshaling response: %v", err), true), nil
	} else {
		return mcputil.NewCallToolResultForAny(string(jsonData), false), nil
	}
}

func GetGoalTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_goal",
		Description: "Get Goal from Aha",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"goal_id": {
					Type:        "string",
					Description: "Goal ID to get",
				},
			},
			Required: []string{"goal_id"},
		},
	}
}
