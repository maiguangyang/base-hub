/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package organization

import (
	"context"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"base-engine/src/services/authorization"
	sessionservice "base-engine/src/services/session"
	"gorm.io/gorm"
)

const organizationSuspendedCode = "ORGANIZATION_SUSPENDED"

// Suspend 暂停加盟组织并在同一事务撤销其全部在线会话。
func (s *Service) Suspend(ctx context.Context, principal *auth.WorkspacePrincipal, organizationID, reasonCode string) error {
	if principal == nil || principal.WorkspaceType != auth.WorkspaceTypeHeadquarters {
		return auth.NewError(auth.CodeWorkspaceForbidden)
	}
	if err := authorization.Authorize(principal, authorization.Intent{Action: "organization:suspend", Mode: authorization.AccessUpdate}); err != nil {
		return err
	}
	reasonCode, err := normalizeReason(reasonCode, 64)
	if err != nil {
		return err
	}
	var revokedSessionIDs []string
	now := time.Now()
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&gen.Organization{}).
			Where("id = ? AND type = ? AND status = ?", organizationID, gen.OrganizationTypeFranchise, gen.OrganizationStatusActive).
			Where("is_delete IS NULL OR is_delete = ?", 1).
			Updates(map[string]any{"status": gen.OrganizationStatusSuspended, "suspended_at": now, "suspension_reason_code": reasonCode})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return auth.NewError(auth.CodePermissionDenied)
		}
		sessions := sessionservice.NewService(0)
		ids, err := sessions.RevokeOrganizationSessions(ctx, tx, organizationID, organizationSuspendedCode, now)
		revokedSessionIDs = ids
		if err != nil {
			return err
		}
		return s.auditGovernance(tx, principal, organizationID, "organization:suspend", reasonCode)
	})
	if err != nil {
		return err
	}
	s.publishSuspension(revokedSessionIDs, organizationID, now)
	return nil
}

// Restore 恢复组织接收新会话，但不复活已撤销会话。
func (s *Service) Restore(ctx context.Context, principal *auth.WorkspacePrincipal, organizationID string) error {
	if principal == nil || principal.WorkspaceType != auth.WorkspaceTypeHeadquarters {
		return auth.NewError(auth.CodeWorkspaceForbidden)
	}
	if err := authorization.Authorize(principal, authorization.Intent{Action: "organization:restore", Mode: authorization.AccessUpdate}); err != nil {
		return err
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&gen.Organization{}).
			Where("id = ? AND type = ? AND status = ?", organizationID, gen.OrganizationTypeFranchise, gen.OrganizationStatusSuspended).
			Where("is_delete IS NULL OR is_delete = ?", 1).
			Updates(map[string]any{"status": gen.OrganizationStatusActive, "suspended_at": nil, "suspension_reason_code": nil})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return auth.NewError(auth.CodePermissionDenied)
		}
		return s.auditGovernance(tx, principal, organizationID, "organization:restore", "")
	})
}

func (s *Service) publishSuspension(sessionIDs []string, organizationID string, occurredAt time.Time) {
	if s.publisher == nil {
		return
	}
	for _, sessionID := range sessionIDs {
		event := &gen.SessionEvent{
			Code: gen.SessionEventCodeOrganizationSuspended, SessionID: sessionID,
			OrganizationID: &organizationID, OccurredAt: occurredAt,
		}
		s.publisher.PublishSession(sessionID, event)
	}
}

func (s *Service) auditGovernance(tx *gorm.DB, principal *auth.WorkspacePrincipal, organizationID, action, reason string) error {
	if s.audit == nil {
		return nil
	}
	return s.audit.Write(tx, audit.Entry{
		ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
		OrganizationID: &organizationID, Action: action,
		ResourceType: "organization", ResourceID: organizationID,
		ResultCode: "SUCCESS", Metadata: audit.Metadata{ReasonCode: reason},
	})
}
