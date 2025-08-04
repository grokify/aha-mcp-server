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

type GetPersonaParams struct {
	PersonaID string `json:"persona_id" description:"Persona ID to get"`
}

func (tc *ToolsClient) GetPersona(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[GetPersonaParams]) (*mcp.CallToolResultFor[any], error) {
	if resp, err := tc.simpleClient.Do(ctx, httpsimple.Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("/api/v1/personas/%s", params.Arguments.PersonaID),
	}); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("error getting Persona: %v", err), true), nil
	} else if personaJSON, err := io.ReadAll(resp.Body); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error unmarshaling API response: %v", err), true), nil
	} else if jsonData, err := json.MarshalIndent(map[string]any{
		"persona":     personaJSON,
		"status_code": resp.StatusCode,
	}, "", "  "); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error marshaling response: %v", err), true), nil
	} else {
		return mcputil.NewCallToolResultForAny(string(jsonData), false), nil
	}
}

func GetPersonaTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_persona",
		Description: "Get Persona from Aha",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"persona_id": {
					Type:        "string",
					Description: "Persona ID to get",
				},
			},
			Required: []string{"persona_id"},
		},
	}
}
