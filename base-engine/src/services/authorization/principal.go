/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package authorization

import (
	"context"
	"errors"
	"time"

	"base-engine/auth"
	"base-engine/config"
	"base-engine/gen"
	"gorm.io/gorm"
)

// PrincipalResolver 从签名令牌定位并解析数据库权威身份。
type PrincipalResolver struct {
	db     *gorm.DB
	config config.SecurityConfig
}

// NewPrincipalResolver 创建请求身份解析器。
func NewPrincipalResolver(db *gorm.DB, cfg config.SecurityConfig) *PrincipalResolver {
	return &PrincipalResolver{db: db, config: cfg}
}

// Resolve 校验令牌、会话和当前授权状态。
func (r *PrincipalResolver) Resolve(ctx context.Context, token string, now time.Time) (*auth.WorkspacePrincipal, error) {
	claims, err := auth.ParseSessionClaims(r.config, token)
	if err != nil {
		return nil, auth.NewError(auth.CodeAuthRequired)
	}
	session, err := r.loadSession(ctx, claims.SessionID, now)
	if err != nil {
		return nil, err
	}
	account, err := r.loadAccount(ctx, session, claims)
	if err != nil {
		return nil, err
	}
	principal := newWorkspacePrincipal(session)
	if session.WorkspaceType == gen.WorkspaceTypeDiscovery {
		return principal, nil
	}
	if err := r.resolveOrganizationWorkspace(ctx, principal, session, account); err != nil {
		return nil, err
	}
	return principal, nil
}

func (r *PrincipalResolver) loadSession(ctx context.Context, sessionID string, now time.Time) (*gen.Session, error) {
	session := &gen.Session{}
	err := activeRecordCondition(r.db.WithContext(ctx)).First(session, "id = ?", sessionID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, auth.NewError(auth.CodeSessionRevoked)
	}
	if err != nil {
		return nil, err
	}
	if session.RevokedAt != nil || !session.ExpiresAt.After(now) {
		return nil, auth.NewError(auth.CodeSessionRevoked)
	}
	return session, nil
}

func (r *PrincipalResolver) loadAccount(ctx context.Context, session *gen.Session, claims *auth.SessionClaims) (*gen.Account, error) {
	if claims.AccountID != session.AccountID || claims.WorkspaceType != auth.WorkspaceType(session.WorkspaceType) {
		return nil, auth.NewError(auth.CodeAuthRequired)
	}
	account := &gen.Account{}
	if err := activeRecordCondition(r.db.WithContext(ctx)).First(account, "id = ?", session.AccountID).Error; err != nil {
		return nil, auth.NewError(auth.CodeAuthRequired)
	}
	if account.Status != gen.AccountStatusActive {
		return nil, auth.NewError(auth.CodeAuthRequired)
	}
	if claims.CredentialVersion != account.CredentialVersion || session.CredentialVersion != account.CredentialVersion {
		return nil, auth.NewError(auth.CodeCredentialsChanged)
	}
	return account, nil
}

func newWorkspacePrincipal(session *gen.Session) *auth.WorkspacePrincipal {
	return &auth.WorkspacePrincipal{
		AccountID: session.AccountID, SessionID: session.ID,
		WorkspaceType:  auth.WorkspaceType(session.WorkspaceType),
		OrganizationID: session.OrganizationID,
		Permissions:    make(map[string]struct{}), StoreIDs: make(map[string]struct{}),
	}
}

func (r *PrincipalResolver) resolveOrganizationWorkspace(ctx context.Context, principal *auth.WorkspacePrincipal, session *gen.Session, account *gen.Account) error {
	if session.OrganizationID == nil {
		return auth.NewError(auth.CodeWorkspaceForbidden)
	}
	organization := &gen.Organization{}
	if err := activeRecordCondition(r.db.WithContext(ctx)).First(organization, "id = ?", *session.OrganizationID).Error; err != nil {
		return auth.NewError(auth.CodeWorkspaceForbidden)
	}
	if organization.Status != gen.OrganizationStatusActive {
		return auth.NewError(auth.CodeOrganizationSuspended)
	}
	if !workspaceMatchesOrganization(session.WorkspaceType, organization.Type) {
		return auth.NewError(auth.CodeWorkspaceForbidden)
	}
	membership, err := r.loadMembership(ctx, account.ID, organization.ID)
	if err != nil {
		return err
	}
	principal.MembershipID = &membership.ID
	if err := r.resolvePermissions(ctx, principal, membership); err != nil {
		return err
	}
	return r.resolveStores(ctx, principal, membership, organization.ID)
}

