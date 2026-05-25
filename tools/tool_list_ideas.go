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

// ListIdeasResponse represents the structured response for listing ideas
type ListIdeasResponse struct {
	Ideas      interface{} `json:"ideas"`
	StatusCode int         `json:"status_code"`
}

// parseTimestamp parses an RFC3339 timestamp string and returns a time.Time or an error
func parseTimestamp(value, fieldName string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid %s timestamp (expected RFC3339 format): %w", fieldName, err)
	}
	return t, nil
}

func (tc *ToolsClient) ListIdeas(ctx context.Context, req *mcp.CallToolRequest, params ListIdeasParams) (*mcp.CallToolResult, any, error) {
	apiReq := tc.client.IdeasAPI.ListIdeas(ctx)

	// Apply string parameters
	if params.Q != "" {
		apiReq = apiReq.Q(params.Q)
	}
	if params.WorkflowStatus != "" {
		apiReq = apiReq.WorkflowStatus(params.WorkflowStatus)
	}
	if params.Sort != "" {
		apiReq = apiReq.Sort(params.Sort)
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

	// Apply and validate timestamp parameters
	if params.CreatedBefore != "" {
		t, err := parseTimestamp(params.CreatedBefore, "created_before")
		if err != nil {
			return mcputil.NewCallToolResultForAny(err.Error(), true), nil, nil
		}
		apiReq = apiReq.CreatedBefore(t)
	}
	if params.CreatedSince != "" {
		t, err := parseTimestamp(params.CreatedSince, "created_since")
		if err != nil {
			return mcputil.NewCallToolResultForAny(err.Error(), true), nil, nil
		}
		apiReq = apiReq.CreatedSince(t)
	}
	if params.UpdatedSince != "" {
		t, err := parseTimestamp(params.UpdatedSince, "updated_since")
		if err != nil {
			return mcputil.NewCallToolResultForAny(err.Error(), true), nil, nil
		}
		apiReq = apiReq.UpdatedSince(t)
	}

	// Apply pointer parameters
	if params.Spam != nil {
		apiReq = apiReq.Spam(*params.Spam)
	}
	if params.Page != nil {
		apiReq = apiReq.Page(*params.Page)
	}
	if params.PerPage != nil {
		apiReq = apiReq.PerPage(*params.PerPage)
	}

	// Execute API request
	ideas, resp, err := apiReq.Execute()
	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error listing ideas: %v", err), true), nil, nil
	}

	// Create typed response
	response := ListIdeasResponse{
		Ideas:      ideas,
		StatusCode: resp.StatusCode,
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error marshaling response: %v", err), true), nil, nil
	}

	return mcputil.NewCallToolResultForAny(string(jsonData), false), string(jsonData), nil
}

func ListIdeasTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "list_ideas",
		Description: "List ideas from Aha with optional filtering and pagination",
		InputSchema: buildListIdeasSchema(),
	}
}

// buildListIdeasSchema constructs the JSON schema for the list_ideas tool.
// This helper improves readability by separating schema construction from tool definition.
func buildListIdeasSchema() *jsonschema.Schema {
	const (
		typeString  = "string"
		typeBoolean = "boolean"
		typeInteger = "integer"
	)

	return &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			// Search and filtering
			"q": {
				Type:        typeString,
				Description: "Search term to match against the idea name",
			},
			"spam": {
				Type:        typeBoolean,
				Description: "When true, shows ideas marked as spam",
			},
			"workflow_status": {
				Type:        typeString,
				Description: "Filter by workflow status ID or name",
			},
			"sort": {
				Type:        typeString,
				Description: "Sort by: recent, trending, or popular",
				Enum:        []any{"recent", "trending", "popular"},
			},

			// Timestamp filters
			"created_before": {
				Type:        typeString,
				Description: "UTC timestamp (ISO8601). Only ideas created before this time",
			},
			"created_since": {
				Type:        typeString,
				Description: "UTC timestamp (ISO8601). Only ideas created after this time",
			},
			"updated_since": {
				Type:        typeString,
				Description: "UTC timestamp (ISO8601). Only ideas updated after this time",
			},

			// Additional filters
			"tag": {
				Type:        typeString,
				Description: "Filter by tag value",
			},
			"user_id": {
				Type:        typeString,
				Description: "Filter by creator user ID",
			},
			"idea_user_id": {
				Type:        typeString,
				Description: "Filter by idea user ID",
			},

			// Pagination
			"page": {
				Type:        typeInteger,
				Description: "Page number",
			},
			"per_page": {
				Type:        typeInteger,
				Description: "Results per page",
			},
		},
	}
}
