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

type GetRequirementParams struct {
	RequirementID string `json:"requirement_id" description:"Requirement ID to get"`
}

// GetRequirementResponse represents the structured response for getting a requirement
type GetRequirementResponse struct {
	Requirement interface{} `json:"requirement"`
	StatusCode  int         `json:"status_code"`
}

func (tc *ToolsClient) GetRequirement(ctx context.Context, req *mcp.CallToolRequest, params GetRequirementParams) (*mcp.CallToolResult, any, error) {
	resp, err := tc.simpleClient.Do(ctx, httpsimple.Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("/api/v1/requirements/%s", params.RequirementID),
	})
	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("error getting Requirement: %v", err), true), nil, err
	}
	defer resp.Body.Close()

	requirementJSON, err := io.ReadAll(resp.Body)
	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error reading API response: %v", err), true), nil, err
	}

	var requirement interface{}
	if err := json.Unmarshal(requirementJSON, &requirement); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error unmarshaling API response: %v", err), true), nil, err
	}

	response := GetRequirementResponse{
		Requirement: requirement,
		StatusCode:  resp.StatusCode,
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error marshaling response: %v", err), true), nil, err
	}

	return mcputil.NewCallToolResultForAny(string(jsonData), false), string(jsonData), nil
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
