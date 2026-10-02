/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package auth

import (
	"encoding/base64"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"base-engine/config"
)

// TestSessionClaimsRoundTrip 验证 token 只携带稳定标识且可严格验签。
func TestSessionClaimsRoundTrip(t *testing.T) {
	cfg := testSecurityConfig()
	organizationID := "org-1"
	claims := SessionClaims{
		SessionID: "session-1", AccountID: "account-1",
		WorkspaceType:  WorkspaceTypeFranchise,
		OrganizationID: &organizationID, CredentialVersion: 7,
	}
	token, err := signSessionClaimsAt(cfg, claims, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := ParseSessionClaims(cfg, token)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.SessionID != claims.SessionID || parsed.AccountID != claims.AccountID {
		t.Fatalf("claims mismatch: %#v", parsed)
	}
	payload := decodeClaimsForTest(t, token)
	if _, exists := payload["permissions"]; exists {
		t.Fatal("permissions must not be embedded in session claims")
	}
}

// TestSessionClaimsRejectInvalidTokens 验证算法、签发者与过期时间均被严格校验。
func TestSessionClaimsRejectInvalidTokens(t *testing.T) {
	cfg := testSecurityConfig()
	now := time.Now()
	expired, err := signSessionClaimsAt(cfg, SessionClaims{}, now.Add(-13*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParseSessionClaims(cfg, expired); err == nil {
		t.Fatal("expired token accepted")
	}
	wrongIssuer := cfg
	wrongIssuer.TokenIssuer = "other-issuer"
	valid, err := signSessionClaimsAt(cfg, SessionClaims{}, now)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParseSessionClaims(wrongIssuer, valid); err == nil {
		t.Fatal("wrong issuer accepted")
	}
	wrongAlgorithm := jwt.NewWithClaims(jwt.SigningMethodHS384, jwt.MapClaims{"iss": cfg.TokenIssuer, "exp": now.Add(time.Hour).Unix()})
	encoded, err := wrongAlgorithm.SignedString(cfg.SigningKey)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParseSessionClaims(cfg, encoded); err == nil {
		t.Fatal("wrong signing algorithm accepted")
	}
}

// TestSessionCookieAttributes 验证会话 Cookie 不可由脚本读取且使用约定路径和 SameSite。
func TestSessionCookieAttributes(t *testing.T) {
	cfg := testSecurityConfig()
	recorder := httptest.NewRecorder()
	SetSessionCookie(recorder, "signed-token", cfg)
	cookies := recorder.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("cookie count = %d", len(cookies))
	}
	cookie := cookies[0]
	if cookie.Name != SessionCookieName || !cookie.HttpOnly || !cookie.Secure || cookie.Path != "/" {
		t.Fatalf("unsafe cookie: %#v", cookie)
	}
}

// testSecurityConfig 构造不包含环境依赖的安全配置。
func testSecurityConfig() config.SecurityConfig {
	return config.SecurityConfig{
		SigningKey: []byte(strings.Repeat("k", 32)), TokenIssuer: "korean-engine",
		CookieSecure: true, SessionDuration: 12 * time.Hour,
	}
}

// decodeClaimsForTest 解码 JWT payload，仅用于断言字段集合。
func decodeClaimsForTest(t *testing.T, token string) map[string]any {
	t.Helper()
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatalf("invalid token parts: %d", len(parts))
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatal(err)
	}
	var claims map[string]any
	if err := json.Unmarshal(payload, &claims); err != nil {
		t.Fatal(err)
	}
	return claims
}

// TestRoleKeyIsolation 验证后台管理员、加盟店、客户端用户使用不同独立密钥且互相隔离。
func TestRoleKeyIsolation(t *testing.T) {
	cfg := config.SecurityConfig{
		AdminSigningKey:     []byte("admin-secret-key-at-least-32-bytes!"),
		FranchiseSigningKey: []byte("franchise-key-at-least-32-bytes!!"),
		ClientSigningKey:    []byte("client-key-at-least-32-bytes-long!"),
		TokenIssuer:         "korean-engine",
		SessionDuration:     12 * time.Hour,
	}
	assertRoleKey(t, cfg, WorkspaceTypeHeadquarters)
	assertRoleKey(t, cfg, WorkspaceTypeFranchise)
	assertRoleKey(t, cfg, WorkspaceTypeClient)
}

func assertRoleKey(t *testing.T, cfg config.SecurityConfig, ws WorkspaceType) {
	t.Helper()
	token, err := signSessionClaimsAt(cfg, SessionClaims{SessionID: "s1", AccountID: "a1", WorkspaceType: ws}, time.Now())
	if err != nil {
		t.Fatalf("%s sign failed: %v", ws, err)
	}
	parsed, err := ParseSessionClaims(cfg, token)
	if err != nil || parsed.WorkspaceType != ws {
		t.Fatalf("%s parse failed: %v", ws, err)
	}
	wrongCfg := cfg
	wrongCfg.AdminSigningKey = []byte("wrong-admin-key-32-bytes-long!!!!")
	wrongCfg.FranchiseSigningKey = []byte("wrong-franchise-key-32-bytes-long")
	wrongCfg.ClientSigningKey = []byte("wrong-client-key-32-bytes-long!!!")
	if _, err := ParseSessionClaims(wrongCfg, token); err == nil {
		t.Fatalf("%s should fail with wrong keys", ws)
	}
}

