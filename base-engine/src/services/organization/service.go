/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package organization

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

const invitationLifetime = 72 * time.Hour
const franchiseOwnerRoleName = "role.franchiseOwner"

// ProvisionInput 描述总部开通加盟组织所需信息。
type ProvisionInput struct {
	Code             string
	Name             string
	OwnerPhone       string
	OwnerDisplayName string
	OwnerEmail       *string
}

// ProvisionResult 返回组织、老板成员及仅一次可见的临时密码。
type ProvisionResult struct {
	Organization      *gen.Organization
	Membership        *gen.OperatorMembership
	TemporaryPassword *string
	InvitationPending bool
}

// Service 管理组织治理事务。
type Service struct {
	db        *gorm.DB
	audit     *audit.Service
	publisher sessionservice.Publisher
}

// NewService 创建组织服务。
func NewService(db *gorm.DB, auditService *audit.Service, publishers ...sessionservice.Publisher) *Service {
	service := &Service{db: db, audit: auditService}
	if len(publishers) > 0 {
		service.publisher = publishers[0]
	}
	return service
}

// ProvisionFranchise 原子开通加盟组织与老板身份。
func (s *Service) ProvisionFranchise(ctx context.Context, principal *auth.WorkspacePrincipal, input ProvisionInput) (*ProvisionResult, error) {
	if principal == nil || principal.WorkspaceType != auth.WorkspaceTypeHeadquarters {
		return nil, auth.NewError(auth.CodeWorkspaceForbidden)
	}
	if err := authorization.Authorize(principal, authorization.Intent{Action: "franchise:provision", Mode: authorization.AccessCreate}); err != nil {
		return nil, err
	}
	input, err := normalizeProvisionInput(input)
	if err != nil {
		return nil, err
	}
	exists, err := organizationCodeExists(ctx, s.db, input.Code)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	var result *ProvisionResult
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		provisioned, err := s.provisionTransaction(tx, principal, input, time.Now())
		result = provisioned
		return err
	})
	if err != nil {
		exists, lookupErr := organizationCodeExists(ctx, s.db, input.Code)
		if lookupErr == nil && exists {
			return nil, auth.NewError(auth.CodeValidationFailed)
		}
	}
	return result, err
}

func (s *Service) provisionTransaction(tx *gorm.DB, principal *auth.WorkspacePrincipal, input ProvisionInput, now time.Time) (*ProvisionResult, error) {
	organization := &gen.Organization{
		ID: uuid.Must(uuid.NewV4()).String(), Code: input.Code, Name: input.Name,
		Type: gen.OrganizationTypeFranchise, Status: gen.OrganizationStatusActive,
	}
	if err := tx.Create(organization).Error; err != nil {
		return nil, err
	}
	account, existing, password, err := findOrCreateOwner(tx, input, now)
	if err != nil {
		return nil, err
	}
	if err := tx.Model(organization).Update("initial_account_id", account.ID).Error; err != nil {
		return nil, err
	}
	organization.InitialAccountID = &account.ID
	if err := tx.Create(&gen.FranchiseOpeningRecord{
		ID: uuid.Must(uuid.NewV4()).String(), RecordNumber: "SYS-" + uuid.Must(uuid.NewV4()).String(),
		Source: gen.FranchiseOpeningSourceSystemProvision, OrganizationID: organization.ID,
		InitialAccountID: account.ID, RecordedByAccountID: principal.AccountID,
	}).Error; err != nil {
		return nil, err
	}
	role, err := createOwnerRole(tx, organization.ID)
	if err != nil {
		return nil, err
	}
	membership, err := createOwnerMembership(tx, account.ID, organization.ID, role, existing, now)
	if err != nil {
		return nil, err
	}
	if existing {
		if err := createOwnerInvitation(tx, membership.ID, principal.AccountID, now); err != nil {
			return nil, err
		}
	}
	if err := s.auditProvision(tx, principal, organization); err != nil {
		return nil, err
	}
	return &ProvisionResult{
		Organization: organization, Membership: membership,
		TemporaryPassword: password, InvitationPending: existing,
	}, nil
}

