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

type GetInitiativeParams struct {
	InitiativeID string `json:"initiative_id" description:"Initiative ID to get"`
}

func (s *ToolsClient) GetInitiative(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[GetInitiativeParams]) (*mcp.CallToolResultFor[any], error) {
	if resp, err := s.simpleClient.Do(ctx, httpsimple.Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("/api/v1/initiatives/%s", params.Arguments.InitiativeID),
	}); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("error getting Initiative: %v", err), true), nil
	} else if initiativeJSON, err := io.ReadAll(resp.Body); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error unmarshaling API response: %v", err), true), nil
	} else if jsonData, err := json.MarshalIndent(map[string]any{
		"initiative":  initiativeJSON,
		"status_code": resp.StatusCode,
	}, "", "  "); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error marshaling response: %v", err), true), nil
	} else {
		return mcputil.NewCallToolResultForAny(string(jsonData), false), nil
	}
}

func GetInitiativeTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_initiative",
		Description: "Get Initiative from Aha",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"initiative_id": {
					Type:        "string",
					Description: "Initiative ID to get",
				},
			},
			Required: []string{"initiative_id"},
		},
	}
}
