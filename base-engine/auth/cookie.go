/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package auth

import (
	"net/http"
	"time"

	"base-engine/config"
)

const SessionCookieName = "base_admin_session"

// SetSessionCookie 写入仅服务端可读的会话 Cookie。
func SetSessionCookie(writer http.ResponseWriter, token string, cfg config.SecurityConfig) {
	http.SetCookie(writer, &http.Cookie{
		Name: SessionCookieName, Value: token, Path: "/", HttpOnly: true,
		Secure: cfg.CookieSecure, SameSite: http.SameSiteLaxMode,
		MaxAge:  int(cfg.SessionDuration.Seconds()),
		Expires: time.Now().Add(cfg.SessionDuration),
	})
}

// ClearSessionCookie 立即删除当前会话 Cookie。
func ClearSessionCookie(writer http.ResponseWriter, cfg config.SecurityConfig) {
	http.SetCookie(writer, &http.Cookie{
		Name: SessionCookieName, Path: "/", HttpOnly: true,
		Secure: cfg.CookieSecure, SameSite: http.SameSiteLaxMode,
		MaxAge: -1, Expires: time.Unix(1, 0),
	})
}
