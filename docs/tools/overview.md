# Tools Reference

This server provides 14 tools for accessing and searching Aha! data.

## Search Tools

### search_documents

Search for documents across your Aha! workspace using GraphQL.

**Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `query` | string | Yes | Search query string |
| `searchable_type` | string | No | Document type to search (default: "Page") |
| `limit` | integer | No | Maximum number of results |

**Example:**

```json
{
  "tool": "search_documents",
  "arguments": {
    "query": "product roadmap",
    "searchable_type": "Page",
    "limit": 10
  }
}
```

## Ideas Tools

### get_idea

Retrieve a specific idea by ID.

**Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `idea_id` | string | Yes | The idea ID (e.g., "IDEA-123") |

### list_ideas

List ideas with filtering and pagination.

**Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `q` | string | No | Search term to match against idea name |
| `spam` | boolean | No | Show ideas marked as spam |
| `workflow_status` | string | No | Filter by workflow status |
| `sort` | string | No | Sort by: recent, trending, or popular |
| `created_before` | string | No | Only ideas created before (ISO8601) |
| `created_since` | string | No | Only ideas created after (ISO8601) |
| `updated_since` | string | No | Only ideas updated after (ISO8601) |
| `tag` | string | No | Filter by tag |
| `user_id` | string | No | Filter by creator user ID |
| `page` | integer | No | Page number |
| `per_page` | integer | No | Results per page |

## Object Retrieval Tools

All object retrieval tools follow the same pattern:

| Tool | Parameter | Description |
|------|-----------|-------------|
| `get_feature` | `feature_id` | Retrieve a feature by ID |
| `get_epic` | `epic_id` | Retrieve an epic by ID |
| `get_release` | `release_id` | Retrieve a release by ID |
| `get_goal` | `goal_id` | Retrieve a goal by ID |
| `get_initiative` | `initiative_id` | Retrieve an initiative by ID |
| `get_key_result` | `key_result_id` | Retrieve a key result by ID |
| `get_persona` | `persona_id` | Retrieve a persona by ID |
| `get_requirement` | `requirement_id` | Retrieve a requirement by ID |
| `get_team` | `team_id` | Retrieve a team by ID |
| `get_user` | `user_id` | Retrieve a user by ID |
| `get_workflow` | `workflow_id` | Retrieve a workflow by ID |
| `get_comment` | `comment_id` | Retrieve a comment by ID |

**Example:**

```json
{
  "tool": "get_feature",
  "arguments": {
    "feature_id": "FEAT-123"
  }
}
```

## Response Format

All tools return JSON data including:

- The requested object data
- HTTP status code
- Any error messages

**Example Response:**

```json
{
  "feature": {
    "id": "6789012345",
    "reference_num": "FEAT-123",
    "name": "User Authentication",
    "workflow_status": {
      "name": "In Development"
    },
    "description": {
      "body": "Implement OAuth2 authentication..."
    }
  },
  "status_code": 200
}
```
