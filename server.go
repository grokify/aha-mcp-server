package ahamcpserver

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/grokify/mogo/net/http/httputilmore"
	"github.com/jessevdk/go-flags"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/grokify/aha-mcp-server/tools"
)

type Options struct {
	HTTPAddr     string `short:"h" long:"http" description:"HTTP address: if set, use streamable HTTP at this address, instead of stdin/stdout"`
	AHASubdomain string
	AHAAPIKey    string
}

func NewOptionsEnv() (*Options, error) {
	opts := &Options{}
	_, err := flags.Parse(opts)
	if err != nil {
		return opts, err
	}
	opts.ReadEnvDefaults()
	return opts, opts.CheckCredentials()
}

func (opts *Options) ReadEnvDefaults() {
	if opts.AHASubdomain == "" {
		opts.AHASubdomain = os.Getenv("AHA_DOMAIN")
	}
	if opts.AHAAPIKey == "" {
		opts.AHAAPIKey = os.Getenv("AHA_API_TOKEN")
	}
}

func (opts *Options) CheckCredentials() error {
	if opts.AHASubdomain == "" {
		return errors.New("AHA_DOMAIN environment variable is required")
	}
	if opts.AHAAPIKey == "" {
		return errors.New("AHA_API_TOKEN environment variable is required")
	}
	return nil
}

func ListenAndServe(ctx context.Context, opts *Options) {
	svr := mcp.NewServer(
		&mcp.Implementation{
			Name:    "aha-mcp-server",
			Title:   "aha-mcp-server",
			Version: "0.3.0"}, nil)

	if toolsClient, err := tools.NewToolsClient(opts.AHASubdomain, opts.AHAAPIKey); err != nil {
		log.Fatal(err)
	} else {
		toolsClient.AddTools(svr)
	}

	if opts.HTTPAddr != "" {
		handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
			return svr
		}, nil)
		log.Printf("MCP handler listening at %s", opts.HTTPAddr)
		httpSvr := httputilmore.NewServerTimeouts(opts.HTTPAddr, handler, time.Second*5)
		log.Fatal(httpSvr.ListenAndServe())
	} else {
		t := mcp.NewLoggingTransport(mcp.NewStdioTransport(), os.Stderr)
		if err := svr.Run(ctx, t); err != nil {
			log.Printf("Server failed: %v", err)
		}
	}
}
