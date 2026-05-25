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

type GetWorkflowParams struct {
	WorkflowID string `json:"workflow_id" description:"Workflow ID to get"`
}

// GetWorkflowResponse represents the structured response for getting a workflow
type GetWorkflowResponse struct {
	Workflow   interface{} `json:"workflow"`
	StatusCode int         `json:"status_code"`
}

func (tc *ToolsClient) GetWorkflow(ctx context.Context, req *mcp.CallToolRequest, params GetWorkflowParams) (*mcp.CallToolResult, any, error) {
	resp, err := tc.simpleClient.Do(ctx, httpsimple.Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("/api/v1/workflows/%s", params.WorkflowID),
	})
	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("error getting Workflow: %v", err), true), nil, err
	}
	defer resp.Body.Close()

	workflowJSON, err := io.ReadAll(resp.Body)
	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error reading API response: %v", err), true), nil, err
	}

	var workflow interface{}
	if err := json.Unmarshal(workflowJSON, &workflow); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error unmarshaling API response: %v", err), true), nil, err
	}

	response := GetWorkflowResponse{
		Workflow:   workflow,
		StatusCode: resp.StatusCode,
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error marshaling response: %v", err), true), nil, err
	}

	return mcputil.NewCallToolResultForAny(string(jsonData), false), string(jsonData), nil
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
