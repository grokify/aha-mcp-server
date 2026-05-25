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

type GetCommentParams struct {
	CommentID string `json:"comment_id" description:"Comment ID to get"`
}

// GetCommentResponse represents the structured response for getting a comment
type GetCommentResponse struct {
	Comment    interface{} `json:"comment"`
	StatusCode int         `json:"status_code"`
}

func (tc *ToolsClient) GetComment(ctx context.Context, req *mcp.CallToolRequest, params GetCommentParams) (*mcp.CallToolResult, any, error) {
	resp, err := tc.simpleClient.Do(ctx, httpsimple.Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("/api/v1/comments/%s", params.CommentID),
	})
	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("error getting Comment: %v", err), true), nil, err
	}
	defer resp.Body.Close()

	commentJSON, err := io.ReadAll(resp.Body)
	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error reading API response: %v", err), true), nil, err
	}

	var comment interface{}
	if err := json.Unmarshal(commentJSON, &comment); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error unmarshaling API response: %v", err), true), nil, err
	}

	response := GetCommentResponse{
		Comment:    comment,
		StatusCode: resp.StatusCode,
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error marshaling response: %v", err), true), nil, err
	}

	return mcputil.NewCallToolResultForAny(string(jsonData), false), string(jsonData), nil
}

func GetCommentTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_comment",
		Description: "Get Comment from Aha",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"comment_id": {
					Type:        "string",
					Description: "Comment ID to get",
				},
			},
			Required: []string{"comment_id"},
		},
	}
}
