/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package src

import (
	"context"

	"base-engine/auth"
	"base-engine/gen"
)

// InitializedOrganizationsIds 通过授权关系返回初始化组织 ID。
func (r *AccountResolver) InitializedOrganizationsIds(ctx context.Context, obj *gen.Account) ([]string, error) {
	items, err := r.Handlers.AccountInitializedOrganizations(ctx, r.GeneratedResolver, obj)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.ID)
	}
	return ids, nil
}

func (r *AccountResolver) OpeningRecordsIds(ctx context.Context, _ *gen.Account) ([]string, error) {
	return nil, denyOpeningRecordProjection(ctx)
}

func (r *AccountResolver) RecordedOpeningRecordsIds(ctx context.Context, _ *gen.Account) ([]string, error) {
	return nil, denyOpeningRecordProjection(ctx)
}

func (r *OrganizationResolver) OpeningRecordsIds(ctx context.Context, _ *gen.Organization) ([]string, error) {
	return nil, denyOpeningRecordProjection(ctx)
}

func denyOpeningRecordProjection(ctx context.Context) error {
	if _, err := auth.RequirePrincipal(ctx); err != nil {
		return err
	}
	return auth.NewError(auth.CodePermissionDenied)
}

// InitialAccountID 防止直接读取标量 ID 绕过关系授权。
func (r *OrganizationResolver) InitialAccountID(ctx context.Context, obj *gen.Organization) (*string, error) {
	principal, err := auth.RequirePrincipal(ctx)
	if err != nil {
		return nil, err
	}
	if principal.WorkspaceType != auth.WorkspaceTypeHeadquarters || !principal.Has("account:update") {
		return nil, auth.NewError(auth.CodePermissionDenied)
	}
	return obj.InitialAccountID, nil
}

// MembershipsIds 通过已授权关系 handler 返回账号成员关系 ID。
func (r *AccountResolver) MembershipsIds(ctx context.Context, obj *gen.Account) ([]string, error) {
	items, err := r.Handlers.AccountMemberships(ctx, r.GeneratedResolver, obj)
	return membershipIDs(items), err
}

// SessionsIds 通过已授权关系 handler 返回账号会话 ID。
func (r *AccountResolver) SessionsIds(ctx context.Context, obj *gen.Account) ([]string, error) {
	items, err := r.Handlers.AccountSessions(ctx, r.GeneratedResolver, obj)
	return sessionIDs(items), err
}

// ReviewedStoresIds 通过已授权关系 handler 返回账号审核过的门店 ID。
func (r *AccountResolver) ReviewedStoresIds(ctx context.Context, obj *gen.Account) ([]string, error) {
	items, err := r.Handlers.AccountReviewedStores(ctx, r.GeneratedResolver, obj)
	return storeIDs(items), err
}

// SentMembershipInvitationsIds 通过已授权关系 handler 返回账号发出的邀请 ID。
func (r *AccountResolver) SentMembershipInvitationsIds(ctx context.Context, obj *gen.Account) ([]string, error) {
	items, err := r.Handlers.AccountSentMembershipInvitations(ctx, r.GeneratedResolver, obj)
	return invitationIDs(items), err
}

// AuditLogsIds 通过已授权关系 handler 返回账号审计日志 ID。
func (r *AccountResolver) AuditLogsIds(ctx context.Context, obj *gen.Account) ([]string, error) {
	items, err := r.Handlers.AccountAuditLogs(ctx, r.GeneratedResolver, obj)
	return auditLogIDs(items), err
}

// MembershipsIds 通过已授权关系 handler 返回组织成员关系 ID。
func (r *OrganizationResolver) MembershipsIds(ctx context.Context, obj *gen.Organization) ([]string, error) {
	items, err := r.Handlers.OrganizationMemberships(ctx, r.GeneratedResolver, obj)
	return membershipIDs(items), err
}

// StoresIds 通过已授权关系 handler 返回组织门店 ID。
func (r *OrganizationResolver) StoresIds(ctx context.Context, obj *gen.Organization) ([]string, error) {
	items, err := r.Handlers.OrganizationStores(ctx, r.GeneratedResolver, obj)
	return storeIDs(items), err
}

// RolesIds 通过已授权关系 handler 返回组织角色 ID。
func (r *OrganizationResolver) RolesIds(ctx context.Context, obj *gen.Organization) ([]string, error) {
	items, err := r.Handlers.OrganizationRoles(ctx, r.GeneratedResolver, obj)
	return roleIDs(items), err
}

