package customer

import (
	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

func prepareAutomaticCouponGrant(tx *gorm.DB, input couponGrantInput, template *gen.CustomerCouponTemplate, member *gen.CustomerMember) (*gen.CustomerCouponGrant, error) {
	replay, err := findCouponGrantReplay(tx, input, template.OrganizationID)
	if err != nil || replay != nil {
		return replay, err
	}
	if err := enforceAutomaticCouponMemberPriority(tx, input, template, member); err != nil {
		return nil, err
	}
	return nil, nil
}

func enforceAutomaticCouponMemberPriority(tx *gorm.DB, input couponGrantInput, template *gen.CustomerCouponTemplate, member *gen.CustomerMember) error {
	if !input.EnforceMemberPriority {
		return nil
	}
	mustYield, err := laterMemberWouldConsumeEarlierCouponCapacity(tx, template, member)
	if err != nil {
		return err
	}
	if mustYield {
		return auth.NewError(auth.CodeConflict)
	}
	return nil
}

func laterMemberWouldConsumeEarlierCouponCapacity(tx *gorm.DB, template *gen.CustomerCouponTemplate, member *gen.CustomerMember) (bool, error) {
	if template.TotalIssueLimit == 0 {
		return false, nil
	}
	remaining := template.TotalIssueLimit - template.IssuedCount
	if remaining <= 0 {
		return false, nil
	}
	grantedMemberIDs := tx.Model(&gen.CustomerCouponGrant{}).Select("member_id").Where("template_id = ?", template.ID)
	var earlierCount int64
	err := tx.Model(&gen.CustomerMember{}).
		Where("organization_id = ? AND status = ?", template.OrganizationID, gen.CustomerMemberStatusActive).
		Where("created_at < ? OR (created_at = ? AND id < ?)", member.CreatedAt, member.CreatedAt, member.ID).
		Where("id NOT IN (?)", grantedMemberIDs).
		Count(&earlierCount).Error
	return earlierCount >= remaining, err
}
