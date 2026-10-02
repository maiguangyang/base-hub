package customer

import (
	"context"
	"strings"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
)

type MemberView struct {
	ID, MemberNumber, PhoneMasked string
	Status                        gen.CustomerMemberStatus
	PointsBalance                 int64
	PointsFrozen                  bool
	PendingCouponCount            int64
	CreatedAt                     time.Time
}

type MemberSearch struct {
	Phone   string
	Status  *gen.CustomerMemberStatus
	Page    int64
	PerPage int64
}

type MemberPage struct {
	Data        []MemberView
	Total       int64
	CurrentPage int64
	PerPage     int64
}

func viewMember(member *gen.CustomerMember) MemberView {
	view := MemberView{
		ID: member.ID, MemberNumber: member.MemberNumber,
		Status: member.Status, PointsBalance: member.PointsBalance,
		PointsFrozen: member.PointsFrozen, CreatedAt: time.UnixMilli(member.CreatedAt),
	}
	if member.Phone != nil {
		view.PhoneMasked = MaskCNPhone(*member.Phone)
	}
	return view
}

func (s *Service) SearchMembers(ctx context.Context, principal *auth.WorkspacePrincipal, search MemberSearch) (*MemberPage, error) {
	hqID, err := headquartersID(principal, "hqCustomer:read", authorization.AccessRead)
	if err != nil {
		return nil, err
	}
	query, err := memberSearchQuery(s.db.WithContext(ctx), hqID, search)
	if err != nil {
		return nil, err
	}
	result := &MemberPage{CurrentPage: search.Page, PerPage: search.PerPage}
	if err := query.Count(&result.Total).Error; err != nil {
		return nil, err
	}
	var records []gen.CustomerMember
	if err := query.Order("created_at DESC").Order("id DESC").
		Limit(int(search.PerPage)).Offset(int((search.Page - 1) * search.PerPage)).
		Find(&records).Error; err != nil {
		return nil, err
	}
	pending, err := s.pendingCouponsForMembers(ctx, records)
	if err != nil {
		return nil, err
	}
	result.Data = make([]MemberView, 0, len(records))
	for index := range records {
		view := viewMember(&records[index])
		view.PendingCouponCount = pending[view.ID]
		result.Data = append(result.Data, view)
	}
	return result, nil
}

func (s *Service) pendingCouponsForMembers(ctx context.Context, members []gen.CustomerMember) (map[string]int64, error) {
	result := make(map[string]int64, len(members))
	if len(members) == 0 {
		return result, nil
	}
	ids := make([]string, 0, len(members))
	for _, member := range members {
		ids = append(ids, member.ID)
	}
	var counts []struct {
		MemberID string
		Pending  int64
	}
	err := s.db.WithContext(ctx).Model(&gen.CustomerCouponGrant{}).
		Select("member_id, COUNT(*) AS pending").
		Where("member_id IN ? AND status = ?", ids, gen.CustomerCouponGrantStatusPendingActivation).
		Group("member_id").Scan(&counts).Error
	for _, count := range counts {
		result[count.MemberID] = count.Pending
	}
	return result, err
}

func memberSearchQuery(db *gorm.DB, hqID string, search MemberSearch) (*gorm.DB, error) {
	if search.Page < 1 || search.PerPage < 1 || search.PerPage > 50 ||
		(search.Status != nil && !search.Status.IsValid()) {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	query := db.Model(&gen.CustomerMember{}).Where(
		"organization_id = ? AND (is_delete IS NULL OR is_delete = 1)", hqID,
	)
	if strings.TrimSpace(search.Phone) != "" {
		prefix, err := NormalizeCNPhonePrefix(search.Phone)
		if err != nil {
			return nil, err
		}
		if len(prefix) == 14 {
			query = query.Where("phone = ?", prefix)
		} else {
			query = query.Where("phone LIKE ?", prefix+"%")
		}
	}
	if search.Status != nil {
		query = query.Where("status = ?", *search.Status)
	}
	return query, nil
}

func (s *Service) GetMember(ctx context.Context, principal *auth.WorkspacePrincipal, memberID string) (*MemberView, error) {
	hqID, err := headquartersID(principal, "hqCustomer:read", authorization.AccessRead)
	if err != nil {
		return nil, err
	}
	var member gen.CustomerMember
	err = s.db.WithContext(ctx).Where("organization_id = ? AND (is_delete IS NULL OR is_delete = 1)", hqID).
		First(&member, "id = ?", memberID).Error
	if err == gorm.ErrRecordNotFound {
		return nil, auth.NewError(auth.CodePermissionDenied)
	}
	if err != nil {
		return nil, err
	}
	view := viewMember(&member)
	err = s.db.WithContext(ctx).Model(&gen.CustomerCouponGrant{}).Where(
		"member_id = ? AND status = ?", memberID, gen.CustomerCouponGrantStatusPendingActivation,
	).Count(&view.PendingCouponCount).Error
	return &view, err
}

func (s *Service) SensitivePhone(ctx context.Context, principal *auth.WorkspacePrincipal, memberID string) (string, error) {
	hqID, err := headquartersID(principal, "customer:read_sensitive", authorization.AccessRead)
	if err != nil {
		return "", err
	}
	var member gen.CustomerMember
	err = s.db.WithContext(ctx).Where("organization_id = ? AND (is_delete IS NULL OR is_delete = 1)", hqID).
		First(&member, "id = ?", memberID).Error
	if err == gorm.ErrRecordNotFound {
		return "", auth.NewError(auth.CodePermissionDenied)
	}
	if err != nil {
		return "", err
	}
	if member.Phone == nil {
		return "", auth.NewError(auth.CodePermissionDenied)
	}
	return *member.Phone, nil
}
