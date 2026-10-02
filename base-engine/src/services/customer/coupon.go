package customer

import (
	"context"
	"strings"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *Service) GrantCoupon(ctx context.Context, principal *auth.WorkspacePrincipal, templateID, memberID, requestKey string) (*gen.CustomerCouponGrant, error) {
	hqID, err := headquartersID(principal, "hqCustomerCoupon:grant", authorization.AccessCreate)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(requestKey) == "" || len(requestKey) > 128 {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	var result *gen.CustomerCouponGrant
	for attempt := 0; attempt < 3; attempt++ {
		err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			var transactionErr error
			result, transactionErr = s.grantCouponInTransaction(tx, principal, hqID, templateID, memberID, requestKey)
			return transactionErr
		})
		if err == nil {
			return result, nil
		}
		if mysqlErrorNumber(err) == 1062 {
			return s.resolveDuplicateCouponGrant(ctx, err, hqID, templateID, memberID, requestKey)
		}
		if mysqlErrorNumber(err) != 1213 && mysqlErrorNumber(err) != 1205 {
			return nil, err
		}
	}
	return nil, auth.NewError(auth.CodeConflict)
}

func (s *Service) resolveDuplicateCouponGrant(ctx context.Context, duplicateErr error, hqID, templateID, memberID, requestKey string) (*gen.CustomerCouponGrant, error) {
	replay, err := findCouponGrant(s.db.WithContext(ctx), hqID, templateID, memberID, requestKey)
	if err != nil {
		return nil, err
	}
	if replay != nil {
		return replay, nil
	}
	return nil, duplicateErr
}

func (s *Service) grantCouponInTransaction(tx *gorm.DB, principal *auth.WorkspacePrincipal, hqID, templateID, memberID, key string) (*gen.CustomerCouponGrant, error) {
	grant, _, err := s.createCouponGrantInTransaction(tx, couponGrantInput{
		TemplateID: templateID, MemberID: memberID, RequestKey: key,
		ExpectedOrganizationID: hqID, Principal: principal, NowMillis: s.now().UnixMilli(),
	})
	return grant, err
}

func checkCouponGrantCaps(tx *gorm.DB, template *gen.CustomerCouponTemplate, memberID string) error {
	if template.TotalIssueLimit > 0 && template.IssuedCount >= template.TotalIssueLimit {
		return auth.NewError(auth.CodeConflict)
	}
	var count int64
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&gen.CustomerCouponGrant{}).
		Where("template_id = ? AND member_id = ?", template.ID, memberID).Count(&count).Error
	if err != nil {
		return err
	}
	if count >= template.PerMemberLimit {
		return auth.NewError(auth.CodeConflict)
	}
	return nil
}

type couponGrantInput struct {
	TemplateID, MemberID, RequestKey string
	ExpectedOrganizationID           string
	ExpectedStoreID                  *string
	IssuerRequestDigest              *string
	Principal                        *auth.WorkspacePrincipal
	AuditOrganizationID              *string
	AuditStoreID                     *string
	Automatic                        bool
	EnforceMemberPriority            bool
	NowMillis                        int64
}

func (s *Service) createCouponGrantInTransaction(tx *gorm.DB, input couponGrantInput) (*gen.CustomerCouponGrant, bool, error) {
	template, member, replay, err := s.prepareCouponGrant(tx, input)
	if err != nil || replay != nil {
		return replay, false, err
	}
	if err := checkCouponGrantCaps(tx, template, member.ID); err != nil {
		return nil, false, err
	}
	grant := newCouponGrant(template, member.ID, input)
	if err := tx.Create(grant).Error; err != nil {
		return nil, false, err
	}
	if err := tx.Model(template).Update("issued_count", template.IssuedCount+1).Error; err != nil {
		return nil, false, err
	}
	if err := s.auditCouponGrantInsert(tx, template, grant.ID, input); err != nil {
		return nil, false, err
	}
	return grant, true, nil
}

func (s *Service) prepareCouponGrant(tx *gorm.DB, input couponGrantInput) (*gen.CustomerCouponTemplate, *gen.CustomerMember, *gen.CustomerCouponGrant, error) {
	if input.Automatic && input.RequestKey != autoCouponRequestKey(input.TemplateID, input.MemberID) {
		return nil, nil, nil, auth.NewError(auth.CodeValidationFailed)
	}
	template, err := lockCouponGrantTemplate(tx, input)
	if err != nil {
		return nil, nil, nil, err
	}
	if !input.Automatic {
		replay, err := findCouponGrantReplay(tx, input, template.OrganizationID)
		if err != nil || replay != nil {
			return template, nil, replay, err
		}
	}
	if !couponTemplateDistributable(template, input.NowMillis) {
		return nil, nil, nil, auth.NewError(auth.CodeConflict)
	}
	member, err := lockedActiveMember(tx, template.OrganizationID, input.MemberID)
	if err != nil {
		return nil, nil, nil, err
	}
	if input.Automatic {
		replay, err := prepareAutomaticCouponGrant(tx, input, template, member)
		return template, member, replay, err
	}
	return template, member, nil, nil
}

