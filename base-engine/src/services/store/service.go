/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package store

import (
	"context"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
)

// CreateInput 描述加盟或总部直营门店基础资料。
type CreateInput struct {
	Code           string
	Name           string
	OrganizationID *string
}

// ReviewInput 描述总部独立批准或退回动作。
type ReviewInput struct {
	StoreID         string
	Approved        bool
	RejectionReason string
}

// Service 管理门店创建与准入生命周期。
type Service struct {
	db               *gorm.DB
	audit            *audit.Service
	documentRoot     string
	documentMu       sync.Mutex
	pendingDocuments map[string]pendingDocument
}

// NewService 创建门店服务。
func NewService(db *gorm.DB, auditService *audit.Service) *Service {
	return &Service{db: db, audit: auditService, pendingDocuments: make(map[string]pendingDocument)}
}

// Create 创建加盟 DRAFT 或总部 ACTIVE 门店。
func (s *Service) Create(ctx context.Context, principal *auth.WorkspacePrincipal, input CreateInput) (*gen.Store, error) {
	organizationID, action, lifecycle, err := createPolicy(principal)
	if err != nil {
		return nil, err
	}
	if input.OrganizationID != nil && *input.OrganizationID != organizationID {
		return nil, auth.NewError(auth.CodePermissionDenied)
	}
	if err := authorization.Authorize(principal, authorization.Intent{Action: action, Mode: authorization.AccessCreate, ResourceOrganizationID: &organizationID}); err != nil {
		return nil, err
	}
	item := &gen.Store{
		ID: uuid.Must(uuid.NewV4()).String(), Code: input.Code, Name: input.Name,
		OrganizationID: organizationID, Lifecycle: lifecycle,
	}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(item).Error; err != nil {
			return err
		}
		return s.writeAudit(tx, principal, item, "store:create", "")
	})
	return item, err
}

// Submit 将 DRAFT 或 REJECTED 门店提交总部审核。
func (s *Service) Submit(ctx context.Context, principal *auth.WorkspacePrincipal, storeID string) (*gen.Store, error) {
	if principal == nil || principal.WorkspaceType != auth.WorkspaceTypeFranchise {
		return nil, auth.NewError(auth.CodeWorkspaceForbidden)
	}
	item := &gen.Store{}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("is_delete IS NULL OR is_delete = ?", 1).First(item, "id = ?", storeID).Error; err != nil {
			return auth.NewError(auth.CodePermissionDenied)
		}
		intent := authorization.Intent{Action: "store:submit", Mode: authorization.AccessUpdate, ResourceOrganizationID: &item.OrganizationID}
		if !principal.AllStores {
			intent.StoreID = &item.ID
		}
		if err := authorization.Authorize(principal, intent); err != nil {
			return err
		}
		now := time.Now()
		updates := map[string]any{"lifecycle": gen.StoreLifecyclePendingApproval, "submitted_at": now, "rejection_reason": nil}
		result := tx.Model(&gen.Store{}).Where("id = ? AND lifecycle IN ?", item.ID, []gen.StoreLifecycle{gen.StoreLifecycleDraft, gen.StoreLifecycleRejected}).Updates(updates)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return auth.NewError(auth.CodeConflict)
		}
		item.Lifecycle = gen.StoreLifecyclePendingApproval
		return s.writeAudit(tx, principal, item, "store:submit", "")
	})
	return item, err
}

// Review 对待审门店执行批准或结构化退回。
func (s *Service) Review(ctx context.Context, principal *auth.WorkspacePrincipal, input ReviewInput) (*gen.Store, error) {
	if principal == nil || principal.WorkspaceType != auth.WorkspaceTypeHeadquarters {
		return nil, auth.NewError(auth.CodeWorkspaceForbidden)
	}
	input, action, target, err := reviewPolicy(input)
	if err != nil {
		return nil, err
	}
	if err := authorization.Authorize(principal, authorization.Intent{Action: action, Mode: authorization.AccessUpdate}); err != nil {
		return nil, err
	}
	item := &gen.Store{}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("is_delete IS NULL OR is_delete = ?", 1).First(item, "id = ?", input.StoreID).Error; err != nil {
			return auth.NewError(auth.CodePermissionDenied)
		}
		now := time.Now()
		updates := reviewUpdates(input, target, principal.AccountID, now)
		result := tx.Model(&gen.Store{}).Where("id = ? AND lifecycle = ?", item.ID, gen.StoreLifecyclePendingApproval).Updates(updates)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return auth.NewError(auth.CodeConflict)
		}
		item.Lifecycle = target
		return s.writeAudit(tx, principal, item, action, input.RejectionReason)
	})
	return item, err
}

func reviewPolicy(input ReviewInput) (ReviewInput, string, gen.StoreLifecycle, error) {
	input.RejectionReason = strings.TrimSpace(input.RejectionReason)
	if utf8.RuneCountInString(input.RejectionReason) > 512 || (!input.Approved && input.RejectionReason == "") {
		return ReviewInput{}, "", "", auth.NewError(auth.CodeValidationFailed)
	}
	if input.Approved {
		return input, "store:approve", gen.StoreLifecycleActive, nil
	}
	return input, "store:reject", gen.StoreLifecycleRejected, nil
}

func reviewUpdates(input ReviewInput, target gen.StoreLifecycle, accountID string, now time.Time) map[string]any {
	updates := map[string]any{"lifecycle": target, "reviewed_at": now, "reviewed_by_account_id": accountID}
	if input.Approved {
		updates["rejection_reason"] = nil
	} else {
		updates["rejection_reason"] = input.RejectionReason
	}
	return updates
}

func createPolicy(principal *auth.WorkspacePrincipal) (string, string, gen.StoreLifecycle, error) {
	if principal == nil || principal.OrganizationID == nil {
		return "", "", "", auth.NewError(auth.CodeWorkspaceForbidden)
	}
	if principal.WorkspaceType == auth.WorkspaceTypeHeadquarters {
		return *principal.OrganizationID, "hqStore:create", gen.StoreLifecycleActive, nil
	}
	if principal.WorkspaceType == auth.WorkspaceTypeFranchise {
		return *principal.OrganizationID, "store:create", gen.StoreLifecycleDraft, nil
	}
	return "", "", "", auth.NewError(auth.CodeWorkspaceForbidden)
}

func (s *Service) writeAudit(tx *gorm.DB, principal *auth.WorkspacePrincipal, item *gen.Store, action, reason string) error {
	if s.audit == nil {
		return nil
	}
	return s.audit.Write(tx, audit.Entry{
		ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
		OrganizationID: &item.OrganizationID, StoreID: &item.ID,
		Action: action, ResourceType: "store", ResourceID: item.ID,
		ResultCode: "SUCCESS", Metadata: audit.Metadata{ReasonCode: reason, TargetStatus: string(item.Lifecycle)},
	})
}
