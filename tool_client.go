package ahamcpserver

import (
	"github.com/grokify/go-aha/v3/oag7/aha"
	"github.com/grokify/go-aha/v3/oag7/client"
	"github.com/grokify/mogo/net/http/httpsimple"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ToolsClient struct {
	client       *aha.APIClient
	config       *aha.Configuration
	simpleClient *httpsimple.Client
}

func NewToolsClient(ahaSubdomain, ahaAPIKey string) (*ToolsClient, error) {
	config, err := client.NewConfiguration(ahaSubdomain, ahaAPIKey)
	if err != nil {
		return nil, err
	}
	sc, err := client.NewSimpleClient(ahaSubdomain, ahaAPIKey)
	if err != nil {
		return nil, err
	}
	return &ToolsClient{
		client:       aha.NewAPIClient(config),
		config:       config,
		simpleClient: sc,
	}, nil
}

func (tc *ToolsClient) AddTools(svr *mcp.Server) {
	mcp.AddTool(svr, GetFeatureTool(), tc.GetFeature)
	mcp.AddTool(svr, GetIdeaTool(), tc.GetIdea)
}
