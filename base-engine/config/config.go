/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type key int

const KeyAppSecret key = iota

const (
	defaultLoginIPAttempts = 20
	defaultLockThreshold   = 5
	defaultLockMinutes     = 15
	defaultSessionHours    = 30 * 24

	// DefaultAdminSigningKey 为后台管理员固定签名密钥。
	DefaultAdminSigningKey = "6f91ccdb1e2a476272c0a308b2f7cd9ce3246fc31fce69269c7ceff092c6f8ac"
	// DefaultFranchiseSigningKey 为加盟店固定签名密钥。
	DefaultFranchiseSigningKey = "0b47cd7336842c24b492869501f75019b6dfd348ccf4f1b64ee426b4f63938f0"
	// DefaultClientSigningKey 为客户端用户固定签名密钥。
	DefaultClientSigningKey = "1a648b482af33db8e8fefe16d2262833cf3bb2fa54448c08804c040c092fb4e3"
	// DefaultAuthSigningKey 兼容保留通用默认密钥（缺省以管理员密钥为基底）。
	DefaultAuthSigningKey = DefaultAdminSigningKey
	// DefaultTokenIssuer 为系统默认的令牌签发者标识。
	DefaultTokenIssuer = "korean-engine"
)

var (
	ErrInvalidSecurityConfig = errors.New("INVALID_SECURITY_CONFIG")
	ErrSessionPubSubRequired = errors.New("SESSION_PUBSUB_REQUIRED")
)

// SecurityConfig 保存认证边界所需的安全配置。
type SecurityConfig struct {
	SigningKey               []byte
	AdminSigningKey          []byte
	FranchiseSigningKey      []byte
	ClientSigningKey         []byte
	TokenIssuer              string
	CookieSecure             bool
	AllowedOrigins           map[string]struct{}
	EngineReplicaCount       int
	LoginIPAttemptsPerMinute int
	LoginLockThreshold       int
	LoginLockDuration        time.Duration
	SessionDuration          time.Duration
}

// LoadSecurityConfig 从环境变量加载配置，支持后台管理员、加盟店、客户端用户独立配置密钥。
func LoadSecurityConfig() (SecurityConfig, error) {
	cfg := SecurityConfig{}
	keys, issuer, err := resolveSigningKeysAndIssuer()
	if err != nil {
		return cfg, err
	}
	cookieSecure, err := requiredBool("AUTH_COOKIE_SECURE")
	if err != nil {
		return cfg, err
	}
	origins, err := requiredOrigins("ALLOWED_ORIGINS")
	if err != nil {
		return cfg, err
	}
	replicas, err := positiveInt("ENGINE_REPLICA_COUNT", 1)
	if err != nil {
		return cfg, err
	}
	if replicas > 1 {
		return cfg, ErrSessionPubSubRequired
	}
	attempts, err := positiveInt("LOGIN_IP_ATTEMPTS_PER_MIN", defaultLoginIPAttempts)
	if err != nil {
		return cfg, err
	}
	threshold, err := positiveInt("LOGIN_LOCK_THRESHOLD", defaultLockThreshold)
	if err != nil {
		return cfg, err
	}
	lockMinutes, err := positiveInt("LOGIN_LOCK_MINUTES", defaultLockMinutes)
	if err != nil {
		return cfg, err
	}
	cfg = SecurityConfig{
		SigningKey:               []byte(keys["general"]),
		AdminSigningKey:          []byte(keys["admin"]),
		FranchiseSigningKey:      []byte(keys["franchise"]),
		ClientSigningKey:         []byte(keys["client"]),
		TokenIssuer:              issuer,
		CookieSecure:             cookieSecure,
		AllowedOrigins:           origins,
		EngineReplicaCount:       replicas,
		LoginIPAttemptsPerMinute: attempts,
		LoginLockThreshold:       threshold,
		LoginLockDuration:        time.Duration(lockMinutes) * time.Minute,
		SessionDuration:          defaultSessionHours * time.Hour,
	}
	return cfg, nil
}

// OriginAllowed 判断请求来源是否符合通配配置或显式白名单。
func (c SecurityConfig) OriginAllowed(origin string) bool {
	origin = strings.TrimSpace(origin)
	if origin == "" {
		return false
	}
	if _, allowAll := c.AllowedOrigins["*"]; allowAll {
		return true
	}
	_, ok := c.AllowedOrigins[origin]
	return ok
}

// SigningKeyForWorkspace 根据工作区或用户类型返回对应的独立签名密钥。
func (c SecurityConfig) SigningKeyForWorkspace(workspace string) []byte {
	switch strings.ToUpper(strings.TrimSpace(workspace)) {
	case "HEADQUARTERS", "ADMIN":
		if len(c.AdminSigningKey) >= 32 {
			return c.AdminSigningKey
		}
	case "FRANCHISE", "STORE", "MERCHANT":
		if len(c.FranchiseSigningKey) >= 32 {
			return c.FranchiseSigningKey
		}
	case "CLIENT", "APP", "KIOSK", "CUSTOMER":
		if len(c.ClientSigningKey) >= 32 {
			return c.ClientSigningKey
		}
	}
	if len(c.SigningKey) >= 32 {
		return c.SigningKey
	}
	return c.AdminSigningKey
}

func resolveSigningKeysAndIssuer() (map[string]string, string, error) {
	keys := make(map[string]string, 4)
	var err error
	if keys["admin"], err = resolveRoleSigningKey("AUTH_ADMIN_SIGNING_KEY", DefaultAdminSigningKey); err != nil {
		return nil, "", err
	}
	if keys["franchise"], err = resolveRoleSigningKey("AUTH_FRANCHISE_SIGNING_KEY", DefaultFranchiseSigningKey); err != nil {
		return nil, "", err
	}
	if keys["client"], err = resolveRoleSigningKey("AUTH_CLIENT_SIGNING_KEY", DefaultClientSigningKey); err != nil {
		return nil, "", err
	}
	legacyKey := strings.TrimSpace(os.Getenv("AUTH_SIGNING_KEY"))
	if legacyKey != "" {
		if len(legacyKey) < 32 {
			return nil, "", invalidConfig("AUTH_SIGNING_KEY")
		}
		keys["general"] = legacyKey
	} else {
		keys["general"] = keys["admin"]
	}
	issuer := strings.TrimSpace(os.Getenv("AUTH_TOKEN_ISSUER"))
	if issuer == "" {
		issuer = DefaultTokenIssuer
	}
	return keys, issuer, nil
}

func resolveRoleSigningKey(envName, defaultKey string) (string, error) {
	key := strings.TrimSpace(os.Getenv(envName))
	if key == "" {
		key = defaultKey
	}
	if len(key) < 32 {
		return "", invalidConfig(envName)
	}
	return key, nil
}

func requiredBool(name string) (bool, error) {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return false, invalidConfig(name)
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false, invalidConfig(name)
	}
	return parsed, nil
}

func requiredOrigins(name string) (map[string]struct{}, error) {
	result := make(map[string]struct{})
	for _, value := range strings.Split(os.Getenv(name), ",") {
		if origin := strings.TrimSpace(value); origin != "" {
			result[origin] = struct{}{}
		}
	}
	if len(result) == 0 {
		return nil, invalidConfig(name)
	}
	return result, nil
}

func positiveInt(name string, fallback int) (int, error) {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 1 {
		return 0, invalidConfig(name)
	}
	return parsed, nil
}

func invalidConfig(name string) error {
	return fmt.Errorf("%w: %s", ErrInvalidSecurityConfig, name)
}
