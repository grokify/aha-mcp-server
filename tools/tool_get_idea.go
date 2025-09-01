package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/grokify/aha-mcp-server/mcputil"
)

type GetIdeaParams struct {
	IdeaID string `json:"idea_id" description:"Idea ID to get"`
}

func (tc *ToolsClient) GetIdea(ctx context.Context, req *mcp.CallToolRequest, params GetIdeaParams) (*mcp.CallToolResult, any, error) {
	idea, resp, err := tc.client.IdeasAPI.GetIdeaExecute(
		tc.client.IdeasAPI.GetIdea(ctx, params.IdeaID))
	if err != nil {
		result := mcputil.NewCallToolResultForAny(fmt.Sprintf("Error getting idea: %v", err), true)
		return result, nil, nil
	}

	if jsonData, err := json.MarshalIndent(map[string]any{
		"idea":        idea,
		"status_code": resp.StatusCode,
	}, "", "  "); err != nil {
		result := mcputil.NewCallToolResultForAny(fmt.Sprintf("Error marshaling response: %v", err), true)
		return result, nil, nil
	} else {
		result := mcputil.NewCallToolResultForAny(string(jsonData), false)
		return result, string(jsonData), nil
	}
}

func GetIdeaTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_idea",
		Description: "Get Idea from Aha",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"idea_id": {
					Type:        "string",
					Description: "Idea ID to get",
				},
			},
			Required: []string{"idea_id"},
		},
	}
}
