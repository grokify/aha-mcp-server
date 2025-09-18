package tools

import (
	"bytes"
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

type SearchDocumentsParams struct {
	Query          string `json:"query" description:"Search query string"`
	SearchableType string `json:"searchable_type,omitempty" description:"Type of document to search for (defaults to Page)"`
}

type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type GraphQLResponse struct {
	Data   SearchData     `json:"data"`
	Errors []GraphQLError `json:"errors,omitempty"`
}

type GraphQLError struct {
	Message string `json:"message"`
}

type SearchData struct {
	Search SearchResults `json:"search"`
}

type SearchResults struct {
	Nodes       []DocumentNode `json:"nodes"`
	CurrentPage int            `json:"currentPage"`
	TotalCount  int            `json:"totalCount"`
	TotalPages  int            `json:"totalPages"`
	IsLastPage  bool           `json:"isLastPage"`
}

type DocumentNode struct {
	Name           string `json:"name"`
	URL            string `json:"url"`
	SearchableID   string `json:"searchableId"`
	SearchableType string `json:"searchableType"`
}

const searchDocumentsQuery = `
query SearchDocuments($query: String!, $searchableType: [String!]) {
  search(query: $query, searchableType: $searchableType, first: 50) {
    nodes {
      name
      url
      searchableId
      searchableType
    }
    currentPage
    totalCount
    totalPages
    isLastPage
  }
}
`

func (tc *ToolsClient) SearchDocuments(ctx context.Context, req *mcp.CallToolRequest, params SearchDocumentsParams) (*mcp.CallToolResult, any, error) {
	if params.Query == "" {
		return mcputil.NewCallToolResultForAny("Query parameter is required", true), nil, fmt.Errorf("query parameter is required")
	}

	// Default to "Page" if no searchable type is specified
	searchableType := params.SearchableType
	if searchableType == "" {
		searchableType = "Page"
	}

	// Prepare GraphQL request
	variables := map[string]interface{}{
		"query":          params.Query,
		"searchableType": []string{searchableType},
	}

	graphqlReq := GraphQLRequest{
		Query:     searchDocumentsQuery,
		Variables: variables,
	}

	// Marshal request body
	requestBody, err := json.Marshal(graphqlReq)
	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error marshaling GraphQL request: %v", err), true), nil, err
	}

	// Use the simpleClient to make the GraphQL request
	httpReq := httpsimple.Request{
		Method: http.MethodPost,
		URL:    "/api/graphql",
		Body:   bytes.NewBuffer(requestBody),
	}
	if httpReq.Headers == nil {
		httpReq.Headers = make(http.Header)
	}
	httpReq.Headers.Set("Content-Type", "application/json")

	resp, err := tc.simpleClient.Do(ctx, httpReq)
	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error making GraphQL request: %v", err), true), nil, err
	}
	defer resp.Body.Close()

	// Read response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error reading response: %v", err), true), nil, err
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("GraphQL request failed with status %d: %s", resp.StatusCode, string(responseBody)), true), nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(responseBody))
	}

	// Parse GraphQL response
	var graphqlResp GraphQLResponse
	if err := json.Unmarshal(responseBody, &graphqlResp); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error parsing GraphQL response: %v", err), true), nil, err
	}

	// Check for GraphQL errors
	if len(graphqlResp.Errors) > 0 {
		errorMessages := make([]string, len(graphqlResp.Errors))
		for i, err := range graphqlResp.Errors {
			errorMessages[i] = err.Message
		}
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("GraphQL errors: %v", errorMessages), true), nil, fmt.Errorf("GraphQL errors: %v", errorMessages)
	}

	// Format the response for MCP
	searchResults := map[string]interface{}{
		"results":       make([]map[string]interface{}, len(graphqlResp.Data.Search.Nodes)),
		"total_results": graphqlResp.Data.Search.TotalCount,
		"current_page":  graphqlResp.Data.Search.CurrentPage,
		"total_pages":   graphqlResp.Data.Search.TotalPages,
		"is_last_page":  graphqlResp.Data.Search.IsLastPage,
	}

	// Transform nodes to the expected format
	results := searchResults["results"].([]map[string]interface{})
	for i, node := range graphqlResp.Data.Search.Nodes {
		results[i] = map[string]interface{}{
			"reference_num": node.SearchableID,
			"name":          node.Name,
			"type":          node.SearchableType,
			"url":           node.URL,
		}
	}

	// Marshal the final response
	jsonData, err := json.MarshalIndent(searchResults, "", "  ")
	if err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error marshaling response: %v", err), true), nil, err
	}

	return mcputil.NewCallToolResultForAny(string(jsonData), false), string(jsonData), nil
}

func SearchDocumentsTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "search_documents",
		Description: "Search for Aha! documents using GraphQL",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"query": {
					Type:        "string",
					Description: "Search query string",
				},
				"searchable_type": {
					Type:        "string",
					Description: "Type of document to search for (defaults to Page). Examples: Page, Feature, Epic, Release, etc.",
				},
			},
			Required: []string{"query"},
		},
	}
}
