package src

import (
	"context"
	"testing"

	"base-engine/auth"
)

func TestHQMembershipDirectiveUsesWorkspacePermission(t *testing.T) {
	organizationID := "hq-organization"
	principal := &auth.WorkspacePrincipal{
		AccountID:      "hq-admin",
		WorkspaceType:  auth.WorkspaceTypeHeadquarters,
		OrganizationID: &organizationID,
		Permissions: map[string]struct{}{
			"hqMembership:create": {},
		},
	}
	ctx := auth.WithPrincipal(context.Background(), principal)
	result, err := hasPermissionDirective(ctx, nil, func(context.Context) (any, error) {
		return "allowed", nil
	}, "operatorMembership:create")
	if err != nil || result != "allowed" {
		t.Fatalf("HQ membership directive result=%v err=%v", result, err)
	}
}
