package ai

import (
	"encoding/json"
	"errors"

	"base-engine/auth"
)

func (s *ApprovalService) Issue(plan Plan) (string, error) {
	if s == nil || plan.ID == "" || len(plan.Steps) == 0 || !s.now().Before(plan.ExpiresAt) {
		return "", errors.New("INVALID_PLAN")
	}
	return s.store.Issue(plan.ID, plan)
}

func (s *ApprovalService) Consume(token string, principal *auth.WorkspacePrincipal) (Plan, error) {
	return s.store.Consume(token, principalBinding(principal), func(plan Plan) bool {
		return planMatchesPrincipal(plan, principal) && s.planPermissionsCurrent(plan, principal)
	})
}

func (s *ApprovalService) planPermissionsCurrent(plan Plan, principal *auth.WorkspacePrincipal) bool {
	for _, step := range plan.Steps {
		spec, ok := s.catalog.Lookup(step.ToolID)
		if !ok || !specAllowed(spec, principal, PhaseRun, map[string]struct{}{step.ToolID: {}}) {
			return false
		}
	}
	return true
}

func planBinding(plan Plan) string {
	encoded, _ := json.Marshal([]string{plan.AccountID, plan.SessionID, plan.WorkspaceType, plan.OrganizationID})
	return string(encoded)
}

func principalBinding(principal *auth.WorkspacePrincipal) string {
	if principal == nil {
		return ""
	}
	organizationID := ""
	if principal.OrganizationID != nil {
		organizationID = *principal.OrganizationID
	}
	encoded, _ := json.Marshal([]string{principal.AccountID, principal.SessionID, string(principal.WorkspaceType), organizationID})
	return string(encoded)
}

func planMatchesPrincipal(plan Plan, principal *auth.WorkspacePrincipal) bool {
	return principal != nil && planBinding(plan) == principalBinding(principal)
}
