/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package organization

import (
	"context"
	"strings"
	"unicode/utf8"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/utils"
	"gorm.io/gorm"
)

func normalizeProvisionInput(input ProvisionInput) (ProvisionInput, error) {
	input.Code = strings.TrimSpace(input.Code)
	input.Name = strings.TrimSpace(input.Name)
	input.OwnerPhone = strings.TrimSpace(input.OwnerPhone)
	input.OwnerDisplayName = strings.TrimSpace(input.OwnerDisplayName)
	input.OwnerEmail = normalizeOptional(input.OwnerEmail)
	valid := bounded(input.Code, 2, 32) && bounded(input.Name, 1, 128) &&
		bounded(input.OwnerDisplayName, 1, 64) && utils.MatchesRule("phone", input.OwnerPhone)
	if !valid || !validEmail(input.OwnerEmail) {
		return ProvisionInput{}, auth.NewError(auth.CodeValidationFailed)
	}
	return input, nil
}

func normalizeReason(reason string, maxLength int) (string, error) {
	reason = strings.TrimSpace(reason)
	if !bounded(reason, 1, maxLength) {
		return "", auth.NewError(auth.CodeValidationFailed)
	}
	return reason, nil
}

func organizationCodeExists(ctx context.Context, db *gorm.DB, code string) (bool, error) {
	var count int64
	err := db.WithContext(ctx).Model(&gen.Organization{}).Where("code = ?", code).Count(&count).Error
	return count > 0, err
}

func normalizeOptional(value *string) *string {
	if value == nil {
		return nil
	}
	normalized := strings.TrimSpace(*value)
	if normalized == "" {
		return nil
	}
	return &normalized
}

func validEmail(value *string) bool {
	return value == nil || (len(*value) <= 128 && utils.MatchesRule("email", *value))
}

func bounded(value string, minimum, maximum int) bool {
	count := utf8.RuneCountInString(value)
	return count >= minimum && count <= maximum
}
