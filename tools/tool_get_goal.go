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

type GetGoalParams struct {
	GoalID string `json:"goal_id" description:"Goal ID to get"`
}

// GetGoalResponse represents the structured response for getting a goal
type GetGoalResponse struct {
	Goal       interface{} `json:"goal"`
	StatusCode int         `json:"status_code"`
}

func (tc *ToolsClient) GetGoal(ctx context.Context, req *mcp.CallToolRequest, params GetGoalParams) (*mcp.CallToolResult, any, error) {
	resp, err := tc.simpleClient.Do(ctx, httpsimple.Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("/api/v1/goals/%s", params.GoalID),
	})
	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("error getting Goal: %v", err), true), nil, err
	}
	defer resp.Body.Close()

	goalJSON, err := io.ReadAll(resp.Body)
	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error reading API response: %v", err), true), nil, err
	}

	var goal interface{}
	if err := json.Unmarshal(goalJSON, &goal); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error unmarshaling API response: %v", err), true), nil, err
	}

	response := GetGoalResponse{
		Goal:       goal,
		StatusCode: resp.StatusCode,
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error marshaling response: %v", err), true), nil, err
	}

	return mcputil.NewCallToolResultForAny(string(jsonData), false), string(jsonData), nil
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
