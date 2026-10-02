package customer

import (
	"context"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
)

type PointEntryPage struct {
	Data  []gen.CustomerPointEntry
	Total int64
}

func (s *Service) ListPointEntries(ctx context.Context, principal *auth.WorkspacePrincipal, memberID string, page, perPage int64) (PointEntryPage, error) {
	hqID, err := headquartersID(principal, "hqCustomerPoints:read", authorization.AccessRead)
	if err != nil {
		return PointEntryPage{}, err
	}
	if page < 1 || perPage < 1 || perPage > 50 || page > 1_000_000 {
		return PointEntryPage{}, auth.NewError(auth.CodeValidationFailed)
	}
	var member gen.CustomerMember
	err = s.db.WithContext(ctx).Where("id = ? AND organization_id = ?", memberID, hqID).First(&member).Error
	if err == gorm.ErrRecordNotFound {
		return PointEntryPage{}, auth.NewError(auth.CodePermissionDenied)
	}
	if err != nil {
		return PointEntryPage{}, err
	}
	query := s.db.WithContext(ctx).Model(&gen.CustomerPointEntry{}).
		Where("member_id = ? AND source_organization_id = ?", memberID, hqID)
	var result PointEntryPage
	if err := query.Count(&result.Total).Error; err != nil {
		return PointEntryPage{}, err
	}
	err = query.Order("created_at DESC, id DESC").
		Offset(int((page - 1) * perPage)).Limit(int(perPage)).Find(&result.Data).Error
	return result, err
}
