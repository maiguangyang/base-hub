/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package membership

import (
	"strings"
	"unicode/utf8"

	"base-engine/auth"
	"base-engine/utils"
)

func normalizeInviteInput(input InviteInput) (InviteInput, error) {
	input.Phone = strings.TrimSpace(input.Phone)
	input.DisplayName = strings.TrimSpace(input.DisplayName)
	input.Email = normalizeOptional(input.Email)
	valid := utils.MatchesRule("phone", input.Phone) &&
		utf8.RuneCountInString(input.DisplayName) >= 1 &&
		utf8.RuneCountInString(input.DisplayName) <= 64
	if !valid || !validEmail(input.Email) {
		return InviteInput{}, auth.NewError(auth.CodeValidationFailed)
	}
	return input, nil
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