func (r *PrincipalResolver) loadMembership(ctx context.Context, accountID, organizationID string) (*gen.OperatorMembership, error) {
	membership := &gen.OperatorMembership{}
	err := activeRecordCondition(r.db.WithContext(ctx)).Where("account_id = ? AND organization_id = ?", accountID, organizationID).First(membership).Error
	if err != nil || membership.Status != gen.MembershipStatusActive {
		return nil, auth.NewError(auth.CodeMembershipInactive)
	}
	return membership, nil
}

func (r *PrincipalResolver) resolvePermissions(ctx context.Context, principal *auth.WorkspacePrincipal, membership *gen.OperatorMembership) error {
	var roles []gen.OperatorRole
	database := r.db.WithContext(ctx).Where("operator_roles.is_delete IS NULL OR operator_roles.is_delete = ?", 1)
	if err := database.Model(membership).Association("Roles").Find(&roles); err != nil {
		return err
	}
	for index := range roles {
		if roles[index].OrganizationID != membership.OrganizationID {
			continue
		}
		if err := r.addRolePermissions(ctx, principal, &roles[index]); err != nil {
			return err
		}
	}
	return nil
}

func (r *PrincipalResolver) addRolePermissions(ctx context.Context, principal *auth.WorkspacePrincipal, role *gen.OperatorRole) error {
	var permissions []gen.Permission
	database := r.db.WithContext(ctx).Where("permissions.is_delete IS NULL OR permissions.is_delete = ?", 1)
	if err := database.Model(role).Association("Permissions").Find(&permissions); err != nil {
		return err
	}
	expectedScope := workspacePermissionScope(principal.WorkspaceType)
	for _, permission := range permissions {
		if permission.Scope == expectedScope {
			principal.Permissions[permission.Action] = struct{}{}
		}
	}
	scope := dynamicPermissionScope(role.Kind, principal.WorkspaceType)
	if scope == "" {
		return nil
	}
	permissions = nil
	if err := activeRecordCondition(r.db.WithContext(ctx)).Where("scope = ?", scope).Find(&permissions).Error; err != nil {
		return err
	}
	for _, permission := range permissions {
		principal.Permissions[permission.Action] = struct{}{}
	}
	return nil
}

func dynamicPermissionScope(kind gen.RoleKind, workspace auth.WorkspaceType) gen.PermissionScope {
	if kind == gen.RoleKindHqSuperAdmin && workspace == auth.WorkspaceTypeHeadquarters {
		return gen.PermissionScopeSystem
	}
	if kind == gen.RoleKindFranchiseOwner && workspace == auth.WorkspaceTypeFranchise {
		return gen.PermissionScopeTenant
	}
	return ""
}

func workspacePermissionScope(workspace auth.WorkspaceType) gen.PermissionScope {
	if workspace == auth.WorkspaceTypeHeadquarters {
		return gen.PermissionScopeSystem
	}
	if workspace == auth.WorkspaceTypeFranchise {
		return gen.PermissionScopeTenant
	}
	return ""
}

func workspaceMatchesOrganization(workspace gen.WorkspaceType, organization gen.OrganizationType) bool {
	if workspace == gen.WorkspaceTypeHeadquarters {
		return organization == gen.OrganizationTypeHeadquarters
	}
	if workspace == gen.WorkspaceTypeFranchise {
		return organization == gen.OrganizationTypeFranchise
	}
	return false
}

func (r *PrincipalResolver) resolveStores(ctx context.Context, principal *auth.WorkspacePrincipal, membership *gen.OperatorMembership, organizationID string) error {
	var stores []gen.Store
	if membership.StoreAccessMode == gen.StoreAccessModeAllStores {
		principal.AllStores = true
		if err := activeRecordCondition(r.db.WithContext(ctx)).Where("organization_id = ? AND lifecycle = ?", organizationID, gen.StoreLifecycleActive).Find(&stores).Error; err != nil {
			return err
		}
	} else if err := r.db.WithContext(ctx).Where("stores.is_delete IS NULL OR stores.is_delete = ?", 1).Model(membership).Association("Stores").Find(&stores); err != nil {
		return err
	}
	for _, store := range stores {
		if store.OrganizationID == organizationID && store.Lifecycle == gen.StoreLifecycleActive {
			principal.StoreIDs[store.ID] = struct{}{}
		}
	}
	return nil
}
