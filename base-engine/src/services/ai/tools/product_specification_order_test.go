package tools

import (
	"strings"
	"testing"

	"base-engine/auth"
	"base-engine/src/services/ai"
)

func TestSpecificationReorderHasHeadquartersWriteTool(t *testing.T) {
	for _, spec := range Specs() {
		if spec.OperationID != "graphql.mutation.hqReorderSpecificationValues" {
			continue
		}
		if spec.Mode != ai.ModeWrite || spec.Permission != "hqProductCatalog:manage" ||
			len(spec.Workspaces) != 1 || spec.Workspaces[0] != auth.WorkspaceTypeHeadquarters {
			t.Fatalf("incorrect reorder tool scope: %+v", spec)
		}
		for _, field := range []string{"specificationId", "orderedIds"} {
			if !strings.Contains(spec.Document, field) {
				t.Fatalf("reorder tool omits %s", field)
			}
		}
		return
	}
	t.Fatal("headquarters specification reorder tool missing")
}
