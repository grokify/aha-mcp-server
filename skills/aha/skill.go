// Package aha provides an omniskill Skill for reading Aha! product management data.
//
// This package can be used standalone with mcp-aha or composed
// with other skills in a multi-service MCP server.
package aha

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	goaha "github.com/grokify/go-aha/v3/oag7/aha"
	goaclient "github.com/grokify/go-aha/v3/oag7/client"
	"github.com/grokify/mogo/net/http/httpsimple"
	"github.com/plexusone/omniskill/skill"
)

// Skill provides Aha! product management tools.
type Skill struct {
	client       *goaha.APIClient
	config       *goaha.Configuration
	simpleClient *httpsimple.Client
}

// New creates a new Aha skill with the given subdomain and API key.
func New(subdomain, apiKey string) (*Skill, error) {
	config, err := goaclient.NewConfiguration(subdomain, apiKey)
	if err != nil {
		return nil, err
	}
	sc, err := goaclient.NewSimpleClient(subdomain, apiKey)
	if err != nil {
		return nil, err
	}
	return &Skill{
		client:       goaha.NewAPIClient(config),
		config:       config,
		simpleClient: sc,
	}, nil
}

// Name returns the skill identifier.
func (s *Skill) Name() string {
	return "aha"
}

// Description returns what this skill does.
func (s *Skill) Description() string {
	return "Read ideas, features, epics, releases, and other product management data from Aha!"
}

// Init initializes the skill (no-op as client is injected).
func (s *Skill) Init(ctx context.Context) error {
	return nil
}

// Close releases resources (no-op for this skill).
func (s *Skill) Close() error {
	return nil
}

// Tools returns all tools provided by this skill.
func (s *Skill) Tools() []skill.Tool {
	return []skill.Tool{
		s.getIdeaTool(),
		s.listIdeasTool(),
		s.getCommentTool(),
		s.getEpicTool(),
		s.getFeatureTool(),
		s.getGoalTool(),
		s.getInitiativeTool(),
		s.getKeyResultTool(),
		s.getPersonaTool(),
		s.getReleaseTool(),
		s.getRequirementTool(),
		s.searchDocumentsTool(),
		s.getTeamTool(),
		s.getUserTool(),
		s.getWorkflowTool(),
	}
}

// Ensure Skill implements skill.Skill.
var _ skill.Skill = (*Skill)(nil)

// resourceConfig defines parameters for a generic GET resource tool.
type resourceConfig struct {
	name     string // e.g., "comment"
	title    string // e.g., "Comment"
	endpoint string // e.g., "comments"
}

