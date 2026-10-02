/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package dbup

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authentication"
	"gorm.io/gorm"
)

const hqBootstrapKey = "initial-hq-administrator"

var hqPhonePattern = regexp.MustCompile(`^((13[0-9])|(15[0-9])|(18[0-9])|(17[0-9])|(14[0-9]))[0-9]{8}$`)

type hqCredentialSpec struct {
	passwordHash       string
	mustChangePassword bool
	temporaryExpires   *time.Time
}

// BootstrapHQResult 是只允许 CLI 输出一次的总部初始化凭据。
type BootstrapHQResult struct {
	AccountID         string `json:"accountId"`
	TemporaryPassword string `json:"temporaryPassword"`
}

// BootstrapHQ 原子创建首个总部管理员，并只返回一次临时密码。
func BootstrapHQ(ctx context.Context, db *gorm.DB, phone, displayName string) (BootstrapHQResult, error) {
	result := BootstrapHQResult{}
	password, err := auth.GenerateTemporaryPassword()
	if err != nil {
		return result, err
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		return result, err
	}
	now := time.Now()
	spec := hqCredentialSpec{passwordHash: hash, mustChangePassword: true, temporaryExpires: authentication.NewTemporaryPasswordExpiry(now)}
	err = bootstrapHQ(ctx, db, phone, displayName, spec, func(accountID string) {
		result = BootstrapHQResult{AccountID: accountID, TemporaryPassword: password}
	})
	return result, err
}

// HQInitialized 返回一次性总部初始化标记是否存在。
func HQInitialized(ctx context.Context, db *gorm.DB) (bool, error) {
	marker := SecurityBootstrap{}
	err := db.WithContext(ctx).Where(map[string]any{"key": hqBootstrapKey}).Take(&marker).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return err == nil, err
}

// InitializeHQ 使用用户提交的正式密码创建首个总部管理员。
func InitializeHQ(ctx context.Context, db *gorm.DB, phone, password, confirmation string) error {
	initialized, err := HQInitialized(ctx, db)
	if err != nil {
		return err
	}
	if initialized {
		return auth.NewError(auth.CodeHQAlreadyBootstrapped)
	}
	phone = strings.TrimSpace(phone)
	if !hqPhonePattern.MatchString(phone) {
		return auth.NewError(auth.CodeValidationFailed)
	}
	if password != confirmation {
		return auth.NewError(auth.CodePasswordConfirmMismatch)
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		return auth.NewError(auth.CodePasswordWeak)
	}
	spec := hqCredentialSpec{passwordHash: hash}
	return bootstrapHQ(ctx, db, phone, phone, spec, nil)
}

func bootstrapHQ(ctx context.Context, db *gorm.DB, phone, displayName string, spec hqCredentialSpec, created func(string)) error {
	var accountID string
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var transactionErr error
		accountID, transactionErr = bootstrapHQTransaction(tx, phone, displayName, spec)
		return transactionErr
	})
	if err != nil {
		return normalizeHQBootstrapError(ctx, db, err)
	}
	if created != nil {
		created(accountID)
	}
	return nil
}

func normalizeHQBootstrapError(ctx context.Context, db *gorm.DB, original error) error {
	if auth.ErrorCode(original) == auth.CodeHQAlreadyBootstrapped {
		return original
	}
	initialized, err := HQInitialized(ctx, db)
	if err == nil && initialized {
		return auth.NewError(auth.CodeHQAlreadyBootstrapped)
	}
	return original
}

func bootstrapHQTransaction(db *gorm.DB, phone, displayName string, spec hqCredentialSpec) (string, error) {
	if err := claimHQBootstrap(db); err != nil {
		return "", err
	}
	organization, role, err := headquartersRecords(db)
	if err != nil {
		return "", err
	}
	accountID := uuid.Must(uuid.NewV4()).String()
	if err := createHQAccount(db, accountID, phone, displayName, spec); err != nil {
		return "", err
	}
	if err := createHQMembership(db, accountID, organization.ID, &role); err != nil {
		return "", err
	}
	return accountID, createBootstrapAudit(db, accountID, organization.ID)
}

func claimHQBootstrap(db *gorm.DB) error {
	marker := SecurityBootstrap{Key: hqBootstrapKey, CreatedAt: time.Now()}
	if err := db.Create(&marker).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || isDuplicateBootstrap(db, err) {
			return auth.NewError(auth.CodeHQAlreadyBootstrapped)
		}
		return err
	}
	return nil
}

func isDuplicateBootstrap(db *gorm.DB, original error) bool {
	var count int64
	err := db.Model(&SecurityBootstrap{}).Where(map[string]any{"key": hqBootstrapKey}).Count(&count).Error
	return err == nil && count > 0 && original != nil
}

func headquartersRecords(db *gorm.DB) (gen.Organization, gen.OperatorRole, error) {
	var organization gen.Organization
	err := db.Where("type = ?", gen.OrganizationTypeHeadquarters).First(&organization).Error
	if err != nil {
		return organization, gen.OperatorRole{}, err
	}
	var role gen.OperatorRole
	err = db.Where("organization_id = ? AND kind = ?", organization.ID, gen.RoleKindHqSuperAdmin).First(&role).Error
	return organization, role, err
}

func createHQAccount(db *gorm.DB, accountID, phone, displayName string, spec hqCredentialSpec) error {
	now := time.Now()
	account := gen.Account{
		ID: accountID, Phone: phone, DisplayName: displayName,
		Status: gen.AccountStatusActive, MustChangePassword: spec.mustChangePassword, CredentialVersion: 1,
	}
	if err := db.Create(&account).Error; err != nil {
		return err
	}
	if !spec.mustChangePassword {
		if err := db.Model(&account).Update("must_change_password", false).Error; err != nil {
			return err
		}
	}
	credential := authentication.AccountCredential{
		AccountID: accountID, PasswordHash: spec.passwordHash,
		TemporaryPasswordExpiresAt: spec.temporaryExpires, PasswordChangedAt: now, UpdatedAt: now,
	}
	return db.Create(&credential).Error
}

func createHQMembership(db *gorm.DB, accountID, organizationID string, role *gen.OperatorRole) error {
	now := time.Now()
	membership := gen.OperatorMembership{
		ID: uuid.Must(uuid.NewV4()).String(), AccountID: accountID,
		OrganizationID: organizationID, Status: gen.MembershipStatusActive,
		StoreAccessMode: gen.StoreAccessModeAllStores, AcceptedAt: &now,
	}
	if err := db.Create(&membership).Error; err != nil {
		return err
	}
	return db.Model(&membership).Association("Roles").Append(role)
}

func createBootstrapAudit(db *gorm.DB, accountID, organizationID string) error {
	metadata, err := json.Marshal(map[string]string{"source": "bootstrap-hq"})
	if err != nil {
		return err
	}
	value := string(metadata)
	audit := gen.AuditLog{
		ID: uuid.Must(uuid.NewV4()).String(), ActorAccountID: &accountID,
		OrganizationID: &organizationID, Action: "hq:bootstrap",
		ResourceType: "account", ResourceID: &accountID,
		ResultCode: "SUCCESS", MetadataJSON: &value,
	}
	return db.Create(&audit).Error
}
