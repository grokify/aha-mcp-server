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

// GetIdeaResponse represents the structured response for getting an idea
type GetIdeaResponse struct {
	Idea       interface{} `json:"idea"`
	StatusCode int         `json:"status_code"`
}

func (tc *ToolsClient) GetIdea(ctx context.Context, req *mcp.CallToolRequest, params GetIdeaParams) (*mcp.CallToolResult, any, error) {
	idea, resp, err := tc.client.IdeasAPI.GetIdeaExecute(
		tc.client.IdeasAPI.GetIdea(ctx, params.IdeaID))
	if err != nil {
		result := mcputil.NewCallToolResultForAny(fmt.Sprintf("Error getting idea: %v", err), true)
		return result, nil, nil
	}

	response := GetIdeaResponse{
		Idea:       idea,
		StatusCode: resp.StatusCode,
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error marshaling response: %v", err), true), nil, nil
	}

	return mcputil.NewCallToolResultForAny(string(jsonData), false), string(jsonData), nil
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
