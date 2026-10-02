package customer

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *Service) GrantStoreCoupon(ctx context.Context, principal *auth.WorkspacePrincipal, storeID, templateID, memberID, requestKey string) (*gen.CustomerCouponGrant, error) {
	if strings.TrimSpace(requestKey) == "" || len(requestKey) > 128 {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	var result *gen.CustomerCouponGrant
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		grant, err := s.grantStoreCouponTx(tx, principal, storeID, templateID, memberID, requestKey)
		result = grant
		return err
	})
	if mysqlErrorNumber(err) == 1062 {
		replay, findErr := findStoreCouponReplay(s.db.WithContext(ctx), storeCouponRequestDigest(storeID, requestKey), templateID, memberID, requestKey)
		if findErr != nil {
			return nil, findErr
		}
		if replay != nil {
			return replay, nil
		}
	}
	return result, err
}

func (s *Service) grantStoreCouponTx(tx *gorm.DB, principal *auth.WorkspacePrincipal, storeID, templateID, memberID, requestKey string) (*gen.CustomerCouponGrant, error) {
	store, err := couponStoreScope(tx, principal, storeID, "franchiseCoupon:grant", authorization.AccessCreate)
	if err != nil {
		return nil, err
	}
	digest := storeCouponRequestDigest(store.ID, requestKey)
	grant, _, err := s.createCouponGrantInTransaction(tx, couponGrantInput{
		TemplateID: templateID, MemberID: memberID, RequestKey: requestKey,
		ExpectedStoreID: &store.ID, IssuerRequestDigest: &digest, Principal: principal,
		AuditOrganizationID: &store.OrganizationID, AuditStoreID: &store.ID, NowMillis: s.now().UnixMilli(),
	})
	return grant, err
}

func findStoreCouponReplay(tx *gorm.DB, digest, templateID, memberID, requestKey string) (*gen.CustomerCouponGrant, error) {
	var grant gen.CustomerCouponGrant
	err := tx.Where("issuer_request_digest = ?", digest).First(&grant).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if grant.TemplateID != templateID || grant.MemberID != memberID || grant.RequestKey != requestKey {
		return nil, auth.NewError(auth.CodeConflict)
	}
	return &grant, nil
}

func storeCouponRequestDigest(storeID, requestKey string) string {
	sum := sha256.Sum256([]byte("STORE\x00" + storeID + "\x00" + requestKey))
	return hex.EncodeToString(sum[:])
}

func (s *Service) RevokeStoreCoupon(ctx context.Context, principal *auth.WorkspacePrincipal, storeID, grantID, reasonCode string) (*gen.CustomerCouponGrant, error) {
	if strings.TrimSpace(reasonCode) == "" || len(reasonCode) > 64 {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	var grant gen.CustomerCouponGrant
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		store, err := couponStoreScope(tx, principal, storeID, "franchiseCoupon:revoke", authorization.AccessUpdate)
		if err != nil {
			return err
		}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&grant, "id = ?", grantID).Error; err != nil {
			return auth.NewError(auth.CodePermissionDenied)
		}
		var template gen.CustomerCouponTemplate
		if err := tx.Where("id = ? AND applicable_store_id = ? AND issuer_scope = ?", grant.TemplateID, store.ID, gen.CouponIssuerScopeStore).First(&template).Error; err != nil {
			return auth.NewError(auth.CodePermissionDenied)
		}
		if grant.Status == gen.CustomerCouponGrantStatusRevoked {
			return nil
		}
		if grant.Status != gen.CustomerCouponGrantStatusPendingActivation {
			return auth.NewError(auth.CodeConflict)
		}
		nowMillis := s.now().UnixMilli()
		if err := tx.Model(&grant).Updates(map[string]any{"status": gen.CustomerCouponGrantStatusRevoked, "revoked_at": nowMillis}).Error; err != nil {
			return err
		}
		grant.Status = gen.CustomerCouponGrantStatusRevoked
		grant.RevokedAt = &nowMillis
		return s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
			OrganizationID: &store.OrganizationID, StoreID: &store.ID, Action: "franchiseCoupon:revoke",
			ResourceType: "customerCouponGrant", ResourceID: grant.ID, ResultCode: "SUCCESS", Metadata: audit.Metadata{ReasonCode: reasonCode}})
	})
	return &grant, err
}
