package tools

import "github.com/modelcontextprotocol/go-sdk/mcp"

func (tc *ToolsClient) AddTools(svr *mcp.Server) {
	mcp.AddTool(svr, GetIdeaTool(), tc.GetIdea)

	mcp.AddTool(svr, GetCommentTool(), tc.GetComment)
	mcp.AddTool(svr, GetEpicTool(), tc.GetEpic)
	mcp.AddTool(svr, GetFeatureTool(), tc.GetFeature)
	mcp.AddTool(svr, GetGoalTool(), tc.GetGoal)
	mcp.AddTool(svr, GetInitiativeTool(), tc.GetInitiative)
	mcp.AddTool(svr, GetKeyResultTool(), tc.GetKeyResult)
	mcp.AddTool(svr, GetPersonaTool(), tc.GetPersona)
	mcp.AddTool(svr, GetReleaseTool(), tc.GetRelease)
	mcp.AddTool(svr, GetRequirementTool(), tc.GetRequirement)
	mcp.AddTool(svr, SearchDocumentsTool(), tc.SearchDocuments)
	mcp.AddTool(svr, GetTeamTool(), tc.GetTeam)
	mcp.AddTool(svr, GetUserTool(), tc.GetUser)
	mcp.AddTool(svr, GetWorkflowTool(), tc.GetWorkflow)
}
