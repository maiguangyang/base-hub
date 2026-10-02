/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package authentication

import (
	"strings"
	"sync"
	"time"
)

type loginWindow struct {
	startedAt time.Time
	attempts  int
}

// LoginLimiter 提供进程内固定窗口登录 IP 限流。
type LoginLimiter struct {
	mu      sync.Mutex
	limit   int
	window  time.Duration
	entries map[string]loginWindow
	lastGC  time.Time
}

// NewLoginLimiter 创建登录尝试限流器。
func NewLoginLimiter(limit int, window time.Duration) *LoginLimiter {
	return &LoginLimiter{limit: limit, window: window, entries: make(map[string]loginWindow)}
}

// Allow 登记一次尝试，并在窗口额度耗尽时拒绝。
func (l *LoginLimiter) Allow(key string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.cleanup(now)
	entry := l.entries[key]
	if entry.startedAt.IsZero() || now.Sub(entry.startedAt) >= l.window {
		entry = loginWindow{startedAt: now}
	}
	if entry.attempts >= l.limit {
		return false
	}
	entry.attempts++
	l.entries[key] = entry
	return true
}

func (l *LoginLimiter) cleanup(now time.Time) {
	if !l.lastGC.IsZero() && now.Sub(l.lastGC) < l.window {
		return
	}
	for key, entry := range l.entries {
		if now.Sub(entry.startedAt) >= l.window {
			delete(l.entries, key)
		}
	}
	l.lastGC = now
}

func loginLimitKey(ip, phone string) string {
	return strings.TrimSpace(ip) + "\x00" + strings.TrimSpace(phone)
}
