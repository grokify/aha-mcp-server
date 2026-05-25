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
	return skill.NewTool(
		"get_comment",
		"Get Comment from Aha",
		map[string]skill.Parameter{
			"comment_id": {
				Type:        "string",
				Description: "Comment ID to get",
				Required:    true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			commentID, _ := params["comment_id"].(string)
			if commentID == "" {
				return nil, fmt.Errorf("comment_id is required")
			}

			resp, err := s.simpleClient.Do(ctx, httpsimple.Request{
				Method: http.MethodGet,
				URL:    fmt.Sprintf("/api/v1/comments/%s", commentID),
			})
			if err != nil {
				return map[string]any{"error": fmt.Sprintf("error getting Comment: %v", err)}, nil
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return map[string]any{"error": fmt.Sprintf("error reading response: %v", err)}, nil
			}

			var comment any
			if err := json.Unmarshal(body, &comment); err != nil {
				return map[string]any{
					"comment":     string(body),
					"status_code": resp.StatusCode,
				}, nil
			}

			return map[string]any{
				"comment":     comment,
				"status_code": resp.StatusCode,
			}, nil
		},
	)
}

func (s *Skill) getEpicTool() skill.Tool {
	return skill.NewTool(
		"get_epic",
		"Get Epic from Aha",
		map[string]skill.Parameter{
			"epic_id": {
				Type:        "string",
				Description: "Epic ID to get",
				Required:    true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			epicID, _ := params["epic_id"].(string)
			if epicID == "" {
				return nil, fmt.Errorf("epic_id is required")
			}

			resp, err := s.simpleClient.Do(ctx, httpsimple.Request{
				Method: http.MethodGet,
				URL:    fmt.Sprintf("/api/v1/epics/%s", epicID),
			})
			if err != nil {
				return map[string]any{"error": fmt.Sprintf("error getting Epic: %v", err)}, nil
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return map[string]any{"error": fmt.Sprintf("error reading response: %v", err)}, nil
			}

			var epic any
			if err := json.Unmarshal(body, &epic); err != nil {
				return map[string]any{
					"epic":        string(body),
					"status_code": resp.StatusCode,
				}, nil
			}

			return map[string]any{
				"epic":        epic,
				"status_code": resp.StatusCode,
			}, nil
		},
	)
}

func (s *Skill) getFeatureTool() skill.Tool {
	return skill.NewTool(
		"get_feature",
		"Get Feature from Aha",
		map[string]skill.Parameter{
			"feature_id": {
				Type:        "string",
				Description: "Feature ID to get",
				Required:    true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			featureID, _ := params["feature_id"].(string)
			if featureID == "" {
				return nil, fmt.Errorf("feature_id is required")
			}

			resp, err := s.simpleClient.Do(ctx, httpsimple.Request{
				Method: http.MethodGet,
				URL:    fmt.Sprintf("/api/v1/features/%s", featureID),
			})
			if err != nil {
				return map[string]any{"error": fmt.Sprintf("error getting Feature: %v", err)}, nil
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return map[string]any{"error": fmt.Sprintf("error reading response: %v", err)}, nil
			}

			var feature any
			if err := json.Unmarshal(body, &feature); err != nil {
				return map[string]any{
					"feature":     string(body),
					"status_code": resp.StatusCode,
				}, nil
			}

			return map[string]any{
				"feature":     feature,
				"status_code": resp.StatusCode,
			}, nil
		},
	)
}

func (s *Skill) getGoalTool() skill.Tool {
	return skill.NewTool(
		"get_goal",
		"Get Goal from Aha",
		map[string]skill.Parameter{
			"goal_id": {
				Type:        "string",
				Description: "Goal ID to get",
				Required:    true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			goalID, _ := params["goal_id"].(string)
			if goalID == "" {
				return nil, fmt.Errorf("goal_id is required")
			}

			resp, err := s.simpleClient.Do(ctx, httpsimple.Request{
				Method: http.MethodGet,
				URL:    fmt.Sprintf("/api/v1/goals/%s", goalID),
			})
			if err != nil {
				return map[string]any{"error": fmt.Sprintf("error getting Goal: %v", err)}, nil
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return map[string]any{"error": fmt.Sprintf("error reading response: %v", err)}, nil
			}

			var goal any
			if err := json.Unmarshal(body, &goal); err != nil {
				return map[string]any{
					"goal":        string(body),
					"status_code": resp.StatusCode,
				}, nil
			}

			return map[string]any{
				"goal":        goal,
				"status_code": resp.StatusCode,
			}, nil
		},
	)
}

func (s *Skill) getInitiativeTool() skill.Tool {
	return skill.NewTool(
		"get_initiative",
		"Get Initiative from Aha",
		map[string]skill.Parameter{
			"initiative_id": {
				Type:        "string",
				Description: "Initiative ID to get",
				Required:    true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			initiativeID, _ := params["initiative_id"].(string)
			if initiativeID == "" {
				return nil, fmt.Errorf("initiative_id is required")
			}

			resp, err := s.simpleClient.Do(ctx, httpsimple.Request{
				Method: http.MethodGet,
				URL:    fmt.Sprintf("/api/v1/initiatives/%s", initiativeID),
			})
			if err != nil {
				return map[string]any{"error": fmt.Sprintf("error getting Initiative: %v", err)}, nil
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return map[string]any{"error": fmt.Sprintf("error reading response: %v", err)}, nil
			}

			var initiative any
			if err := json.Unmarshal(body, &initiative); err != nil {
				return map[string]any{
					"initiative":  string(body),
					"status_code": resp.StatusCode,
				}, nil
			}

			return map[string]any{
				"initiative":  initiative,
				"status_code": resp.StatusCode,
			}, nil
		},
	)
}

func (s *Skill) getKeyResultTool() skill.Tool {
	return skill.NewTool(
		"get_key_result",
		"Get Key Result from Aha",
		map[string]skill.Parameter{
			"key_result_id": {
				Type:        "string",
				Description: "Key Result ID to get",
				Required:    true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			keyResultID, _ := params["key_result_id"].(string)
			if keyResultID == "" {
				return nil, fmt.Errorf("key_result_id is required")
			}

			resp, err := s.simpleClient.Do(ctx, httpsimple.Request{
				Method: http.MethodGet,
				URL:    fmt.Sprintf("/api/v1/key_results/%s", keyResultID),
			})
			if err != nil {
				return map[string]any{"error": fmt.Sprintf("error getting Key Result: %v", err)}, nil
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return map[string]any{"error": fmt.Sprintf("error reading response: %v", err)}, nil
			}

			var keyResult any
			if err := json.Unmarshal(body, &keyResult); err != nil {
				return map[string]any{
					"key_result":  string(body),
					"status_code": resp.StatusCode,
				}, nil
			}

			return map[string]any{
				"key_result":  keyResult,
				"status_code": resp.StatusCode,
			}, nil
		},
	)
}

func (s *Skill) getPersonaTool() skill.Tool {
	return skill.NewTool(
		"get_persona",
		"Get Persona from Aha",
		map[string]skill.Parameter{
			"persona_id": {
				Type:        "string",
				Description: "Persona ID to get",
				Required:    true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			personaID, _ := params["persona_id"].(string)
			if personaID == "" {
				return nil, fmt.Errorf("persona_id is required")
			}

			resp, err := s.simpleClient.Do(ctx, httpsimple.Request{
				Method: http.MethodGet,
				URL:    fmt.Sprintf("/api/v1/personas/%s", personaID),
			})
			if err != nil {
				return map[string]any{"error": fmt.Sprintf("error getting Persona: %v", err)}, nil
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return map[string]any{"error": fmt.Sprintf("error reading response: %v", err)}, nil
			}

			var persona any
			if err := json.Unmarshal(body, &persona); err != nil {
				return map[string]any{
					"persona":     string(body),
					"status_code": resp.StatusCode,
				}, nil
			}

			return map[string]any{
				"persona":     persona,
				"status_code": resp.StatusCode,
			}, nil
		},
	)
}

func (s *Skill) getReleaseTool() skill.Tool {
	return skill.NewTool(
		"get_release",
		"Get Release from Aha",
		map[string]skill.Parameter{
			"release_id": {
				Type:        "string",
				Description: "Release ID to get",
				Required:    true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			releaseID, _ := params["release_id"].(string)
			if releaseID == "" {
				return nil, fmt.Errorf("release_id is required")
			}

			resp, err := s.simpleClient.Do(ctx, httpsimple.Request{
				Method: http.MethodGet,
				URL:    fmt.Sprintf("/api/v1/releases/%s", releaseID),
			})
			if err != nil {
				return map[string]any{"error": fmt.Sprintf("error getting Release: %v", err)}, nil
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return map[string]any{"error": fmt.Sprintf("error reading response: %v", err)}, nil
			}

			var release any
			if err := json.Unmarshal(body, &release); err != nil {
				return map[string]any{
					"release":     string(body),
					"status_code": resp.StatusCode,
				}, nil
			}

			return map[string]any{
				"release":     release,
				"status_code": resp.StatusCode,
			}, nil
		},
	)
}

func (s *Skill) getRequirementTool() skill.Tool {
	return skill.NewTool(
		"get_requirement",
		"Get Requirement from Aha",
		map[string]skill.Parameter{
			"requirement_id": {
				Type:        "string",
				Description: "Requirement ID to get",
				Required:    true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			requirementID, _ := params["requirement_id"].(string)
			if requirementID == "" {
				return nil, fmt.Errorf("requirement_id is required")
			}

			resp, err := s.simpleClient.Do(ctx, httpsimple.Request{
				Method: http.MethodGet,
				URL:    fmt.Sprintf("/api/v1/requirements/%s", requirementID),
			})
			if err != nil {
				return map[string]any{"error": fmt.Sprintf("error getting Requirement: %v", err)}, nil
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return map[string]any{"error": fmt.Sprintf("error reading response: %v", err)}, nil
			}

			var requirement any
			if err := json.Unmarshal(body, &requirement); err != nil {
				return map[string]any{
					"requirement": string(body),
					"status_code": resp.StatusCode,
				}, nil
			}

			return map[string]any{
				"requirement": requirement,
				"status_code": resp.StatusCode,
			}, nil
		},
	)
}

func (s *Skill) getTeamTool() skill.Tool {
	return skill.NewTool(
		"get_team",
		"Get Team from Aha",
		map[string]skill.Parameter{
			"team_id": {
				Type:        "string",
				Description: "Team ID to get",
				Required:    true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			teamID, _ := params["team_id"].(string)
			if teamID == "" {
				return nil, fmt.Errorf("team_id is required")
			}

			resp, err := s.simpleClient.Do(ctx, httpsimple.Request{
				Method: http.MethodGet,
				URL:    fmt.Sprintf("/api/v1/teams/%s", teamID),
			})
			if err != nil {
				return map[string]any{"error": fmt.Sprintf("error getting Team: %v", err)}, nil
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return map[string]any{"error": fmt.Sprintf("error reading response: %v", err)}, nil
			}

			var team any
			if err := json.Unmarshal(body, &team); err != nil {
				return map[string]any{
					"team":        string(body),
					"status_code": resp.StatusCode,
				}, nil
			}

			return map[string]any{
				"team":        team,
				"status_code": resp.StatusCode,
			}, nil
		},
	)
}

func (s *Skill) getUserTool() skill.Tool {
	return skill.NewTool(
		"get_user",
		"Get User from Aha",
		map[string]skill.Parameter{
			"user_id": {
				Type:        "string",
				Description: "User ID to get",
				Required:    true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			userID, _ := params["user_id"].(string)
			if userID == "" {
				return nil, fmt.Errorf("user_id is required")
			}

			resp, err := s.simpleClient.Do(ctx, httpsimple.Request{
				Method: http.MethodGet,
				URL:    fmt.Sprintf("/api/v1/users/%s", userID),
			})
			if err != nil {
				return map[string]any{"error": fmt.Sprintf("error getting User: %v", err)}, nil
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return map[string]any{"error": fmt.Sprintf("error reading response: %v", err)}, nil
			}

			var user any
			if err := json.Unmarshal(body, &user); err != nil {
				return map[string]any{
					"user":        string(body),
					"status_code": resp.StatusCode,
				}, nil
			}

			return map[string]any{
				"user":        user,
				"status_code": resp.StatusCode,
			}, nil
		},
	)
}

func (s *Skill) getWorkflowTool() skill.Tool {
	return skill.NewTool(
		"get_workflow",
		"Get Workflow from Aha",
		map[string]skill.Parameter{
			"workflow_id": {
				Type:        "string",
				Description: "Workflow ID to get",
				Required:    true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			workflowID, _ := params["workflow_id"].(string)
			if workflowID == "" {
				return nil, fmt.Errorf("workflow_id is required")
			}

			resp, err := s.simpleClient.Do(ctx, httpsimple.Request{
				Method: http.MethodGet,
				URL:    fmt.Sprintf("/api/v1/workflows/%s", workflowID),
			})
			if err != nil {
				return map[string]any{"error": fmt.Sprintf("error getting Workflow: %v", err)}, nil
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return map[string]any{"error": fmt.Sprintf("error reading response: %v", err)}, nil
			}

			var workflow any
			if err := json.Unmarshal(body, &workflow); err != nil {
				return map[string]any{
					"workflow":    string(body),
					"status_code": resp.StatusCode,
				}, nil
			}

			return map[string]any{
				"workflow":    workflow,
				"status_code": resp.StatusCode,
			}, nil
		},
	)
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
