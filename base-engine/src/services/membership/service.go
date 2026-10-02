/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package membership

import (
	"context"
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"base-engine/src/services/authentication"
	"base-engine/src/services/authorization"
	sessionservice "base-engine/src/services/session"
	"gorm.io/gorm"
)

// InviteInput 描述加盟组织邀请操作员及其初始范围。
type InviteInput struct {
	Phone           string
	DisplayName     string
	Email           *string
	RoleIDs         []string
	StoreAccessMode gen.StoreAccessMode
	StoreIDs        []string
}

// InviteResult 返回成员与只对新账号出现一次的临时密码。
type InviteResult struct {
	Membership        *gen.OperatorMembership
	TemporaryPassword *string
	InvitationPending bool
}

// Service 管理成员邀请、接受与状态不变量。
type Service struct {
	db        *gorm.DB
	audit     *audit.Service
	sessions  *sessionservice.Service
	publisher sessionservice.Publisher
}

// ServiceDependencies supplies session revocation side effects for membership governance.
type ServiceDependencies struct {
	Sessions  *sessionservice.Service
	Publisher sessionservice.Publisher
}

// NewService 创建成员服务。
func NewService(db *gorm.DB, auditService *audit.Service, dependencies ...ServiceDependencies) *Service {
	service := &Service{db: db, audit: auditService}
	if len(dependencies) > 0 {
		service.sessions = dependencies[0].Sessions
		service.publisher = dependencies[0].Publisher
	}
	return service
}

// InviteOperator 原子创建账号或邀请既有账号进入当前组织。
func (s *Service) InviteOperator(ctx context.Context, principal *auth.WorkspacePrincipal, input InviteInput) (*InviteResult, error) {
	organizationID, err := authorizeMembershipChange(principal, "operatorMembership:create", authorization.AccessCreate)
	if err != nil {
		auth.LogAuthorizationDenied(principal, inviteAuditAction(principal), "operatorMembership", "", err)
		return nil, err
	}
	input, err = normalizeInviteInput(input)
	if err != nil {
		return nil, err
	}
	input = normalizeWorkspaceInviteInput(principal, input)
	var result *InviteResult
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if txErr := validateAccessScope(tx, principal, organizationID, input); txErr != nil {
			auth.LogAuthorizationDenied(principal, inviteAuditAction(principal), "operatorMembership", "", txErr)
			return txErr
		}
		created, txErr := s.inviteTransaction(tx, principal, organizationID, input, time.Now())
		result = created
		return txErr
	})
	return result, err
}

func normalizeWorkspaceInviteInput(principal *auth.WorkspacePrincipal, input InviteInput) InviteInput {
	if principal != nil && principal.WorkspaceType == auth.WorkspaceTypeHeadquarters {
		input.StoreAccessMode = gen.StoreAccessModeAllStores
		input.StoreIDs = nil
	}
	return input
}

func (s *Service) inviteTransaction(tx *gorm.DB, principal *auth.WorkspacePrincipal, organizationID string, input InviteInput, now time.Time) (*InviteResult, error) {
	account, existing, password, err := findOrCreateAccount(tx, input, now)
	if err != nil {
		return nil, err
	}
	membership, invitation, err := pendingMembershipForReissue(tx, account.ID, organizationID, now)
	if err != nil {
		return nil, err
	}
	if membership == nil {
		membership = buildMembership(account.ID, organizationID, input.StoreAccessMode, existing, now)
		if err := tx.Create(membership).Error; err != nil {
			return nil, err
		}
	} else if err := resetPendingMembership(tx, membership, input.StoreAccessMode, now); err != nil {
		return nil, err
	}
	if err := replaceMembershipScope(tx, membership, input.RoleIDs, input.StoreIDs); err != nil {
		return nil, err
	}
	if existing {
		if err := refreshMembershipInvitation(tx, invitation, membership.ID, principal.AccountID, now); err != nil {
			return nil, err
		}
	}
	metadata := audit.Metadata{RoleIDs: input.RoleIDs, StoreIDs: input.StoreIDs}
	if err := s.writeAudit(tx, principal, membership, inviteAuditAction(principal), metadata); err != nil {
		return nil, err
	}
	return &InviteResult{Membership: membership, TemporaryPassword: password, InvitationPending: existing}, nil
}

func inviteAuditAction(principal *auth.WorkspacePrincipal) string {
	if principal != nil && principal.WorkspaceType == auth.WorkspaceTypeHeadquarters {
		return "hqAdministrator:invite"
	}
	return "operator:invite"
}

