package customer

import (
	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func scheduleTemplateCouponDistribution(tx *gorm.DB, template *gen.CustomerCouponTemplate, nowMillis int64) error {
	if !template.Enabled {
		return nil
	}
	availableAt := nowMillis
	if template.EffectiveAt > availableAt {
		availableAt = template.EffectiveAt
	}
	status := gen.CustomerCouponDistributionJobStatusPending
	if template.TotalIssueLimit > 0 && template.IssuedCount >= template.TotalIssueLimit {
		status = gen.CustomerCouponDistributionJobStatusCompleted
	}
	job := gen.CustomerCouponDistributionJob{
		ID: uuid.Must(uuid.NewV4()).String(), Kind: gen.CustomerCouponDistributionJobKindTemplateFanout,
		Status: status, RequestKey: templateCouponDistributionKey(template.ID), AvailableAt: availableAt,
		TemplateID: &template.ID, CreatedAt: nowMillis,
	}
	updates := map[string]any{
		"kind": gen.CustomerCouponDistributionJobKindTemplateFanout, "status": status,
		"available_at": availableAt, "lease_expires_at": nil, "lease_token": nil,
		"cursor_created_at": nil, "cursor_key": nil, "attempts": 0, "last_error_code": nil,
		"template_id": template.ID, "member_id": nil,
	}
	return tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "request_key"}}, DoUpdates: clause.Assignments(updates)}).Create(&job).Error
}

func scheduleMemberCouponCatchup(tx *gorm.DB, member *gen.CustomerMember, nowMillis int64) error {
	job := gen.CustomerCouponDistributionJob{
		ID: uuid.Must(uuid.NewV4()).String(), Kind: gen.CustomerCouponDistributionJobKindMemberCatchup,
		Status: gen.CustomerCouponDistributionJobStatusPending, RequestKey: memberCouponDistributionKey(member.ID),
		AvailableAt: nowMillis, MemberID: &member.ID, CreatedAt: nowMillis,
	}
	updates := map[string]any{
		"kind": gen.CustomerCouponDistributionJobKindMemberCatchup, "status": gen.CustomerCouponDistributionJobStatusPending,
		"available_at": nowMillis, "lease_expires_at": nil, "lease_token": nil,
		"cursor_created_at": nil, "cursor_key": nil, "attempts": 0, "last_error_code": nil,
		"template_id": nil, "member_id": member.ID,
	}
	return tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "request_key"}}, DoUpdates: clause.Assignments(updates)}).Create(&job).Error
}

func templateCouponDistributionKey(templateID string) string {
	return "TEMPLATE_FANOUT:" + templateID
}

func memberCouponDistributionKey(memberID string) string {
	return "MEMBER_CATCHUP:" + memberID
}

func scheduleReactivatedMemberCouponCatchup(tx *gorm.DB, member *gen.CustomerMember, previous gen.CustomerMemberStatus, nowMillis int64) error {
	if previous != gen.CustomerMemberStatusSuspended || member.Status != gen.CustomerMemberStatusActive {
		return nil
	}
	return scheduleMemberCouponCatchup(tx, member, nowMillis)
}

func (s *Service) finishStoreCouponTemplateCreate(tx *gorm.DB, principal *auth.WorkspacePrincipal, store *gen.Store, template *gen.CustomerCouponTemplate) error {
	if err := scheduleTemplateCouponDistribution(tx, template, template.CreatedAt); err != nil {
		return err
	}
	return s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
		OrganizationID: &store.OrganizationID, StoreID: &store.ID, Action: "franchiseCoupon:manage",
		ResourceType: "customerCouponTemplate", ResourceID: template.ID, ResultCode: "SUCCESS"})
}
