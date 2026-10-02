package customer

import (
	"context"
	"strings"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CreateMemberInput struct {
	Phone, RequestKey string
}

func (s *Service) CreateMember(ctx context.Context, principal *auth.WorkspacePrincipal, input CreateMemberInput) (*gen.CustomerMember, error) {
	hqID, err := headquartersID(principal, "hqCustomer:create", authorization.AccessCreate)
	if err != nil {
		return nil, err
	}
	if !validMobileDigits(input.Phone, 11, 11) || !validMemberInput(input) {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	phone, err := NormalizeCNPhone(input.Phone)
	if err != nil {
		return nil, err
	}
	var result *gen.CustomerMember
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var createErr error
		result, createErr = s.createMemberInTransaction(tx, principal, hqID, phone, input)
		return createErr
	})
	if auth.ErrorCode(err) == auth.CodeConflict {
		replayed, findErr := findMemberByRequest(s.db.WithContext(ctx), hqID, input.RequestKey)
		if findErr != nil {
			return nil, findErr
		}
		if replayed != nil && sameCreateIntent(replayed, phone) {
			return replayed, nil
		}
	}
	return result, err
}

func (s *Service) createMemberInTransaction(tx *gorm.DB, principal *auth.WorkspacePrincipal, hqID, phone string, input CreateMemberInput) (*gen.CustomerMember, error) {
	replayed, err := findMemberByRequest(tx, hqID, input.RequestKey)
	if err != nil {
		return nil, err
	}
	if replayed != nil {
		if !sameCreateIntent(replayed, phone) {
			return nil, auth.NewError(auth.CodeConflict)
		}
		return replayed, nil
	}
	id := uuid.Must(uuid.NewV4()).String()
	member := &gen.CustomerMember{
		ID: id, MemberNumber: "CM" + strings.ReplaceAll(id, "-", "")[:30],
		CreatedAt:  s.now().UnixMilli(),
		RequestKey: input.RequestKey, Phone: &phone,
		Status: gen.CustomerMemberStatusActive, OrganizationID: hqID,
	}
	inserted := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(member)
	if inserted.Error != nil {
		return nil, inserted.Error
	}
	if inserted.RowsAffected != 1 {
		return nil, auth.NewError(auth.CodeConflict)
	}
	if err := scheduleMemberCouponCatchup(tx, member, member.CreatedAt); err != nil {
		return nil, err
	}
	return member, s.auditMember(tx, principal, "hqCustomer:create", member.ID, "")
}

func sameCreateIntent(member *gen.CustomerMember, phone string) bool {
	return member.Phone != nil && *member.Phone == phone
}

func validMemberInput(input CreateMemberInput) bool {
	input.RequestKey = strings.TrimSpace(input.RequestKey)
	return input.RequestKey != "" && len(input.RequestKey) <= 128
}

func findMemberByRequest(tx *gorm.DB, hqID, key string) (*gen.CustomerMember, error) {
	var member gen.CustomerMember
	err := tx.Where("organization_id = ? AND request_key = ?", hqID, key).First(&member).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &member, nil
}
