// Command mcp-aha runs an Aha! MCP server that exposes tools for
// reading product management data from Aha!
// It can also be used as a CLI tool for testing and scripting.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	runtime "github.com/plexusone/omniskill/mcp/server"
	"github.com/plexusone/omnitoken"
	"github.com/spf13/cobra"

	ahaskill "github.com/grokify/aha-mcp-server/skills/aha"

	// Register desktop vault providers (1Password, etc.)
	_ "github.com/plexusone/omnivault-desktop"
)

const (
	serverName    = "mcp-aha"
	serverVersion = "v0.7.0"
)

var (
	// Credential flags (persistent across all commands)
	subdomain       string
	apiKey          string
	vaultURI        string
	credentialsName string

	// Output format flag
	outputFormat string

	// list-ideas flags
	listQuery          string
	listSpam           bool
	listWorkflowStatus string
	listSort           string
	listCreatedBefore  string
	listCreatedSince   string
	listUpdatedSince   string
	listTag            string
	listUserID         string
	listPage           int
	listPerPage        int

	// search-documents flags
	searchQuery         string
	searchSearchableType string
	searchLimit         int
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "mcp-aha",
	Short: "MCP server and CLI for Aha!",
	Long: `An MCP (Model Context Protocol) server for reading product management data from Aha!
Can also be used as a CLI tool for testing and scripting.

Running without a subcommand starts the MCP server (default behavior).

Credentials can be provided via:
  - Direct credentials (subdomain and API key)
  - Vault-backed credentials via omnitoken`,
	Example: `  # Start MCP server (default)
  mcp-aha --subdomain mycompany --api-key your-api-key

  # CLI: Get an idea
  mcp-aha get-idea IDEA-123 --subdomain mycompany --api-key your-api-key

  # CLI: List ideas with search
  mcp-aha list-ideas --query "authentication" --subdomain mycompany --api-key your-api-key

  # CLI: Get a feature
  mcp-aha get-feature FEAT-456 --subdomain mycompany --api-key your-api-key`,
	SilenceUsage: true,
	RunE:         runServer,
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the MCP server",
	Long:  "Start the MCP server using stdio transport for communication with MCP clients.",
	RunE:  runServer,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s %s\n", serverName, serverVersion)
	},
}

// CLI commands for Aha! tools
var getIdeaCmd = &cobra.Command{
	Use:   "get-idea <idea-id>",
	Short: "Get an idea by ID",
	Long:  "Retrieve an Aha! idea by its ID.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTool("get_idea", map[string]any{
			"idea_id": args[0],
		})
	},
}

var listIdeasCmd = &cobra.Command{
	Use:   "list-ideas",
	Short: "List ideas with filtering",
	Long:  "List ideas from Aha! with optional filtering and pagination.",
	RunE: func(cmd *cobra.Command, args []string) error {
		params := map[string]any{}
		if listQuery != "" {
			params["q"] = listQuery
		}
		if listSpam {
			params["spam"] = listSpam
		}
		if listWorkflowStatus != "" {
			params["workflow_status"] = listWorkflowStatus
		}
		if listSort != "" {
			params["sort"] = listSort
		}
		if listCreatedBefore != "" {
			params["created_before"] = listCreatedBefore
		}
		if listCreatedSince != "" {
			params["created_since"] = listCreatedSince
		}
		if listUpdatedSince != "" {
			params["updated_since"] = listUpdatedSince
		}
		if listTag != "" {
			params["tag"] = listTag
		}
		if listUserID != "" {
			params["user_id"] = listUserID
		}
		if listPage > 0 {
			params["page"] = listPage
		}
		if listPerPage > 0 {
			params["per_page"] = listPerPage
		}
		return runTool("list_ideas", params)
	},
}

var getFeatureCmd = &cobra.Command{
	Use:   "get-feature <feature-id>",
	Short: "Get a feature by ID",
	Long:  "Retrieve an Aha! feature by its ID.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTool("get_feature", map[string]any{
			"feature_id": args[0],
		})
	},
}

var getEpicCmd = &cobra.Command{
	Use:   "get-epic <epic-id>",
	Short: "Get an epic by ID",
	Long:  "Retrieve an Aha! epic by its ID.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTool("get_epic", map[string]any{
			"epic_id": args[0],
		})
	},
}

var getReleaseCmd = &cobra.Command{
	Use:   "get-release <release-id>",
	Short: "Get a release by ID",
	Long:  "Retrieve an Aha! release by its ID.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTool("get_release", map[string]any{
			"release_id": args[0],
		})
	},
}

var getGoalCmd = &cobra.Command{
	Use:   "get-goal <goal-id>",
	Short: "Get a goal by ID",
	Long:  "Retrieve an Aha! goal by its ID.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTool("get_goal", map[string]any{
			"goal_id": args[0],
		})
	},
}

var getInitiativeCmd = &cobra.Command{
	Use:   "get-initiative <initiative-id>",
	Short: "Get an initiative by ID",
	Long:  "Retrieve an Aha! initiative by its ID.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTool("get_initiative", map[string]any{
			"initiative_id": args[0],
		})
	},
}

