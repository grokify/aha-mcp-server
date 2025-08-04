
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

type GetRequirementParams struct {
	RequirementID string `json:"requirement_id" description:"Requirement ID to get"`
}

func (s *ToolsClient) GetRequirement(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[GetRequirementParams]) (*mcp.CallToolResultFor[any], error) {
	if resp, err := s.simpleClient.Do(ctx, httpsimple.Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("/api/v1/requirements/%s", params.Arguments.RequirementID),
	}); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("error getting Requirement: %v", err), true), nil
	} else if requirementJSON, err := io.ReadAll(resp.Body); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error unmarshaling API response: %v", err), true), nil
	} else if jsonData, err := json.MarshalIndent(map[string]any{
		"requirement": requirementJSON,
		"status_code": resp.StatusCode,
	}, "", "  "); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error marshaling response: %v", err), true), nil
	} else {
		return mcputil.NewCallToolResultForAny(string(jsonData), false), nil
	}
}

func GetRequirementTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_requirement",
		Description: "Get Requirement from Aha",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"requirement_id": {
					Type:        "string",
					Description: "Requirement ID to get",
				},
			},
			Required: []string{"requirement_id"},
		},
	}
}
