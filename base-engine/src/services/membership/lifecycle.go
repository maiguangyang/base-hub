/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package membership

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

// AcceptInvitation 激活属于当前账号且未失效的成员邀请。
func (s *Service) AcceptInvitation(ctx context.Context, principal *auth.WorkspacePrincipal, invitationID string) (*gen.OperatorMembership, error) {
	var membership *gen.OperatorMembership
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		invitation := &gen.MembershipInvitation{}
		if err := tx.Where("is_delete IS NULL OR is_delete = ?", 1).First(invitation, "id = ?", invitationID).Error; err != nil {
			return auth.NewError(auth.CodeInvitationExpired)
		}
		if !invitationUsable(invitation, time.Now()) {
			return auth.NewError(auth.CodeInvitationExpired)
		}
		membership = &gen.OperatorMembership{}
		if err := tx.First(membership, "id = ?", invitation.MembershipID).Error; err != nil {
			return err
		}
		if principal == nil || membership.AccountID != principal.AccountID {
			return auth.NewError(auth.CodePermissionDenied)
		}
		if err := validateInvitationActivation(tx, membership); err != nil {
			return err
		}
		now := time.Now()
		if err := tx.Model(membership).Updates(map[string]any{"status": gen.MembershipStatusActive, "accepted_at": now}).Error; err != nil {
			return err
		}
		if err := tx.Model(invitation).Update("accepted_at", now).Error; err != nil {
			return err
		}
		membership.Status = gen.MembershipStatusActive
		membership.AcceptedAt = &now
		return s.writeAudit(tx, principal, membership, "invitation:accept", audit.Metadata{})
	})
	return membership, err
}

// ChangeMembershipStatus 修改成员状态，并保护最后一个活跃老板。
func (s *Service) ChangeMembershipStatus(ctx context.Context, principal *auth.WorkspacePrincipal, membershipID string, status gen.MembershipStatus) error {
	organizationID, err := authorizeMembershipChange(principal, "operatorMembership:update", authorization.AccessUpdate)
	if err != nil {
		auth.LogAuthorizationDenied(principal, membershipStatusAuditAction(principal), "operatorMembership", membershipID, err)
		return err
	}
	if !targetMembershipStatusAllowed(status) {
		return auth.NewError(auth.CodeValidationFailed)
	}
	outcome, err := s.changeMembershipStatus(ctx, principal, organizationID, membershipID, status)
	if err == nil {
		publishMembershipRevocations(s.publisher, outcome.sessionIDs, outcome.organizationID, outcome.occurredAt)
	}
	return err
}

type membershipStatusOutcome struct {
	sessionIDs     []string
	organizationID string
	occurredAt     time.Time
}

func (s *Service) changeMembershipStatus(ctx context.Context, principal *auth.WorkspacePrincipal, organizationID, membershipID string, status gen.MembershipStatus) (membershipStatusOutcome, error) {
	outcome := membershipStatusOutcome{occurredAt: time.Now()}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return s.applyMembershipStatus(ctx, tx, principal, organizationID, membershipID, status, &outcome)
	})
	return outcome, err
}

func (s *Service) applyMembershipStatus(ctx context.Context, tx *gorm.DB, principal *auth.WorkspacePrincipal, organizationID, membershipID string, status gen.MembershipStatus, outcome *membershipStatusOutcome) error {
	membership := &gen.OperatorMembership{}
	if err := tx.Where("is_delete IS NULL OR is_delete = ?", 1).First(membership, "id = ? AND organization_id = ?", membershipID, organizationID).Error; err != nil {
		return auth.NewError(auth.CodePermissionDenied)
	}
	if !membershipStatusManageable(membership.Status) {
		return auth.NewError(auth.CodePermissionDenied)
	}
	if err := authorization.ValidateHQMembershipTarget(tx, principal, membership); err != nil {
		auth.LogAuthorizationDenied(principal, membershipStatusAuditAction(principal), "operatorMembership", membership.ID, err)
		return err
	}
	if status != gen.MembershipStatusActive {
		if err := authorization.EnsurePrivilegedMembershipRemains(tx, membership); err != nil {
			return err
		}
	}
	if err := tx.Model(membership).Update("status", status).Error; err != nil {
		return err
	}
	membership.Status = status
	if err := s.revokeStatusSessions(ctx, tx, principal, membership, status, outcome); err != nil {
		return err
	}
	return s.writeAudit(tx, principal, membership, membershipStatusAuditAction(principal), audit.Metadata{})
}