func findOrCreateOwner(tx *gorm.DB, input ProvisionInput, now time.Time) (*gen.Account, bool, *string, error) {
	account := &gen.Account{}
	err := tx.Where("phone = ?", input.OwnerPhone).First(account).Error
	if err == nil {
		if err := validateExistingOwner(tx, account); err != nil {
			return nil, false, nil, err
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
		ID: uuid.Must(uuid.NewV4()).String(), Phone: input.OwnerPhone,
		DisplayName: input.OwnerDisplayName, Email: input.OwnerEmail,
		Status: gen.AccountStatusActive, MustChangePassword: true, CredentialVersion: 1,
	}
	if err := tx.Create(account).Error; err != nil {
		return nil, false, nil, err
	}
	credential := authentication.AccountCredential{
		AccountID: account.ID, PasswordHash: hash,
		TemporaryPasswordExpiresAt: authentication.NewTemporaryPasswordExpiry(now), PasswordChangedAt: now, UpdatedAt: now,
	}
	if err := tx.Create(&credential).Error; err != nil {
		return nil, false, nil, err
	}
	return account, false, &password, nil
}

func validateExistingOwner(tx *gorm.DB, account *gen.Account) error {
	if account.Status != gen.AccountStatusActive || !recordActive(account.IsDelete) {
		return auth.NewError(auth.CodePermissionDenied)
	}
	headquartersMembership, err := ownerHasHeadquartersMembership(tx, account.ID)
	if err != nil {
		return err
	}
	if headquartersMembership {
		return auth.NewError(auth.CodePermissionDenied)
	}
	return nil
}

func ownerHasHeadquartersMembership(tx *gorm.DB, accountID string) (bool, error) {
	var count int64
	err := tx.Model(&gen.OperatorMembership{}).
		Joins("JOIN organizations ON organizations.id = operator_memberships.organization_id").
		Where("operator_memberships.account_id = ? AND organizations.type = ?", accountID, gen.OrganizationTypeHeadquarters).
		Where("operator_memberships.is_delete IS NULL OR operator_memberships.is_delete = ?", 1).
		Count(&count).Error
	return count != 0, err
}

func recordActive(isDelete *int64) bool {
	return isDelete == nil || *isDelete == 1
}

func createOwnerRole(tx *gorm.DB, organizationID string) (*gen.OperatorRole, error) {
	role := &gen.OperatorRole{
		ID: uuid.Must(uuid.NewV4()).String(), Name: franchiseOwnerRoleName,
		Kind: gen.RoleKindFranchiseOwner, OrganizationID: organizationID,
	}
	return role, tx.Create(role).Error
}

func createOwnerMembership(tx *gorm.DB, accountID, organizationID string, role *gen.OperatorRole, invited bool, now time.Time) (*gen.OperatorMembership, error) {
	status := gen.MembershipStatusActive
	var acceptedAt *time.Time = &now
	if invited {
		status = gen.MembershipStatusInvited
		acceptedAt = nil
	}
	membership := &gen.OperatorMembership{
		ID: uuid.Must(uuid.NewV4()).String(), AccountID: accountID,
		OrganizationID: organizationID, Status: status,
		StoreAccessMode: gen.StoreAccessModeAllStores, AcceptedAt: acceptedAt,
	}
	if err := tx.Create(membership).Error; err != nil {
		return nil, err
	}
	return membership, tx.Model(membership).Association("Roles").Append(role)
}

func createOwnerInvitation(tx *gorm.DB, membershipID, inviterID string, now time.Time) error {
	invitation := gen.MembershipInvitation{
		ID: uuid.Must(uuid.NewV4()).String(), MembershipID: membershipID,
		InvitedByAccountID: inviterID, ExpiresAt: now.Add(invitationLifetime),
	}
	return tx.Create(&invitation).Error
}

func (s *Service) auditProvision(tx *gorm.DB, principal *auth.WorkspacePrincipal, organization *gen.Organization) error {
	if s.audit == nil {
		return nil
	}
	return s.audit.Write(tx, audit.Entry{
		ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
		OrganizationID: &organization.ID, Action: "franchise:provision",
		ResourceType: "organization", ResourceID: organization.ID,
		ResultCode: "SUCCESS", Metadata: audit.Metadata{TargetStatus: string(organization.Status)},
	})
}
