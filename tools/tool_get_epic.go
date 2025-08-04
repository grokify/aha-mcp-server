
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

type GetEpicParams struct {
	EpicID string `json:"epic_id" description:"Epic ID to get"`
}

func (s *ToolsClient) GetEpic(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[GetEpicParams]) (*mcp.CallToolResultFor[any], error) {
	if resp, err := s.simpleClient.Do(ctx, httpsimple.Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("/api/v1/epics/%s", params.Arguments.EpicID),
	}); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("error getting Epic: %v", err), true), nil
	} else if epicJSON, err := io.ReadAll(resp.Body); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error unmarshaling API response: %v", err), true), nil
	} else if jsonData, err := json.MarshalIndent(map[string]any{
		"epic": epicJSON,
		"status_code": resp.StatusCode,
	}, "", "  "); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error marshaling response: %v", err), true), nil
	} else {
		return mcputil.NewCallToolResultForAny(string(jsonData), false), nil
	}
}

func GetEpicTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_epic",
		Description: "Get Epic from Aha",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"epic_id": {
					Type:        "string",
					Description: "Epic ID to get",
				},
			},
			Required: []string{"epic_id"},
		},
	}
}