// getResourceTool creates a generic tool for fetching a single resource by ID.
func (s *Skill) getResourceTool(cfg resourceConfig) skill.Tool {
	paramName := cfg.name + "_id"
	return skill.NewTool(
		"get_"+cfg.name,
		"Get "+cfg.title+" from Aha",
		map[string]skill.Parameter{
			paramName: {
				Type:        "string",
				Description: cfg.title + " ID to get",
				Required:    true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			id, _ := params[paramName].(string)
			if id == "" {
				return nil, fmt.Errorf("%s is required", paramName)
			}

			resp, err := s.simpleClient.Do(ctx, httpsimple.Request{
				Method: http.MethodGet,
				URL:    fmt.Sprintf("/api/v1/%s/%s", cfg.endpoint, id),
			})
			if err != nil {
				return map[string]any{"error": fmt.Sprintf("error getting %s: %v", cfg.title, err)}, nil
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return map[string]any{"error": fmt.Sprintf("error reading response: %v", err)}, nil
			}

			var result any
			if err := json.Unmarshal(body, &result); err != nil {
				return map[string]any{
					cfg.name:      string(body),
					"status_code": resp.StatusCode,
				}, nil
			}

			return map[string]any{
				cfg.name:      result,
				"status_code": resp.StatusCode,
			}, nil
		},
	)
}

func (s *Skill) getIdeaTool() skill.Tool {
	return skill.NewTool(
		"get_idea",
		"Get Idea from Aha",
		map[string]skill.Parameter{
			"idea_id": {
				Type:        "string",
				Description: "Idea ID to get",
				Required:    true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			ideaID, _ := params["idea_id"].(string)
			if ideaID == "" {
				return nil, fmt.Errorf("idea_id is required")
			}

			idea, resp, err := s.client.IdeasAPI.GetIdeaExecute(
				s.client.IdeasAPI.GetIdea(ctx, ideaID))
			if err != nil {
				return map[string]any{
					"error": fmt.Sprintf("Error getting idea: %v", err),
				}, nil
			}

			return map[string]any{
				"idea":        idea,
				"status_code": resp.StatusCode,
			}, nil
		},
	)
}

func (s *Skill) listIdeasTool() skill.Tool {
	return skill.NewTool(
		"list_ideas",
		"List ideas from Aha with optional filtering and pagination",
		map[string]skill.Parameter{
			"q": {
				Type:        "string",
				Description: "Search term to match against the idea name",
				Required:    false,
			},
			"spam": {
				Type:        "boolean",
				Description: "When true, shows ideas marked as spam",
				Required:    false,
			},
			"workflow_status": {
				Type:        "string",
				Description: "Filter by workflow status ID or name",
				Required:    false,
			},
			"sort": {
				Type:        "string",
				Description: "Sort by: recent, trending, or popular",
				Required:    false,
			},
			"created_before": {
				Type:        "string",
				Description: "UTC timestamp (ISO8601). Only ideas created before this time",
				Required:    false,
			},
			"created_since": {
				Type:        "string",
				Description: "UTC timestamp (ISO8601). Only ideas created after this time",
				Required:    false,
			},
			"updated_since": {
				Type:        "string",
				Description: "UTC timestamp (ISO8601). Only ideas updated after this time",
				Required:    false,
			},
			"tag": {
				Type:        "string",
				Description: "Filter by tag value",
				Required:    false,
			},
			"user_id": {
				Type:        "string",
				Description: "Filter by creator user ID",
				Required:    false,
			},
			"idea_user_id": {
				Type:        "string",
				Description: "Filter by idea user ID",
				Required:    false,
			},
			"page": {
				Type:        "integer",
				Description: "Page number",
				Required:    false,
			},
			"per_page": {
				Type:        "integer",
				Description: "Results per page",
				Required:    false,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			apiReq := s.client.IdeasAPI.ListIdeas(ctx)

			if q, ok := params["q"].(string); ok && q != "" {
				apiReq = apiReq.Q(q)
			}
			if spam, ok := params["spam"].(bool); ok {
				apiReq = apiReq.Spam(spam)
			}
			if ws, ok := params["workflow_status"].(string); ok && ws != "" {
				apiReq = apiReq.WorkflowStatus(ws)
			}
			if sort, ok := params["sort"].(string); ok && sort != "" {
				apiReq = apiReq.Sort(sort)
			}
			if cb, ok := params["created_before"].(string); ok && cb != "" {
				if t, err := time.Parse(time.RFC3339, cb); err == nil {
					apiReq = apiReq.CreatedBefore(t)
				}
			}
			if cs, ok := params["created_since"].(string); ok && cs != "" {
				if t, err := time.Parse(time.RFC3339, cs); err == nil {
					apiReq = apiReq.CreatedSince(t)
				}
			}
			if us, ok := params["updated_since"].(string); ok && us != "" {
				if t, err := time.Parse(time.RFC3339, us); err == nil {
					apiReq = apiReq.UpdatedSince(t)
				}
			}
			if tag, ok := params["tag"].(string); ok && tag != "" {
				apiReq = apiReq.Tag(tag)
			}
			if uid, ok := params["user_id"].(string); ok && uid != "" {
				apiReq = apiReq.UserId(uid)
			}
			if iuid, ok := params["idea_user_id"].(string); ok && iuid != "" {
				apiReq = apiReq.IdeaUserId(iuid)
			}
			if page, ok := params["page"].(float64); ok {
				apiReq = apiReq.Page(int32(page))
			}
			if perPage, ok := params["per_page"].(float64); ok {
				apiReq = apiReq.PerPage(int32(perPage))
			}

			ideas, resp, err := apiReq.Execute()
			if err != nil {
				return map[string]any{
					"error": fmt.Sprintf("Error listing ideas: %v", err),
				}, nil
			}

			return map[string]any{
				"ideas":       ideas,
				"status_code": resp.StatusCode,
			}, nil
		},
	)
}

func (s *Skill) getCommentTool() skill.Tool {
	return s.getResourceTool(resourceConfig{name: "comment", title: "Comment", endpoint: "comments"})
}

func (s *Skill) getEpicTool() skill.Tool {
	return s.getResourceTool(resourceConfig{name: "epic", title: "Epic", endpoint: "epics"})
}

func (s *Skill) getFeatureTool() skill.Tool {
	return s.getResourceTool(resourceConfig{name: "feature", title: "Feature", endpoint: "features"})
}

func (s *Skill) getGoalTool() skill.Tool {
	return s.getResourceTool(resourceConfig{name: "goal", title: "Goal", endpoint: "goals"})
}

func (s *Skill) getInitiativeTool() skill.Tool {
	return s.getResourceTool(resourceConfig{name: "initiative", title: "Initiative", endpoint: "initiatives"})
}

func (s *Skill) getKeyResultTool() skill.Tool {
	return s.getResourceTool(resourceConfig{name: "key_result", title: "Key Result", endpoint: "key_results"})
}

func (s *Skill) getPersonaTool() skill.Tool {
	return s.getResourceTool(resourceConfig{name: "persona", title: "Persona", endpoint: "personas"})
}

func (s *Skill) getReleaseTool() skill.Tool {
	return s.getResourceTool(resourceConfig{name: "release", title: "Release", endpoint: "releases"})
}

func (s *Skill) getRequirementTool() skill.Tool {
	return s.getResourceTool(resourceConfig{name: "requirement", title: "Requirement", endpoint: "requirements"})
}

func (s *Skill) getTeamTool() skill.Tool {
	return s.getResourceTool(resourceConfig{name: "team", title: "Team", endpoint: "teams"})
}

func (s *Skill) getUserTool() skill.Tool {
	return s.getResourceTool(resourceConfig{name: "user", title: "User", endpoint: "users"})
}

func (s *Skill) getWorkflowTool() skill.Tool {
	return s.getResourceTool(resourceConfig{name: "workflow", title: "Workflow", endpoint: "workflows"})
}

// GraphQL types for search
type graphQLRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

type graphQLResponse struct {
	Data   searchData     `json:"data"`
	Errors []graphQLError `json:"errors,omitempty"`
}

type graphQLError struct {
	Message string `json:"message"`
}

type searchData struct {
	Search searchResults `json:"search"`
}

type searchResults struct {
	Nodes       []documentNode `json:"nodes"`
	CurrentPage int            `json:"currentPage"`
	TotalCount  int            `json:"totalCount"`
	TotalPages  int            `json:"totalPages"`
	IsLastPage  bool           `json:"isLastPage"`
}

type documentNode struct {
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

func (s *Skill) searchDocumentsTool() skill.Tool {
	return skill.NewTool(
		"search_documents",
		"Search for Aha! documents using GraphQL",
		map[string]skill.Parameter{
			"query": {
				Type:        "string",
				Description: "Search query string",
				Required:    true,
			},
			"searchable_type": {
				Type:        "string",
				Description: "Type of document to search for (defaults to Page). Examples: Page, Feature, Epic, Release, etc.",
				Required:    false,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			query, _ := params["query"].(string)
			if query == "" {
				return nil, fmt.Errorf("query is required")
			}

			searchableType := "Page"
			if st, ok := params["searchable_type"].(string); ok && st != "" {
				searchableType = st
			}

			variables := map[string]any{
				"query":          query,
				"searchableType": []string{searchableType},
			}

			graphqlReq := graphQLRequest{
				Query:     searchDocumentsQuery,
				Variables: variables,
			}

			requestBody, err := json.Marshal(graphqlReq)
			if err != nil {
				return nil, fmt.Errorf("error marshaling GraphQL request: %w", err)
			}

			httpReq := httpsimple.Request{
				Method:  http.MethodPost,
				URL:     "/api/graphql",
				Body:    bytes.NewBuffer(requestBody),
				Headers: http.Header{"Content-Type": []string{"application/json"}},
			}

			resp, err := s.simpleClient.Do(ctx, httpReq)
			if err != nil {
				return map[string]any{"error": fmt.Sprintf("error making GraphQL request: %v", err)}, nil
			}
			defer resp.Body.Close()

			responseBody, err := io.ReadAll(resp.Body)
			if err != nil {
				return map[string]any{"error": fmt.Sprintf("error reading response: %v", err)}, nil
			}

			if resp.StatusCode != http.StatusOK {
				return map[string]any{
					"error":       fmt.Sprintf("GraphQL request failed with status %d", resp.StatusCode),
					"status_code": resp.StatusCode,
					"body":        string(responseBody),
				}, nil
			}

			var graphqlResp graphQLResponse
			if err := json.Unmarshal(responseBody, &graphqlResp); err != nil {
				return map[string]any{"error": fmt.Sprintf("error parsing GraphQL response: %v", err)}, nil
			}

			if len(graphqlResp.Errors) > 0 {
				errorMessages := make([]string, len(graphqlResp.Errors))
				for i, err := range graphqlResp.Errors {
					errorMessages[i] = err.Message
				}
				return map[string]any{"errors": errorMessages}, nil
			}

			results := make([]map[string]any, len(graphqlResp.Data.Search.Nodes))
			for i, node := range graphqlResp.Data.Search.Nodes {
				results[i] = map[string]any{
					"reference_num": node.SearchableID,
					"name":          node.Name,
					"type":          node.SearchableType,
					"url":           node.URL,
				}
			}

			return map[string]any{
				"results":       results,
				"total_results": graphqlResp.Data.Search.TotalCount,
				"current_page":  graphqlResp.Data.Search.CurrentPage,
				"total_pages":   graphqlResp.Data.Search.TotalPages,
				"is_last_page":  graphqlResp.Data.Search.IsLastPage,
			}, nil
		},
	)
}