func (s *Service) revokeStatusSessions(ctx context.Context, tx *gorm.DB, principal *auth.WorkspacePrincipal, membership *gen.OperatorMembership, status gen.MembershipStatus, outcome *membershipStatusOutcome) error {
	if !shouldRevokeMembershipSessions(principal, status) || s.sessions == nil {
		return nil
	}
	ids, err := s.sessions.RevokeWorkspaceSessions(ctx, tx, membership.AccountID, membership.OrganizationID, sessionservice.RevocationCodeAuthorityChanged, outcome.occurredAt)
	outcome.sessionIDs = ids
	outcome.organizationID = membership.OrganizationID
	return err
}

func shouldRevokeMembershipSessions(principal *auth.WorkspacePrincipal, status gen.MembershipStatus) bool {
	return principal != nil && principal.WorkspaceType == auth.WorkspaceTypeHeadquarters && status != gen.MembershipStatusActive
}

func publishMembershipRevocations(publisher sessionservice.Publisher, sessionIDs []string, organizationID string, occurredAt time.Time) {
	if publisher == nil {
		return
	}
	for _, sessionID := range sessionIDs {
		organization := organizationID
		publisher.PublishSession(sessionID, &gen.SessionEvent{Code: gen.SessionEventCodeSessionRevoked, SessionID: sessionID, OrganizationID: &organization, OccurredAt: occurredAt})
	}
}

func membershipStatusAuditAction(principal *auth.WorkspacePrincipal) string {
	if principal != nil && principal.WorkspaceType == auth.WorkspaceTypeHeadquarters {
		return "hqAdministrator:change_status"
	}
	return "membership:status"
}

func invitationUsable(invitation *gen.MembershipInvitation, now time.Time) bool {
	return invitation.AcceptedAt == nil && invitation.RevokedAt == nil && invitation.ExpiresAt.After(now)
}

func targetMembershipStatusAllowed(status gen.MembershipStatus) bool {
	return status == gen.MembershipStatusActive || status == gen.MembershipStatusSuspended || status == gen.MembershipStatusLeft
}

func membershipStatusManageable(status gen.MembershipStatus) bool {
	return status == gen.MembershipStatusActive || status == gen.MembershipStatusSuspended
}

func authorizeMembershipChange(principal *auth.WorkspacePrincipal, action string, mode authorization.AccessMode) (string, error) {
	if principal == nil || principal.OrganizationID == nil {
		return "", auth.NewError(auth.CodeWorkspaceForbidden)
	}
	workspaceAction, err := authorization.WorkspaceAction(principal, action)
	if err != nil {
		return "", err
	}
	organizationID := *principal.OrganizationID
	err = authorization.Authorize(principal, authorization.Intent{Action: workspaceAction, Mode: mode, ResourceOrganizationID: &organizationID})
	return organizationID, err
}

func (s *Service) writeAudit(tx *gorm.DB, principal *auth.WorkspacePrincipal, membership *gen.OperatorMembership, action string, metadata audit.Metadata) error {
	if s.audit == nil {
		return nil
	}
	metadata.TargetStatus = string(membership.Status)
	return s.audit.Write(tx, audit.Entry{
		ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
		OrganizationID: &membership.OrganizationID, Action: action,
		ResourceType: "operatorMembership", ResourceID: membership.ID,
		ResultCode: "SUCCESS", Metadata: audit.MetadataForPrincipal(principal, metadata),
	})
}
