package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/grokify/aha-mcp-server/mcputil"
)

type ListIdeasParams struct {
	Q              string `json:"q,omitempty" description:"Search term to match against the idea name"`
	Spam           *bool  `json:"spam,omitempty" description:"When true, shows ideas marked as spam"`
	WorkflowStatus string `json:"workflow_status,omitempty" description:"Filter by workflow status ID or name"`
	Sort           string `json:"sort,omitempty" description:"Sort by: recent, trending, or popular"`
	CreatedBefore  string `json:"created_before,omitempty" description:"UTC timestamp (ISO8601). Only ideas created before this time"`
	CreatedSince   string `json:"created_since,omitempty" description:"UTC timestamp (ISO8601). Only ideas created after this time"`
	UpdatedSince   string `json:"updated_since,omitempty" description:"UTC timestamp (ISO8601). Only ideas updated after this time"`
	Tag            string `json:"tag,omitempty" description:"Filter by tag value"`
	UserID         string `json:"user_id,omitempty" description:"Filter by creator user ID"`
	IdeaUserID     string `json:"idea_user_id,omitempty" description:"Filter by idea user ID"`
	Page           *int32 `json:"page,omitempty" description:"Page number"`
	PerPage        *int32 `json:"per_page,omitempty" description:"Results per page"`
}

func (tc *ToolsClient) ListIdeas(ctx context.Context, req *mcp.CallToolRequest, params ListIdeasParams) (*mcp.CallToolResult, any, error) {
	apiReq := tc.client.IdeasAPI.ListIdeas(ctx)

	if params.Q != "" {
		apiReq = apiReq.Q(params.Q)
	}
	if params.Spam != nil {
		apiReq = apiReq.Spam(*params.Spam)
	}
	if params.WorkflowStatus != "" {
		apiReq = apiReq.WorkflowStatus(params.WorkflowStatus)
	}
	if params.Sort != "" {
		apiReq = apiReq.Sort(params.Sort)
	}
	if params.CreatedBefore != "" {
		if t, err := time.Parse(time.RFC3339, params.CreatedBefore); err == nil {
			apiReq = apiReq.CreatedBefore(t)
		}
	}
	if params.CreatedSince != "" {
		if t, err := time.Parse(time.RFC3339, params.CreatedSince); err == nil {
			apiReq = apiReq.CreatedSince(t)
		}
	}
	if params.UpdatedSince != "" {
		if t, err := time.Parse(time.RFC3339, params.UpdatedSince); err == nil {
			apiReq = apiReq.UpdatedSince(t)
		}
	}
	if params.Tag != "" {
		apiReq = apiReq.Tag(params.Tag)
	}
	if params.UserID != "" {
		apiReq = apiReq.UserId(params.UserID)
	}
	if params.IdeaUserID != "" {
		apiReq = apiReq.IdeaUserId(params.IdeaUserID)
	}
	if params.Page != nil {
		apiReq = apiReq.Page(*params.Page)
	}
	if params.PerPage != nil {
		apiReq = apiReq.PerPage(*params.PerPage)
	}

	ideas, resp, err := apiReq.Execute()
	if err != nil {
		result := mcputil.NewCallToolResultForAny(fmt.Sprintf("Error listing ideas: %v", err), true)
		return result, nil, nil
	}

	if jsonData, err := json.MarshalIndent(map[string]any{
		"ideas":       ideas,
		"status_code": resp.StatusCode,
	}, "", "  "); err != nil {
		result := mcputil.NewCallToolResultForAny(fmt.Sprintf("Error marshaling response: %v", err), true)
		return result, nil, nil
	} else {
		result := mcputil.NewCallToolResultForAny(string(jsonData), false)
		return result, string(jsonData), nil
	}
}

func ListIdeasTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "list_ideas",
		Description: "List ideas from Aha with optional filtering and pagination",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"q": {
					Type:        "string",
					Description: "Search term to match against the idea name",
				},
				"spam": {
					Type:        "boolean",
					Description: "When true, shows ideas marked as spam",
				},
				"workflow_status": {
					Type:        "string",
					Description: "Filter by workflow status ID or name",
				},
				"sort": {
					Type:        "string",
					Description: "Sort by: recent, trending, or popular",
					Enum:        []any{"recent", "trending", "popular"},
				},
				"created_before": {
					Type:        "string",
					Description: "UTC timestamp (ISO8601). Only ideas created before this time",
				},
				"created_since": {
					Type:        "string",
					Description: "UTC timestamp (ISO8601). Only ideas created after this time",
				},
				"updated_since": {
					Type:        "string",
					Description: "UTC timestamp (ISO8601). Only ideas updated after this time",
				},
				"tag": {
					Type:        "string",
					Description: "Filter by tag value",
				},
				"user_id": {
					Type:        "string",
					Description: "Filter by creator user ID",
				},
				"idea_user_id": {
					Type:        "string",
					Description: "Filter by idea user ID",
				},
				"page": {
					Type:        "integer",
					Description: "Page number",
				},
				"per_page": {
					Type:        "integer",
					Description: "Results per page",
				},
			},
		},
	}
}
