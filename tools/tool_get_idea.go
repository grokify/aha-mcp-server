package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/grokify/aha-mcp-server/mcputil"
)

type GetIdeaParams struct {
	IdeaID string `json:"idea_id" description:"Idea ID to get"`
}

func (tc *ToolsClient) GetIdea(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[GetIdeaParams]) (*mcp.CallToolResultFor[any], error) {
	idea, resp, err := tc.client.IdeasAPI.GetIdeaExecute(
		tc.client.IdeasAPI.GetIdea(ctx, params.Arguments.IdeaID))

	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error getting idea: %v", err), true), nil
	}

	if jsonData, err := json.MarshalIndent(map[string]any{
		"idea":        idea,
		"status_code": resp.StatusCode,
	}, "", "  "); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error marshaling response: %v", err), true), nil
	} else {
		return mcputil.NewCallToolResultForAny(string(jsonData), false), nil
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
