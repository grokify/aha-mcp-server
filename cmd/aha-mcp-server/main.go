package main

import (
	"context"
	"log"

	ahamcpserver "github.com/grokify/aha-mcp-server"
)

func main() {
	opts, err := ahamcpserver.NewOptionsEnv()
	if err != nil {
		log.Fatal(err)
	}

	ahamcpserver.ListenAndServe(context.Background(), opts)
}
