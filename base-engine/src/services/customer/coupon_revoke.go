package customer

import (
	"context"
	"strings"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *Service) RevokeCoupon(ctx context.Context, principal *auth.WorkspacePrincipal, grantID, reasonCode string) (*gen.CustomerCouponGrant, error) {
	hqID, err := headquartersID(principal, "hqCustomerCoupon:revoke", authorization.AccessUpdate)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(reasonCode) == "" || len(reasonCode) > 64 {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	var grant gen.CustomerCouponGrant
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&grant, "id = ?", grantID).Error; err != nil {
			return err
		}
		var template gen.CustomerCouponTemplate
		err := tx.Where("id = ? AND organization_id = ? AND issuer_scope = ?", grant.TemplateID, hqID, gen.CouponIssuerScopeHeadquarters).First(&template).Error
		if err == gorm.ErrRecordNotFound {
			return auth.NewError(auth.CodePermissionDenied)
		}
		if err != nil {
			return err
		}
		if grant.Status == gen.CustomerCouponGrantStatusRevoked {
			return nil
		}
		nowMillis := s.now().UnixMilli()
		if err := tx.Model(&grant).Updates(map[string]any{
			"status": gen.CustomerCouponGrantStatusRevoked, "revoked_at": nowMillis,
		}).Error; err != nil {
			return err
		}
		grant.Status = gen.CustomerCouponGrantStatusRevoked
		grant.RevokedAt = &nowMillis
		return s.auditCouponGrant(tx, principal, grant.ID, "revoke", reasonCode)
	})
	return &grant, err
}
