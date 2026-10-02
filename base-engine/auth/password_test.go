/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package auth

import (
	"errors"
	"regexp"
	"strings"
	"testing"
)

func TestGenerateTemporaryPasswordUsesEightLettersAndDigits(t *testing.T) {
	letters := regexp.MustCompile(`[A-Za-z]`)
	digits := regexp.MustCompile(`[0-9]`)
	allowed := regexp.MustCompile(`^[A-Za-z0-9]{8}$`)
	for range 100 {
		password, err := GenerateTemporaryPassword()
		if err != nil {
			t.Fatal(err)
		}
		if !allowed.MatchString(password) || !letters.MatchString(password) || !digits.MatchString(password) {
			t.Fatalf("temporary password violates eight-character alphanumeric policy: %q", password)
		}
		if err := ValidatePassword(password); err != nil {
			t.Fatalf("generated password rejected by credential policy: %v", err)
		}
	}
}

// TestValidatePassword 验证后台账号密码的字符长度边界。
func TestValidatePassword(t *testing.T) {
	valid := []string{
		"password",
		strings.Repeat("密", 20),
	}
	for _, value := range valid {
		if err := ValidatePassword(value); err != nil {
			t.Fatalf("valid password %q rejected: %v", value, err)
		}
	}
	invalid := []string{
		"1234567",
		strings.Repeat("密", 21),
	}
	for _, value := range invalid {
		if !errors.Is(ValidatePassword(value), ErrWeakPassword) {
			t.Errorf("invalid password accepted: %q", value)
		}
	}
}

// TestPasswordHashIsSaltedAndVerifiable 验证 bcrypt 使用随机盐且只接受正确明文。
func TestPasswordHashIsSaltedAndVerifiable(t *testing.T) {
	first, err := HashPassword("Correct-Horse-42")
	if err != nil {
		t.Fatal(err)
	}
	second, err := HashPassword("Correct-Horse-42")
	if err != nil {
		t.Fatal(err)
	}
	if first == second {
		t.Fatal("bcrypt hashes must use distinct salts")
	}
	if err := VerifyPassword(first, "Correct-Horse-42"); err != nil {
		t.Fatalf("valid password rejected: %v", err)
	}
	if err := VerifyPassword(first, "Wrong-Horse-42"); err == nil {
		t.Fatal("wrong password accepted")
	}
}
