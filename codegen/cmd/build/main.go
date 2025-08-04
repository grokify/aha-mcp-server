package main

import (
	"fmt"
	"log"

	"github.com/grokify/aha-mcp-server/codegen"
)

func main() {
	err := codegen.BuildCodeToolsGet("tools")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DONE")
}
