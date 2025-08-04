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

type GetWorkflowParams struct {
	WorkflowID string `json:"workflow_id" description:"Workflow ID to get"`
}

func (s *ToolsClient) GetWorkflow(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[GetWorkflowParams]) (*mcp.CallToolResultFor[any], error) {
	if resp, err := s.simpleClient.Do(ctx, httpsimple.Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("/api/v1/workflows/%s", params.Arguments.WorkflowID),
	}); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("error getting Workflow: %v", err), true), nil
	} else if workflowJSON, err := io.ReadAll(resp.Body); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error unmarshaling API response: %v", err), true), nil
	} else if jsonData, err := json.MarshalIndent(map[string]any{
		"workflow":    workflowJSON,
		"status_code": resp.StatusCode,
	}, "", "  "); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error marshaling response: %v", err), true), nil
	} else {
		return mcputil.NewCallToolResultForAny(string(jsonData), false), nil
	}
}

func GetWorkflowTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_workflow",
		Description: "Get Workflow from Aha",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"workflow_id": {
					Type:        "string",
					Description: "Workflow ID to get",
				},
			},
			Required: []string{"workflow_id"},
		},
	}
}
