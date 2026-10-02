/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package src

import (
	"context"
	"sort"
	"time"

	"base-engine/auth"
	"base-engine/gen"
)

// Viewer 返回当前数据库权威身份视图。
func (r *QueryResolver) Viewer(ctx context.Context) (*gen.Viewer, error) {
	principal, err := auth.RequirePrincipal(ctx)
	if err != nil {
		return nil, err
	}
	return buildViewer(ctx, r.Services, principal)
}

// Workspaces 返回账号可进入的有效组织工作台。
func (r *QueryResolver) Workspaces(ctx context.Context) ([]*gen.Workspace, error) {
	principal, err := auth.RequirePrincipal(ctx)
	if err != nil {
		return nil, err
	}
	return loadWorkspaces(ctx, r.Services, principal.AccountID)
}

// PendingMembershipInvitations 返回当前账号仍可接受的邀请。
func (r *QueryResolver) PendingMembershipInvitations(ctx context.Context) ([]*gen.MembershipInvitation, error) {
	principal, err := auth.RequirePrincipal(ctx)
	if err != nil {
		return nil, err
	}
	var membershipIDs []string
	err = r.Services.DB.Model(&gen.OperatorMembership{}).
		Where("account_id = ? AND status = ?", principal.AccountID, gen.MembershipStatusInvited).
		Where("is_delete IS NULL OR is_delete = ?", 1).
		Pluck("id", &membershipIDs).Error
	if err != nil || len(membershipIDs) == 0 {
		return []*gen.MembershipInvitation{}, err
	}
	var invitations []*gen.MembershipInvitation
	err = r.Services.DB.WithContext(ctx).
		Where("membership_id IN ? AND accepted_at IS NULL AND revoked_at IS NULL AND expires_at > ?", membershipIDs, time.Now()).
		Where("is_delete IS NULL OR is_delete = ?", 1).
		Find(&invitations).Error
	return invitations, err
}

func buildViewer(ctx context.Context, services Dependencies, principal *auth.WorkspacePrincipal) (*gen.Viewer, error) {
	account := &gen.Account{}
	if err := services.DB.WithContext(ctx).Where("is_delete IS NULL OR is_delete = ?", 1).First(account, "id = ?", principal.AccountID).Error; err != nil {
		return nil, err
	}
	workspaces, err := loadWorkspaces(ctx, services, principal.AccountID)
	if err != nil {
		return nil, err
	}
	permissions := make([]string, 0, len(principal.Permissions))
	for action := range principal.Permissions {
		permissions = append(permissions, action)
	}
	sort.Strings(permissions)
	return &gen.Viewer{
		Account: account, CurrentWorkspace: currentWorkspace(workspaces, principal),
		Workspaces: workspaces, Permissions: permissions,
	}, nil
}

func loadWorkspaces(ctx context.Context, services Dependencies, accountID string) ([]*gen.Workspace, error) {
	workspaces := []*gen.Workspace{{WorkspaceType: gen.WorkspaceTypeDiscovery, OrganizationName: "workspace.discovery", HomePath: "/admin"}}
	var memberships []gen.OperatorMembership
	err := services.DB.WithContext(ctx).Where("account_id = ? AND status = ?", accountID, gen.MembershipStatusActive).
		Where("is_delete IS NULL OR is_delete = ?", 1).Find(&memberships).Error
	if err != nil {
		return nil, err
	}
	for _, membership := range memberships {
		organization := &gen.Organization{}
		if err := services.DB.WithContext(ctx).Where("is_delete IS NULL OR is_delete = ?", 1).First(organization, "id = ? AND status = ?", membership.OrganizationID, gen.OrganizationStatusActive).Error; err != nil {
			continue
		}
		workspaces = append(workspaces, organizationWorkspace(organization))
	}
	return workspaces, nil
}