func lockCouponGrantTemplate(tx *gorm.DB, input couponGrantInput) (*gen.CustomerCouponTemplate, error) {
	var template gen.CustomerCouponTemplate
	query := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", input.TemplateID)
	if input.ExpectedStoreID != nil {
		query = query.Where("applicable_store_id = ? AND issuer_scope = ?", *input.ExpectedStoreID, gen.CouponIssuerScopeStore)
	} else {
		query = query.Where("organization_id = ? AND issuer_scope = ?", input.ExpectedOrganizationID, gen.CouponIssuerScopeHeadquarters)
	}
	err := query.First(&template).Error
	if err == gorm.ErrRecordNotFound {
		return nil, auth.NewError(auth.CodePermissionDenied)
	}
	if err != nil {
		return nil, err
	}
	if template.PerMemberLimit < 1 || template.TotalIssueLimit < 0 {
		return nil, auth.NewError(auth.CodeConflict)
	}
	return &template, nil
}

func findCouponGrantReplay(tx *gorm.DB, input couponGrantInput, organizationID string) (*gen.CustomerCouponGrant, error) {
	if input.Automatic {
		var grant gen.CustomerCouponGrant
		err := tx.Where("template_id = ? AND member_id = ?", input.TemplateID, input.MemberID).First(&grant).Error
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return &grant, err
	}
	if input.ExpectedStoreID != nil && input.IssuerRequestDigest != nil {
		return findStoreCouponReplay(tx, *input.IssuerRequestDigest, input.TemplateID, input.MemberID, input.RequestKey)
	}
	return findCouponGrant(tx, organizationID, input.TemplateID, input.MemberID, input.RequestKey)
}

func findExactAutoCouponGrant(tx *gorm.DB, templateID, memberID string) (*gen.CustomerCouponGrant, error) {
	var grant gen.CustomerCouponGrant
	err := tx.Where("template_id = ? AND member_id = ? AND request_key = ?", templateID, memberID,
		autoCouponRequestKey(templateID, memberID)).First(&grant).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &grant, nil
}

func (s *Service) resolveAutomaticCouponGrantDuplicate(ctx context.Context, duplicateErr error, templateID, memberID string) (*gen.CustomerCouponGrant, bool, error) {
	replay, err := findExactAutoCouponGrant(s.db.WithContext(ctx), templateID, memberID)
	if err != nil {
		return nil, false, err
	}
	if replay != nil {
		return replay, false, nil
	}
	return nil, false, duplicateErr
}

func newCouponGrant(template *gen.CustomerCouponTemplate, memberID string, input couponGrantInput) *gen.CustomerCouponGrant {
	return &gen.CustomerCouponGrant{
		ID: uuid.Must(uuid.NewV4()).String(), TemplateID: template.ID, MemberID: memberID,
		RequestKey: input.RequestKey, IssuerRequestDigest: input.IssuerRequestDigest,
		Status:    gen.CustomerCouponGrantStatusPendingActivation,
		AmountFen: template.AmountFen, MinSpendFen: template.MinSpendFen,
		DaysAfterActivation: template.DaysAfterActivation,
		IssuedAt:            input.NowMillis, CreatedAt: input.NowMillis,
	}
}

func (s *Service) auditCouponGrantInsert(tx *gorm.DB, template *gen.CustomerCouponTemplate, grantID string, input couponGrantInput) error {
	if input.Automatic {
		action := "hqCustomerCoupon:grant"
		var storeID *string
		if template.IssuerScope == gen.CouponIssuerScopeStore {
			if template.ApplicableStoreID == nil {
				return auth.NewError(auth.CodeInternalError)
			}
			action, storeID = "franchiseCoupon:grant", template.ApplicableStoreID
		}
		return s.audit.Write(tx, audit.Entry{OrganizationID: &template.OrganizationID, StoreID: storeID,
			Action: action, ResourceType: "customerCouponGrant", ResourceID: grantID, ResultCode: "SUCCESS",
			Metadata: audit.Metadata{Source: "AUTO_COUPON_DISTRIBUTION"}})
	}
	if input.ExpectedStoreID == nil {
		return s.auditCouponGrant(tx, input.Principal, grantID, "grant", "grant")
	}
	return s.audit.Write(tx, audit.Entry{ActorAccountID: input.Principal.AccountID, SessionID: &input.Principal.SessionID,
		OrganizationID: input.AuditOrganizationID, StoreID: input.AuditStoreID, Action: "franchiseCoupon:grant",
		ResourceType: "customerCouponGrant", ResourceID: grantID, ResultCode: "SUCCESS"})
}

func findCouponGrant(tx *gorm.DB, hqID, templateID, memberID, key string) (*gen.CustomerCouponGrant, error) {
	var grant gen.CustomerCouponGrant
	err := tx.Model(&gen.CustomerCouponGrant{}).Select("customer_coupon_grants.*").
		Joins("JOIN customer_coupon_templates ON customer_coupon_templates.id = customer_coupon_grants.template_id").
		Joins("JOIN customer_members ON customer_members.id = customer_coupon_grants.member_id").
		Where("customer_coupon_templates.organization_id = ? AND customer_coupon_templates.issuer_scope = ? AND customer_members.organization_id = ?", hqID, gen.CouponIssuerScopeHeadquarters, hqID).
		Where("customer_coupon_grants.template_id = ? AND customer_coupon_grants.member_id = ? AND customer_coupon_grants.request_key = ?", templateID, memberID, key).
		Take(&grant).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &grant, nil
}

func (s *Service) auditCouponGrant(tx *gorm.DB, principal *auth.WorkspacePrincipal, grantID, operation, reasonCode string) error {
	return s.audit.Write(tx, audit.Entry{
		ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
		OrganizationID: principal.OrganizationID, Action: "hqCustomerCoupon:" + operation,
		ResourceType: "customerCouponGrant", ResourceID: grantID, ResultCode: "SUCCESS",
		Metadata: audit.MetadataForPrincipal(principal, audit.Metadata{ReasonCode: reasonCode}),
	})
}
