/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-22
 */
package authorization

import (
	"context"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
)

func TestLastHQSuperAdminRoleCannotBeRemoved(t *testing.T) {
	fixture := newPrincipalFixture(t)
	membership := &gen.OperatorMembership{}
	if err := fixture.db.First(membership, "id = ?", "hq-membership").Error; err != nil {
		t.Fatal(err)
	}
	custom := gen.OperatorRole{ID: "hq-custom-role", Name: "Custom", Kind: gen.RoleKindCustom, OrganizationID: "hq-org"}
	if err := fixture.db.Create(&custom).Error; err != nil {
		t.Fatal(err)
	}
	resolver := &gen.GeneratedResolver{DB: gen.NewDB(fixture.db)}
	err := ensureOwnerRoleRemains(context.Background(), resolver, membership, map[string]interface{}{"rolesIds": []string{custom.ID}})
	if auth.ErrorCode(err) != auth.CodeLastHQSuperAdminRequired {
		t.Fatalf("last HQ super administrator error = %v", err)
	}
}

func TestSuspendedHQSuperAdminIsNotPrivilegedActor(t *testing.T) {
	fixture := newPrincipalFixture(t)
	database := fixture.db
	organizationID := "hq-org"
	principal := &auth.WorkspacePrincipal{AccountID: "suspended-super-account", OrganizationID: &organizationID, MembershipID: authorizationStringPointer("suspended-super-membership"), WorkspaceType: auth.WorkspaceTypeHeadquarters}
	values := []any{
		&gen.Account{ID: principal.AccountID, Phone: "13900000031", DisplayName: "Suspended", Status: gen.AccountStatusActive},
		&gen.OperatorMembership{ID: "suspended-super-membership", AccountID: principal.AccountID, OrganizationID: *principal.OrganizationID, Status: gen.MembershipStatusSuspended},
		&gen.OperatorRole{ID: "suspended-super-role", Name: "Suspended Super", Kind: gen.RoleKindHqSuperAdmin, OrganizationID: *principal.OrganizationID},
		&membershipRole{OperatorMembershipID: "suspended-super-membership", OperatorRoleID: "suspended-super-role"},
	}
	for _, value := range values {
		if err := database.Create(value).Error; err != nil {
			t.Fatal(err)
		}
	}
	super, err := IsHQSuperAdmin(database, principal)
	if err != nil {
		t.Fatal(err)
	}
	if super {
		t.Fatal("suspended HQ membership was treated as an active super administrator")
	}
}

func TestRoleUsedOnlyBySuspendedMemberCanBeDeleted(t *testing.T) {
	fixture := newPrincipalFixture(t)
	database := fixture.db
	organizationID := "hq-org"
	principal := &auth.WorkspacePrincipal{AccountID: "suspended-member-account", OrganizationID: &organizationID, WorkspaceType: auth.WorkspaceTypeHeadquarters}
	role := &gen.OperatorRole{ID: "suspended-member-role", Name: "Unused", Kind: gen.RoleKindCustom, OrganizationID: *principal.OrganizationID}
	membership := &gen.OperatorMembership{ID: "suspended-member", AccountID: principal.AccountID, OrganizationID: *principal.OrganizationID, Status: gen.MembershipStatusSuspended}
	for _, value := range []any{&gen.Account{ID: principal.AccountID, Phone: "13900000032", DisplayName: "Suspended", Status: gen.AccountStatusActive}, role, membership, &membershipRole{OperatorMembershipID: membership.ID, OperatorRoleID: role.ID}} {
		if err := database.Create(value).Error; err != nil {
			t.Fatal(err)
		}
	}
	if err := EnsureRoleUnused(database, role); err != nil {
		t.Fatalf("suspended membership kept role in use: %v", err)
	}
}

func authorizationStringPointer(value string) *string { return &value }
