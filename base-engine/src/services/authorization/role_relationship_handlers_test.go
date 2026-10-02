package authorization

import (
	"context"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
)

func TestRolePermissionsMatchEffectiveRolePermissions(t *testing.T) {
	fixture := newPrincipalFixture(t)
	resolver := &gen.GeneratedResolver{DB: gen.NewDB(fixture.db)}

	t.Run("HQ super administrator resolves all system permissions", func(t *testing.T) {
		permissions := resolveRolePermissionsForTest(t, fixture, resolver, "hq-session", "hq-role")
		assertPermissionActions(t, permissions, []string{"account:create", "store:approve"}, []string{"store:submit"})
	})

	t.Run("franchise owner resolves all tenant permissions", func(t *testing.T) {
		permissions := resolveRolePermissionsForTest(t, fixture, resolver, "owner-session", "owner-role")
		assertPermissionActions(t, permissions, []string{"store:create", "store:submit"}, []string{"account:create"})
	})

	t.Run("custom role resolves only its active associations", func(t *testing.T) {
		permissions := resolveRolePermissionsForTest(t, fixture, resolver, "owner-session", "custom-role")
		assertPermissionActions(t, permissions, []string{"store:read"}, []string{"store:create", "account:read"})
	})
}

func resolveRolePermissionsForTest(t *testing.T, fixture *principalFixture, resolver *gen.GeneratedResolver, sessionID, roleID string) []*gen.Permission {
	t.Helper()
	ctx := auth.WithPrincipal(context.Background(), fixture.resolve(t, sessionID))
	permissions, err := rolePermissions(ctx, resolver, &gen.OperatorRole{ID: roleID})
	if err != nil {
		t.Fatal(err)
	}
	return permissions
}

func assertPermissionActions(t *testing.T, permissions []*gen.Permission, included, excluded []string) {
	t.Helper()
	actions := make(map[string]struct{}, len(permissions))
	for _, permission := range permissions {
		actions[permission.Action] = struct{}{}
	}
	for _, action := range included {
		if _, ok := actions[action]; !ok {
			t.Errorf("missing permission %q", action)
		}
	}
	for _, action := range excluded {
		if _, ok := actions[action]; ok {
			t.Errorf("unexpected permission %q", action)
		}
	}
}