func organizationWorkspace(organization *gen.Organization) *gen.Workspace {
	workspaceType := gen.WorkspaceTypeFranchise
	homePath := "/admin/franchise"
	if organization.Type == gen.OrganizationTypeHeadquarters {
		workspaceType = gen.WorkspaceTypeHeadquarters
		homePath = "/admin/hq"
	}
	return &gen.Workspace{
		WorkspaceType: workspaceType, OrganizationID: &organization.ID,
		OrganizationName: organization.Name, HomePath: homePath,
	}
}

func currentWorkspace(workspaces []*gen.Workspace, principal *auth.WorkspacePrincipal) *gen.Workspace {
	for _, workspace := range workspaces {
		if auth.WorkspaceType(workspace.WorkspaceType) != principal.WorkspaceType {
			continue
		}
		if sameOptionalID(workspace.OrganizationID, principal.OrganizationID) {
			return workspace
		}
	}
	return nil
}

func sameOptionalID(left, right *string) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return *left == *right
}

func selectWorkspaceRecords(ctx context.Context, services Dependencies, principal *auth.WorkspacePrincipal, input gen.SelectWorkspaceInput) (*gen.Account, *gen.Session, error) {
	account := &gen.Account{}
	if err := services.DB.WithContext(ctx).Where("is_delete IS NULL OR is_delete = ?", 1).First(account, "id = ?", principal.AccountID).Error; err != nil {
		return nil, nil, auth.NewError(auth.CodeAuthRequired)
	}
	if account.MustChangePassword && input.WorkspaceType != gen.WorkspaceTypeDiscovery {
		return nil, nil, auth.NewError(auth.CodeTempPasswordChangeRequired)
	}
	organizationID, err := validateWorkspaceSelection(ctx, services, account.ID, input)
	if err != nil {
		return nil, nil, err
	}
	session := &gen.Session{}
	updates := map[string]any{"workspace_type": input.WorkspaceType, "organization_id": organizationID, "last_seen_at": time.Now()}
	result := services.DB.WithContext(ctx).Model(session).
		Where("id = ? AND account_id = ? AND revoked_at IS NULL", principal.SessionID, account.ID).
		Where("is_delete IS NULL OR is_delete = ?", 1).Updates(updates)
	if result.Error != nil || result.RowsAffected == 0 {
		return nil, nil, auth.NewError(auth.CodeSessionRevoked)
	}
	if err := services.DB.WithContext(ctx).Where("is_delete IS NULL OR is_delete = ?", 1).First(session, "id = ?", principal.SessionID).Error; err != nil {
		return nil, nil, err
	}
	return account, session, nil
}

func validateWorkspaceSelection(ctx context.Context, services Dependencies, accountID string, input gen.SelectWorkspaceInput) (*string, error) {
	if input.WorkspaceType == gen.WorkspaceTypeDiscovery {
		return nil, nil
	}
	if input.OrganizationID == nil {
		return nil, auth.NewError(auth.CodeWorkspaceForbidden)
	}
	organization := &gen.Organization{}
	if err := services.DB.WithContext(ctx).Where("is_delete IS NULL OR is_delete = ?", 1).
		First(organization, "id = ? AND status = ?", *input.OrganizationID, gen.OrganizationStatusActive).Error; err != nil {
		return nil, auth.NewError(auth.CodeWorkspaceForbidden)
	}
	expected := gen.OrganizationTypeFranchise
	if input.WorkspaceType == gen.WorkspaceTypeHeadquarters {
		expected = gen.OrganizationTypeHeadquarters
	}
	if organization.Type != expected {
		return nil, auth.NewError(auth.CodeWorkspaceForbidden)
	}
	var count int64
	err := services.DB.WithContext(ctx).Model(&gen.OperatorMembership{}).
		Where("account_id = ? AND organization_id = ? AND status = ?", accountID, organization.ID, gen.MembershipStatusActive).
		Where("is_delete IS NULL OR is_delete = ?", 1).Count(&count).Error
	if err != nil || count != 1 {
		return nil, auth.NewError(auth.CodeWorkspaceForbidden)
	}
	return &organization.ID, nil
}
