package codegen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/grokify/mogo/text/stringcase"
)

func BuildCodeClientToolsGet(dir string) error {
	objs := Objects()
	ht := ClientAddTools(objs)
	fn := "client_tools_add.go"
	fp := filepath.Join(dir, fn)
	return os.WriteFile(fp, []byte(ht), 0600)
}

func BuildCodeToolsGet(dir string) error {
	if d := strings.TrimSpace(dir); d == "" {
		dir = "."
	}
	if err := BuildCodeClientToolsGet(dir); err != nil {
		return err
	}
	objs := Objects()
	for _, obj := range objs {
		ht := ToolGetObject(obj)
		fn := fmt.Sprintf("tool_get_%s.go", obj.Singular.SnakeCase())
		fp := filepath.Join(dir, fn)
		if err := os.WriteFile(fp, []byte(ht), 0600); err != nil {
			return err
		}
	}
	return nil
}

func Objects() []stringcase.CaserNoun {
	raw := [][]string{
		{"Comment", "Comments"},
		{"Epic", "Epics"},
		{"Feature", "Features"},
		{"Goal", "Goals"},
		{"Initiative", "Initiatives"},
		{"Key Result", "Key Results"},
		{"Persona", "Personas"},
		{"Release", "Releases"},
		{"Requirement", "Requirements"},
		{"Team", "Teams"},
		{"User", "Users"},
		{"Workflow", "Workflows"},
	}
	var objs []stringcase.CaserNoun
	for _, r := range raw {
		objs = append(objs, stringcase.NewCaserNoun(r[0], r[1]))
	}
	return objs
}
