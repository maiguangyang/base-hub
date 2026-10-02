package customer

import (
	"strings"

	"base-engine/auth"
)

// NormalizeCNPhone stores mainland mobile numbers in one canonical format.
func NormalizeCNPhone(input string) (string, error) {
	cleaned := compactCNPhone(input)
	if !validMobileDigits(cleaned, 11, 11) {
		return "", auth.NewError(auth.CodeValidationFailed)
	}
	return "+86" + cleaned, nil
}

// MaskCNPhone keeps the country prefix, mobile prefix, and final four digits.
func MaskCNPhone(canonical string) string {
	if len(canonical) != 14 || !strings.HasPrefix(canonical, "+86") {
		return ""
	}
	return canonical[:6] + "****" + canonical[10:]
}

func NormalizeCNPhonePrefix(input string) (string, error) {
	cleaned := compactCNPhone(input)
	if !validMobileDigits(cleaned, 7, 11) {
		return "", auth.NewError(auth.CodeValidationFailed)
	}
	return "+86" + cleaned, nil
}

func compactCNPhone(input string) string {
	cleaned := strings.Map(func(r rune) rune {
		if r == ' ' || r == '-' {
			return -1
		}
		return r
	}, input)
	cleaned = strings.TrimPrefix(cleaned, "+86")
	if strings.HasPrefix(cleaned, "86") {
		cleaned = cleaned[2:]
	}
	return cleaned
}

func validMobileDigits(value string, minimum, maximum int) bool {
	if len(value) < minimum || len(value) > maximum ||
		value[0] != '1' || value[1] < '3' || value[1] > '9' {
		return false
	}
	for _, digit := range value {
		if digit < '0' || digit > '9' {
			return false
		}
	}
	return true
}
