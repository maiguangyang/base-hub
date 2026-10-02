package tools

import (
	"strings"

	"base-engine/auth"
	"base-engine/src/services/ai"
)

func reviewedSpec(id, operation, document string, mode ai.ToolMode, risk, permission string, workspace auth.WorkspaceType, kind ai.WriteKind, targets, arguments []string) ai.ToolSpec {
	title, description := toolMetadata(id, mode, workspace)
	root := operation[strings.LastIndex(operation, ".")+1:]
	fields := []string{root}
	if mode == ai.ModeWrite {
		fields = append(fields, "id")
	}
	return ai.ToolSpec{ID: id, Name: id, Title: title, OperationID: operation, Document: document, Description: description,
		Mode: mode, Risk: risk, Permission: permission, Workspaces: []auth.WorkspaceType{workspace},
		OutputFields: fields, WriteKind: kind, TargetFields: targets, ArgumentFields: arguments}
}

func remainingSpecs() []ai.ToolSpec {
	var all []ai.ToolSpec
	for _, group := range [][]ai.ToolSpec{task8Specs(), task9Specs(), task10Specs(), task11Specs(), task12Specs(), task13Specs(), paymentConfigSpecs(), contractGapSpecs(), customerReadSpecs(), customerWriteSpecs(), productReadSpecs(), productWriteSpecs()} {
		all = append(all, group...)
	}
	return all
}
