package tools

import (
	"net/http"

	"base-engine/auth"
	"base-engine/src/services/ai"
)

func paymentConfigSpecs() []ai.ToolSpec {
	read := paymentSpec("HqPaymentConfigRead", "http.paymentConfig.read", http.MethodGet, "/api/payment-config", ai.ModeReadOnly, "paymentConfig:read", "LOW")
	state := []ai.ToolSpec{
		paymentWriteSpec("HqGlobalPaymentState", "http.paymentConfig.state", "/api/payment-config/state", "channel"),
		paymentWriteSpec("HqFranchisePaymentState", "http.paymentConfig.state", "/api/payment-config/state", "organizationId"),
		paymentWriteSpec("HqStorePaymentState", "http.paymentConfig.state", "/api/payment-config/state", "storeId"),
	}
	restore := []ai.ToolSpec{
		paymentWriteSpec("HqFranchisePaymentRestore", "http.paymentConfig.restoreInheritance", "/api/payment-config/restore-inheritance", "organizationId"),
		paymentWriteSpec("HqStorePaymentRestore", "http.paymentConfig.restoreInheritance", "/api/payment-config/restore-inheritance", "storeId"),
	}
	return append(append([]ai.ToolSpec{read}, state...), restore...)
}

func paymentSpec(id, operation, method, path string, mode ai.ToolMode, permission, risk string) ai.ToolSpec {
	title, description := toolMetadata(id, mode, auth.WorkspaceTypeHeadquarters)
	return ai.ToolSpec{ID: id, Name: id, Title: title, Description: description, OperationID: operation,
		Method: method, Path: path, Mode: mode, Risk: risk, Permission: permission,
		Workspaces: []auth.WorkspaceType{auth.WorkspaceTypeHeadquarters}, OutputFields: []string{"channels", "validationOnly"}}
}

func paymentWriteSpec(id, operation, path, target string) ai.ToolSpec {
	spec := paymentSpec(id, operation, http.MethodPost, path, ai.ModeWrite, "paymentConfig:manage", "HIGH")
	spec.AdditionalPermissions = []string{"paymentConfig:read"}
	spec.WriteKind, spec.TargetFields = ai.WriteExisting, []string{target}
	spec.ArgumentFields = []string{"scope", "organizationId", "storeId", "channel", "state", "recordId", "version"}
	return spec
}
