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

type GetCommentParams struct {
	CommentID string `json:"comment_id" description:"Comment ID to get"`
}

func (tc *ToolsClient) GetComment(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[GetCommentParams]) (*mcp.CallToolResultFor[any], error) {
	if resp, err := tc.simpleClient.Do(ctx, httpsimple.Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("/api/v1/comments/%s", params.Arguments.CommentID),
	}); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("error getting Comment: %v", err), true), nil
	} else if commentJSON, err := io.ReadAll(resp.Body); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error unmarshaling API response: %v", err), true), nil
	} else if jsonData, err := json.MarshalIndent(map[string]any{
		"comment":     commentJSON,
		"status_code": resp.StatusCode,
	}, "", "  "); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error marshaling response: %v", err), true), nil
	} else {
		return mcputil.NewCallToolResultForAny(string(jsonData), false), nil
	}
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