// SessionsIds 通过已授权关系 handler 返回组织会话 ID。
func (r *OrganizationResolver) SessionsIds(ctx context.Context, obj *gen.Organization) ([]string, error) {
	items, err := r.Handlers.OrganizationSessions(ctx, r.GeneratedResolver, obj)
	return sessionIDs(items), err
}

// AuditLogsIds 通过已授权关系 handler 返回组织审计日志 ID。
func (r *OrganizationResolver) AuditLogsIds(ctx context.Context, obj *gen.Organization) ([]string, error) {
	items, err := r.Handlers.OrganizationAuditLogs(ctx, r.GeneratedResolver, obj)
	return auditLogIDs(items), err
}

// RolesIds 通过已授权关系 handler 返回成员角色 ID。
func (r *OperatorMembershipResolver) RolesIds(ctx context.Context, obj *gen.OperatorMembership) ([]string, error) {
	items, err := r.Handlers.OperatorMembershipRoles(ctx, r.GeneratedResolver, obj)
	return roleIDs(items), err
}

// StoresIds 通过已授权关系 handler 返回成员门店 ID。
func (r *OperatorMembershipResolver) StoresIds(ctx context.Context, obj *gen.OperatorMembership) ([]string, error) {
	items, err := r.Handlers.OperatorMembershipStores(ctx, r.GeneratedResolver, obj)
	return storeIDs(items), err
}

// InvitationsIds 通过已授权关系 handler 返回成员邀请 ID。
func (r *OperatorMembershipResolver) InvitationsIds(ctx context.Context, obj *gen.OperatorMembership) ([]string, error) {
	items, err := r.Handlers.OperatorMembershipInvitations(ctx, r.GeneratedResolver, obj)
	return invitationIDs(items), err
}

// RolesIds 通过已授权关系 handler 返回权限关联角色 ID。
func (r *PermissionResolver) RolesIds(ctx context.Context, obj *gen.Permission) ([]string, error) {
	items, err := r.Handlers.PermissionRoles(ctx, r.GeneratedResolver, obj)
	return roleIDs(items), err
}

// MembersIds 通过已授权关系 handler 返回角色成员 ID。
func (r *OperatorRoleResolver) MembersIds(ctx context.Context, obj *gen.OperatorRole) ([]string, error) {
	items, err := r.Handlers.OperatorRoleMembers(ctx, r.GeneratedResolver, obj)
	return membershipIDs(items), err
}

// PermissionsIds 通过已授权关系 handler 返回角色权限 ID。
func (r *OperatorRoleResolver) PermissionsIds(ctx context.Context, obj *gen.OperatorRole) ([]string, error) {
	items, err := r.Handlers.OperatorRolePermissions(ctx, r.GeneratedResolver, obj)
	return permissionIDs(items), err
}

// MembersIds 通过已授权关系 handler 返回门店成员 ID。
func (r *StoreResolver) MembersIds(ctx context.Context, obj *gen.Store) ([]string, error) {
	items, err := r.Handlers.StoreMembers(ctx, r.GeneratedResolver, obj)
	return membershipIDs(items), err
}

// AuditLogsIds 通过已授权关系 handler 返回门店审计日志 ID。
func (r *StoreResolver) AuditLogsIds(ctx context.Context, obj *gen.Store) ([]string, error) {
	items, err := r.Handlers.StoreAuditLogs(ctx, r.GeneratedResolver, obj)
	return auditLogIDs(items), err
}

// AuditLogsIds 通过已授权关系 handler 返回会话审计日志 ID。
func (r *SessionResolver) AuditLogsIds(ctx context.Context, obj *gen.Session) ([]string, error) {
	items, err := r.Handlers.SessionAuditLogs(ctx, r.GeneratedResolver, obj)
	return auditLogIDs(items), err
}

func membershipIDs(items []*gen.OperatorMembership) []string {
	ids := make([]string, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.ID)
	}
	return ids
}

func storeIDs(items []*gen.Store) []string {
	ids := make([]string, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.ID)
	}
	return ids
}

func roleIDs(items []*gen.OperatorRole) []string {
	ids := make([]string, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.ID)
	}
	return ids
}

func permissionIDs(items []*gen.Permission) []string {
	ids := make([]string, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.ID)
	}
	return ids
}

func sessionIDs(items []*gen.Session) []string {
	ids := make([]string, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.ID)
	}
	return ids
}

func invitationIDs(items []*gen.MembershipInvitation) []string {
	ids := make([]string, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.ID)
	}
	return ids
}

func auditLogIDs(items []*gen.AuditLog) []string {
	ids := make([]string, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.ID)
	}
	return ids
}
