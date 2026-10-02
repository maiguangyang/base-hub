/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package config

import (
	"errors"
	"strings"
	"testing"
	"time"
)

// TestLoadSecurityConfig 验证安全配置的必填值、默认值和来源解析。
func TestLoadSecurityConfig(t *testing.T) {
	setValidSecurityEnv(t)
	cfg, err := LoadSecurityConfig()
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.OriginAllowed("https://admin.example.com") || cfg.OriginAllowed("https://evil.example.com") {
		t.Fatalf("unexpected origins: %#v", cfg.AllowedOrigins)
	}
	if cfg.LoginIPAttemptsPerMinute != 20 || cfg.LoginLockThreshold != 5 {
		t.Fatalf("unexpected defaults: %#v", cfg)
	}
	if cfg.LoginLockDuration != 15*time.Minute || cfg.SessionDuration != 30*24*time.Hour {
		t.Fatalf("unexpected durations: %#v", cfg)
	}
}

// TestSecurityConfigAllowsAnyOriginWithWildcard 验证通配配置接受任意非空来源。
func TestSecurityConfigAllowsAnyOriginWithWildcard(t *testing.T) {
	cfg := SecurityConfig{AllowedOrigins: map[string]struct{}{"*": {}}}
	for _, origin := range []string{"http://localhost:4322", "https://admin.example.com"} {
		if !cfg.OriginAllowed(origin) {
			t.Fatalf("origin %q was rejected", origin)
		}
	}
}

// TestLoadSecurityConfigDefaultSigningKey 验证缺省签名密钥时自动采用各角色独立固定密钥。
func TestLoadSecurityConfigDefaultSigningKey(t *testing.T) {
	setValidSecurityEnv(t)
	t.Setenv("AUTH_SIGNING_KEY", "")
	t.Setenv("AUTH_ADMIN_SIGNING_KEY", "")
	t.Setenv("AUTH_FRANCHISE_SIGNING_KEY", "")
	t.Setenv("AUTH_CLIENT_SIGNING_KEY", "")
	t.Setenv("AUTH_TOKEN_ISSUER", "")
	cfg, err := LoadSecurityConfig()
	if err != nil {
		t.Fatalf("expected fallback to default signing key, got error: %v", err)
	}
	if string(cfg.AdminSigningKey) != DefaultAdminSigningKey {
		t.Fatalf("got admin key %s, want %s", string(cfg.AdminSigningKey), DefaultAdminSigningKey)
	}
	if string(cfg.FranchiseSigningKey) != DefaultFranchiseSigningKey {
		t.Fatalf("got franchise key %s, want %s", string(cfg.FranchiseSigningKey), DefaultFranchiseSigningKey)
	}
	if string(cfg.ClientSigningKey) != DefaultClientSigningKey {
		t.Fatalf("got client key %s, want %s", string(cfg.ClientSigningKey), DefaultClientSigningKey)
	}
	if string(cfg.SigningKeyForWorkspace("HEADQUARTERS")) != DefaultAdminSigningKey {
		t.Fatalf("unexpected hq key")
	}
	if string(cfg.SigningKeyForWorkspace("FRANCHISE")) != DefaultFranchiseSigningKey {
		t.Fatalf("unexpected franchise key")
	}
	if string(cfg.SigningKeyForWorkspace("CLIENT")) != DefaultClientSigningKey {
		t.Fatalf("unexpected client key")
	}
	if cfg.TokenIssuer != DefaultTokenIssuer {
		t.Fatalf("got issuer %s, want %s", cfg.TokenIssuer, DefaultTokenIssuer)
	}
}

// TestLoadSecurityConfigRejectsUnsafeValues 验证短密钥和多副本部署均失败关闭。
func TestLoadSecurityConfigRejectsUnsafeValues(t *testing.T) {
	tests := []struct {
		name, key, replicas string
		want                error
	}{
		{name: "short key", key: "short", replicas: "1", want: ErrInvalidSecurityConfig},
		{name: "multiple replicas", key: strings.Repeat("k", 32), replicas: "2", want: ErrSessionPubSubRequired},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setValidSecurityEnv(t)
			t.Setenv("AUTH_SIGNING_KEY", test.key)
			t.Setenv("ENGINE_REPLICA_COUNT", test.replicas)
			_, err := LoadSecurityConfig()
			if !errors.Is(err, test.want) {
				t.Fatalf("error = %v, want %v", err, test.want)
			}
		})
	}
	t.Run("empty origins", func(t *testing.T) {
		setValidSecurityEnv(t)
		t.Setenv("ALLOWED_ORIGINS", "")
		_, err := LoadSecurityConfig()
		if !errors.Is(err, ErrInvalidSecurityConfig) {
			t.Fatalf("error = %v, want %v", err, ErrInvalidSecurityConfig)
		}
	})
}

// setValidSecurityEnv 注入一组可通过校验的测试环境变量。
func setValidSecurityEnv(t *testing.T) {
	t.Helper()
	t.Setenv("AUTH_SIGNING_KEY", strings.Repeat("k", 32))
	t.Setenv("AUTH_TOKEN_ISSUER", "base-engine")
	t.Setenv("AUTH_COOKIE_SECURE", "true")
	t.Setenv("ALLOWED_ORIGINS", "https://admin.example.com, https://ops.example.com")
	t.Setenv("ENGINE_REPLICA_COUNT", "1")
	t.Setenv("LOGIN_IP_ATTEMPTS_PER_MIN", "")
	t.Setenv("LOGIN_LOCK_THRESHOLD", "")
	t.Setenv("LOGIN_LOCK_MINUTES", "")
}