var getCommentCmd = &cobra.Command{
	Use:   "get-comment <comment-id>",
	Short: "Get a comment by ID",
	Long:  "Retrieve an Aha! comment by its ID.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTool("get_comment", map[string]any{
			"comment_id": args[0],
		})
	},
}

var getKeyResultCmd = &cobra.Command{
	Use:   "get-key-result <key-result-id>",
	Short: "Get a key result by ID",
	Long:  "Retrieve an Aha! key result by its ID.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTool("get_key_result", map[string]any{
			"key_result_id": args[0],
		})
	},
}

var getPersonaCmd = &cobra.Command{
	Use:   "get-persona <persona-id>",
	Short: "Get a persona by ID",
	Long:  "Retrieve an Aha! persona by its ID.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTool("get_persona", map[string]any{
			"persona_id": args[0],
		})
	},
}

var getRequirementCmd = &cobra.Command{
	Use:   "get-requirement <requirement-id>",
	Short: "Get a requirement by ID",
	Long:  "Retrieve an Aha! requirement by its ID.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTool("get_requirement", map[string]any{
			"requirement_id": args[0],
		})
	},
}

var getTeamCmd = &cobra.Command{
	Use:   "get-team <team-id>",
	Short: "Get a team by ID",
	Long:  "Retrieve an Aha! team by its ID.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTool("get_team", map[string]any{
			"team_id": args[0],
		})
	},
}

var getUserCmd = &cobra.Command{
	Use:   "get-user <user-id>",
	Short: "Get a user by ID",
	Long:  "Retrieve an Aha! user by their ID.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTool("get_user", map[string]any{
			"user_id": args[0],
		})
	},
}

var getWorkflowCmd = &cobra.Command{
	Use:   "get-workflow <workflow-id>",
	Short: "Get a workflow by ID",
	Long:  "Retrieve an Aha! workflow by its ID.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTool("get_workflow", map[string]any{
			"workflow_id": args[0],
		})
	},
}

var searchDocumentsCmd = &cobra.Command{
	Use:   "search-documents <query>",
	Short: "Search documents",
	Long:  "Search for documents across your Aha! workspace using GraphQL.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		params := map[string]any{
			"query": args[0],
		}
		if searchSearchableType != "" {
			params["searchable_type"] = searchSearchableType
		}
		if searchLimit > 0 {
			params["limit"] = searchLimit
		}
		return runTool("search_documents", params)
	},
}

func init() {
	// Persistent flags (available to all commands)
	rootCmd.PersistentFlags().StringVar(&subdomain, "subdomain", "",
		"Aha! subdomain (env: AHA_DOMAIN)")
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "",
		"Aha! API key (env: AHA_API_TOKEN)")
	rootCmd.PersistentFlags().StringVar(&vaultURI, "vault", "",
		"vault URI for credentials (env: OMNITOKEN_VAULT_URI)")
	rootCmd.PersistentFlags().StringVar(&credentialsName, "credentials-name", "",
		"name of credentials in vault (env: OMNITOKEN_CREDENTIALS_NAME)")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "json",
		"output format: json, pretty (default: json)")

	// list-ideas flags
	listIdeasCmd.Flags().StringVarP(&listQuery, "query", "q", "", "search term to match against idea name")
	listIdeasCmd.Flags().BoolVar(&listSpam, "spam", false, "show ideas marked as spam")
	listIdeasCmd.Flags().StringVar(&listWorkflowStatus, "workflow-status", "", "filter by workflow status")
	listIdeasCmd.Flags().StringVar(&listSort, "sort", "", "sort by: recent, trending, or popular")
	listIdeasCmd.Flags().StringVar(&listCreatedBefore, "created-before", "", "only ideas created before (ISO8601)")
	listIdeasCmd.Flags().StringVar(&listCreatedSince, "created-since", "", "only ideas created after (ISO8601)")
	listIdeasCmd.Flags().StringVar(&listUpdatedSince, "updated-since", "", "only ideas updated after (ISO8601)")
	listIdeasCmd.Flags().StringVar(&listTag, "tag", "", "filter by tag")
	listIdeasCmd.Flags().StringVar(&listUserID, "user-id", "", "filter by creator user ID")
	listIdeasCmd.Flags().IntVar(&listPage, "page", 0, "page number")
	listIdeasCmd.Flags().IntVar(&listPerPage, "per-page", 0, "results per page")

	// search-documents flags
	searchDocumentsCmd.Flags().StringVar(&searchSearchableType, "type", "", "document type (e.g., Page)")
	searchDocumentsCmd.Flags().IntVar(&searchLimit, "limit", 0, "maximum results to return")

	// Add commands
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(versionCmd)

	// Aha! CLI commands
	rootCmd.AddCommand(getIdeaCmd)
	rootCmd.AddCommand(listIdeasCmd)
	rootCmd.AddCommand(getFeatureCmd)
	rootCmd.AddCommand(getEpicCmd)
	rootCmd.AddCommand(getReleaseCmd)
	rootCmd.AddCommand(getGoalCmd)
	rootCmd.AddCommand(getInitiativeCmd)
	rootCmd.AddCommand(getCommentCmd)
	rootCmd.AddCommand(getKeyResultCmd)
	rootCmd.AddCommand(getPersonaCmd)
	rootCmd.AddCommand(getRequirementCmd)
	rootCmd.AddCommand(getTeamCmd)
	rootCmd.AddCommand(getUserCmd)
	rootCmd.AddCommand(getWorkflowCmd)
	rootCmd.AddCommand(searchDocumentsCmd)
}

