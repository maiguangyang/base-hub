/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"base-engine/config"
)

// WorkspaceType 是认证协议使用的工作台类型，避免认证层依赖生成模型。
type WorkspaceType string

const (
	WorkspaceTypeDiscovery    WorkspaceType = "DISCOVERY"
	WorkspaceTypeHeadquarters WorkspaceType = "HEADQUARTERS"
	WorkspaceTypeFranchise    WorkspaceType = "FRANCHISE"
	WorkspaceTypeClient       WorkspaceType = "CLIENT"
)

// SessionClaims 仅保存可用于重新加载权限的稳定标识。
type SessionClaims struct {
	SessionID         string        `json:"sid"`
	AccountID         string        `json:"aid"`
	WorkspaceType     WorkspaceType `json:"wst"`
	OrganizationID    *string       `json:"oid,omitempty"`
	CredentialVersion int64         `json:"cv"`
	jwt.RegisteredClaims
}

// SignSessionClaims 为会话签发 HS256 token。
func SignSessionClaims(cfg config.SecurityConfig, claims SessionClaims) (string, error) {
	return signSessionClaimsAt(cfg, claims, time.Now())
}

func signSessionClaimsAt(cfg config.SecurityConfig, claims SessionClaims, now time.Time) (string, error) {
	key := cfg.SigningKeyForWorkspace(string(claims.WorkspaceType))
	if len(key) < 32 || cfg.TokenIssuer == "" || cfg.SessionDuration <= 0 {
		return "", config.ErrInvalidSecurityConfig
	}
	claims.Issuer = cfg.TokenIssuer
	claims.IssuedAt = jwt.NewNumericDate(now)
	claims.NotBefore = jwt.NewNumericDate(now)
	claims.ExpiresAt = jwt.NewNumericDate(now.Add(cfg.SessionDuration))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(key)
}

// ParseSessionClaims 严格校验算法、签发者与过期时间。
func ParseSessionClaims(cfg config.SecurityConfig, encoded string) (*SessionClaims, error) {
	encoded = strings.TrimSpace(strings.TrimPrefix(encoded, "Bearer "))
	claims := &SessionClaims{}
	token, err := jwt.ParseWithClaims(
		encoded,
		claims,
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("INVALID_SIGNING_METHOD")
			}
			return cfg.SigningKeyForWorkspace(string(claims.WorkspaceType)), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(cfg.TokenIssuer),
		jwt.WithExpirationRequired(),
	)
	if err != nil || !token.Valid {
		return nil, errors.Join(ErrAuthRequired, err)
	}
	return claims, nil
}
