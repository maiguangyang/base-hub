/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package auth

import (
	"crypto/rand"
	"errors"
	"math/big"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

var ErrWeakPassword = errors.New("PASSWORD_WEAK")

const passwordHashCost = 12

// ValidatePassword 校验密码字符长度。
func ValidatePassword(value string) error {
	if !utf8.ValidString(value) {
		return ErrWeakPassword
	}
	length := utf8.RuneCountInString(value)
	if length < 8 || length > 20 {
		return ErrWeakPassword
	}
	return nil
}

// HashPassword 校验后使用 bcrypt cost 12 生成不可逆摘要。
func HashPassword(value string) (string, error) {
	if err := ValidatePassword(value); err != nil {
		return "", err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(value), passwordHashCost)
	return string(hash), err
}

// VerifyPassword 对照 bcrypt 摘要验证明文密码。
func VerifyPassword(hash, value string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(value))
}

// GenerateTemporaryPassword 创建 8 位英数字一次性密码。
func GenerateTemporaryPassword() (string, error) {
	const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	limit := big.NewInt(int64(len(alphabet)))
	for {
		password := make([]byte, 8)
		hasLetter, hasDigit := false, false
		for index := range password {
			choice, err := rand.Int(rand.Reader, limit)
			if err != nil {
				return "", err
			}
			password[index] = alphabet[choice.Int64()]
			hasDigit = hasDigit || choice.Int64() < 10
			hasLetter = hasLetter || choice.Int64() >= 10
		}
		if hasLetter && hasDigit {
			return string(password), nil
		}
	}
}