// applyEnvDefaults applies environment variable defaults to flags
func applyEnvDefaults() {
	if subdomain == "" {
		subdomain = os.Getenv("AHA_DOMAIN")
	}
	if apiKey == "" {
		apiKey = os.Getenv("AHA_API_TOKEN")
	}
	if vaultURI == "" {
		vaultURI = os.Getenv("OMNITOKEN_VAULT_URI")
	}
	if credentialsName == "" {
		credentialsName = os.Getenv("OMNITOKEN_CREDENTIALS_NAME")
	}
	if credentialsName == "" {
		credentialsName = "aha"
	}
}

// getSkill creates and initializes an Aha skill with proper credentials
func getSkill(ctx context.Context) (*ahaskill.Skill, func(), error) {
	applyEnvDefaults()

	hasDirectCreds := subdomain != "" && apiKey != ""
	hasVaultCreds := vaultURI != ""

	var skill *ahaskill.Skill
	var tokenMgr *omnitoken.TokenManager
	var err error

	cleanup := func() {}

	if hasVaultCreds {
		tokenMgr, err = omnitoken.NewFromVaultURI(vaultURI)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create token manager: %w", err)
		}
		cleanup = func() {
			if err := tokenMgr.Close(); err != nil {
				log.Printf("Warning: failed to close token manager: %v", err)
			}
		}

		creds, err := tokenMgr.GetCredentials(ctx, credentialsName)
		if err != nil {
			cleanup()
			return nil, nil, fmt.Errorf("failed to get credentials: %w", err)
		}

		if subdomain == "" {
			subdomain = creds.Subdomain
		}

		if apiKey == "" && creds.HeaderQuery != nil {
			if authHeader := creds.HeaderQuery.Header.Get("Authorization"); authHeader != "" {
				if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
					apiKey = authHeader[7:]
				} else {
					apiKey = authHeader
				}
			}
		}

		if apiKey == "" && creds.Additional != nil {
			if key := creds.Additional.Get("api_key"); key != "" {
				apiKey = key
			} else if key := creds.Additional.Get("api_token"); key != "" {
				apiKey = key
			}
		}

		if subdomain == "" || apiKey == "" {
			cleanup()
			return nil, nil, fmt.Errorf("vault credentials must contain subdomain and api_key")
		}

		skill, err = ahaskill.New(subdomain, apiKey)
		if err != nil {
			cleanup()
			return nil, nil, fmt.Errorf("failed to create Aha skill: %w", err)
		}
	} else if hasDirectCreds {
		skill, err = ahaskill.New(subdomain, apiKey)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create Aha skill: %w", err)
		}
	} else {
		return nil, nil, fmt.Errorf("credentials required: use --subdomain/--api-key or --vault")
	}

	if err := skill.Init(ctx); err != nil {
		cleanup()
		return nil, nil, fmt.Errorf("failed to initialize Aha skill: %w", err)
	}

	fullCleanup := func() {
		if err := skill.Close(); err != nil {
			log.Printf("Warning: failed to close Aha skill: %v", err)
		}
		cleanup()
	}

	return skill, fullCleanup, nil
}

// outputResult outputs the result in the specified format
func outputResult(result any) error {
	var data []byte
	var err error

	switch outputFormat {
	case "pretty":
		data, err = json.MarshalIndent(result, "", "  ")
	default:
		data, err = json.Marshal(result)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	fmt.Println(string(data))
	return nil
}

// runTool runs a tool by name with the given params
func runTool(toolName string, params map[string]any) error {
	ctx := context.Background()

	skill, cleanup, err := getSkill(ctx)
	if err != nil {
		return err
	}
	defer cleanup()

	for _, tool := range skill.Tools() {
		if tool.Name() == toolName {
			result, err := tool.Call(ctx, params)
			if err != nil {
				return err
			}
			return outputResult(result)
		}
	}
	return fmt.Errorf("tool not found: %s", toolName)
}

func runServer(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	skill, cleanup, err := getSkill(ctx)
	if err != nil {
		return err
	}
	defer cleanup()

	rt := runtime.New(&mcp.Implementation{
		Name:    serverName,
		Version: serverVersion,
	}, nil)

	rt.RegisterSkill(skill)

	if err := rt.ServeStdio(ctx); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}