func findOrCreateAccount(tx *gorm.DB, input InviteInput, now time.Time) (*gen.Account, bool, *string, error) {
	account := &gen.Account{}
	err := inviteAccountQuery(tx, input.Phone).First(account).Error
	if err == nil {
		if account.Status != gen.AccountStatusActive || !recordActive(account.IsDelete) {
			return nil, false, nil, auth.NewError(auth.CodePermissionDenied)
		}
		return account, true, nil, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil, err
	}
	password, err := auth.GenerateTemporaryPassword()
	if err != nil {
		return nil, false, nil, err
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		return nil, false, nil, err
	}
	account = &gen.Account{
		ID: uuid.Must(uuid.NewV4()).String(), Phone: input.Phone,
		DisplayName: input.DisplayName, Email: input.Email,
		Status: gen.AccountStatusActive, MustChangePassword: true, CredentialVersion: 1,
	}
	if err := tx.Create(account).Error; err != nil {
		return nil, false, nil, err
	}
	credential := authentication.AccountCredential{
		AccountID: account.ID, PasswordHash: hash, TemporaryPasswordExpiresAt: authentication.NewTemporaryPasswordExpiry(now),
		PasswordChangedAt: now, UpdatedAt: now,
	}
	if err := tx.Create(&credential).Error; err != nil {
		return nil, false, nil, err
	}
	return account, false, &password, nil
}

func buildMembership(accountID, organizationID string, mode gen.StoreAccessMode, invited bool, now time.Time) *gen.OperatorMembership {
	status := gen.MembershipStatusActive
	acceptedAt := &now
	if invited {
		status = gen.MembershipStatusInvited
		acceptedAt = nil
	}
	return &gen.OperatorMembership{
		ID: uuid.Must(uuid.NewV4()).String(), AccountID: accountID,
		OrganizationID: organizationID, Status: status,
		StoreAccessMode: mode, AcceptedAt: acceptedAt, InvitedAt: &now,
	}
}

func replaceMembershipScope(tx *gorm.DB, membership *gen.OperatorMembership, roleIDs, storeIDs []string) error {
	var roles []*gen.OperatorRole
	if len(roleIDs) > 0 {
		if err := tx.Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
			return err
		}
	}
	if err := tx.Model(membership).Association("Roles").Replace(roles); err != nil {
		return err
	}
	var stores []*gen.Store
	if len(storeIDs) > 0 {
		if err := tx.Where("id IN ?", storeIDs).Find(&stores).Error; err != nil {
			return err
		}
	}
	return tx.Model(membership).Association("Stores").Replace(stores)
}

func validateAccessScope(db *gorm.DB, principal *auth.WorkspacePrincipal, organizationID string, input InviteInput) error {
	if err := validateMembershipWorkspace(db, principal, organizationID); err != nil {
		return err
	}
	if principal != nil && principal.WorkspaceType == auth.WorkspaceTypeHeadquarters {
		return authorization.ValidateHQRoleDelegation(db, principal, organizationID, input.RoleIDs)
	}
	if input.StoreAccessMode == gen.StoreAccessModeAllStores && len(input.StoreIDs) != 0 {
		return auth.NewError(auth.CodeValidationFailed)
	}
	if input.StoreAccessMode == gen.StoreAccessModeSelectedStores && len(input.StoreIDs) == 0 {
		return auth.NewError(auth.CodeValidationFailed)
	}
	if err := validateInviterStoreScope(principal, input); err != nil {
		return err
	}
	if err := validateRelatedCount(db, &gen.OperatorRole{}, input.RoleIDs, organizationID, ""); err != nil {
		return err
	}
	return validateRelatedCount(db, &gen.Store{}, input.StoreIDs, organizationID, string(gen.StoreLifecycleActive))
}

func validateMembershipWorkspace(db *gorm.DB, principal *auth.WorkspacePrincipal, organizationID string) error {
	organization := &gen.Organization{}
	if err := db.Where("is_delete IS NULL OR is_delete = ?", 1).First(organization, "id = ?", organizationID).Error; err != nil {
		return auth.NewError(auth.CodeWorkspaceForbidden)
	}
	expected := gen.OrganizationTypeFranchise
	if principal != nil && principal.WorkspaceType == auth.WorkspaceTypeHeadquarters {
		expected = gen.OrganizationTypeHeadquarters
	}
	if organization.Type != expected || organization.Status != gen.OrganizationStatusActive {
		return auth.NewError(auth.CodeWorkspaceForbidden)
	}
	return nil
}

func validateInviterStoreScope(principal *auth.WorkspacePrincipal, input InviteInput) error {
	if principal == nil || principal.WorkspaceType != auth.WorkspaceTypeFranchise || principal.AllStores {
		return nil
	}
	if input.StoreAccessMode == gen.StoreAccessModeAllStores {
		return auth.NewError(auth.CodeStoreScopeDenied)
	}
	for _, storeID := range input.StoreIDs {
		if !principal.HasStore(storeID) {
			return auth.NewError(auth.CodeStoreScopeDenied)
		}
	}
	return nil
}

func validateRelatedCount(db *gorm.DB, model any, ids []string, organizationID, lifecycle string) error {
	if len(ids) == 0 {
		return nil
	}
	query := db.Model(model).Where("id IN ? AND organization_id = ?", ids, organizationID).
		Where("is_delete IS NULL OR is_delete = ?", 1)
	if lifecycle != "" {
		query = query.Where("lifecycle = ?", lifecycle)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return err
	}
	if count != int64(len(ids)) {
		if lifecycle != "" {
			return auth.NewError(auth.CodeStoreNotActive)
		}
		return auth.NewError(auth.CodePermissionDenied)
	}
	return nil
}
