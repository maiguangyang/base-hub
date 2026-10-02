/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package authorization

import (
	"testing"

	"base-engine/auth"
)

// TestAuthorizeWorkspaceAndScope 验证发现区、组织边界与门店边界。
func TestAuthorizeWorkspaceAndScope(t *testing.T) {
	organizationID := "org-a"
	storeID := "store-a"
	principal := &auth.WorkspacePrincipal{
		WorkspaceType: auth.WorkspaceTypeFranchise, OrganizationID: &organizationID,
		Permissions: map[string]struct{}{"store:update": {}}, StoreIDs: map[string]struct{}{storeID: {}},
	}
	if err := Authorize(principal, Intent{Action: "store:update", Mode: AccessUpdate, ResourceOrganizationID: &organizationID, StoreID: &storeID}); err != nil {
		t.Fatal(err)
	}
	otherOrganization := "org-b"
	if code := auth.ErrorCode(Authorize(principal, Intent{Action: "store:update", Mode: AccessUpdate, ResourceOrganizationID: &otherOrganization})); code != auth.CodePermissionDenied {
		t.Fatalf("cross organization code = %s", code)
	}
	otherStore := "store-b"
	if code := auth.ErrorCode(Authorize(principal, Intent{Action: "store:update", Mode: AccessUpdate, ResourceOrganizationID: &organizationID, StoreID: &otherStore})); code != auth.CodeStoreScopeDenied {
		t.Fatalf("cross store code = %s", code)
	}
}

// TestAuthorizeWorkspaceRestrictions 验证发现区与总部不能执行加盟 generated 写入。
func TestAuthorizeWorkspaceRestrictions(t *testing.T) {
	discovery := &auth.WorkspacePrincipal{WorkspaceType: auth.WorkspaceTypeDiscovery}
	if auth.ErrorCode(Authorize(discovery, Intent{Action: "store:read", Mode: AccessRead})) != auth.CodeWorkspaceForbidden {
		t.Fatal("discovery workspace accessed business data")
	}
	hq := &auth.WorkspacePrincipal{WorkspaceType: auth.WorkspaceTypeHeadquarters, Permissions: map[string]struct{}{"store:update": {}}}
	if auth.ErrorCode(Authorize(hq, Intent{Action: "store:update", Mode: AccessUpdate})) != auth.CodeWorkspaceForbidden {
		t.Fatal("HQ used franchise generated write")
	}
	hqOrganization := "org-hq"
	franchiseOrganization := "org-franchise"
	hq = &auth.WorkspacePrincipal{
		WorkspaceType: auth.WorkspaceTypeHeadquarters, OrganizationID: &hqOrganization,
		Permissions: map[string]struct{}{"hqMembership:update": {}},
	}
	if err := Authorize(hq, Intent{Action: "hqMembership:update", Mode: AccessUpdate, ResourceOrganizationID: &hqOrganization}); err != nil {
		t.Fatalf("HQ own-organization write denied: %v", err)
	}
	if auth.ErrorCode(Authorize(hq, Intent{Action: "hqMembership:update", Mode: AccessUpdate, ResourceOrganizationID: &franchiseOrganization})) != auth.CodePermissionDenied {
		t.Fatal("HQ changed a franchise membership through generated CRUD")
	}
}
