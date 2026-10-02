/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package authentication

import "time"

// TemporaryPasswordLifetime 限定一次性临时密码的最长有效期。
const TemporaryPasswordLifetime = 24 * time.Hour

// AccountCredential 隔离保存不可暴露给 GraphQL 的认证秘密与锁定状态。
type AccountCredential struct {
	AccountID                  string     `gorm:"type:varchar(36);primaryKey"`
	PasswordHash               string     `gorm:"type:varchar(255);not null"`
	FailedLoginCount           int        `gorm:"not null;default:0"`
	LockedUntil                *time.Time `gorm:"default:null;index"`
	TemporaryPasswordExpiresAt *time.Time `gorm:"default:null;index"`
	PasswordChangedAt          time.Time  `gorm:"not null"`
	UpdatedAt                  time.Time  `gorm:"not null"`
}

// NewTemporaryPasswordExpiry 返回新临时密码的服务端失效时间。
func NewTemporaryPasswordExpiry(now time.Time) *time.Time {
	expiresAt := now.Add(TemporaryPasswordLifetime)
	return &expiresAt
}
