/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package integration_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
	"base-engine/auth"
	"base-engine/gen"
)

// TestWebSocketRevalidatesRevokedSession 验证同一升级连接不会继续复用已经撤销的握手身份。
func TestWebSocketRevalidatesRevokedSession(t *testing.T) {
	fixture := newSecurityFixture(t)
	server := httptest.NewServer(fixture.handler)
	defer server.Close()
	fixture.cfg.AllowedOrigins[server.URL] = struct{}{}

	header := http.Header{"Origin": []string{server.URL}}
	header.Set("Cookie", fixture.cookies["session-a-1"].String())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	connection, _, err := websocket.Dial(ctx, "ws"+strings.TrimPrefix(server.URL, "http")+"/graphql", &websocket.DialOptions{
		HTTPHeader: header, Subprotocols: []string{"graphql-transport-ws"},
	})
	if err != nil {
		if strings.Contains(err.Error(), "operation not permitted") {
			t.Skipf("skipping WebSocket dial in sandbox environment without loopback network permission: %v", err)
		}
		t.Fatal(err)
	}
	defer connection.CloseNow()

	writeWebSocketMessage(t, ctx, connection, map[string]any{"type": "connection_init"})
	assertWebSocketMessageType(t, ctx, connection, "connection_ack")
	writeWebSocketOperation(t, ctx, connection, "before-revoke")
	assertWebSocketOperationCode(t, ctx, connection, "before-revoke", "")

	now := time.Now()
	if err := fixture.db.Model(&gen.Session{}).Where("id = ?", "session-a-1").Update("revoked_at", now).Error; err != nil {
		t.Fatal(err)
	}
	writeWebSocketOperation(t, ctx, connection, "after-revoke")
	assertWebSocketOperationCode(t, ctx, connection, "after-revoke", auth.CodeSessionRevoked)
}

func writeWebSocketOperation(t *testing.T, ctx context.Context, connection *websocket.Conn, id string) {
	t.Helper()
	writeWebSocketMessage(t, ctx, connection, map[string]any{
		"id": id, "type": "subscribe", "payload": map[string]string{"query": `query { viewer { account { id } } }`},
	})
}

func writeWebSocketMessage(t *testing.T, ctx context.Context, connection *websocket.Conn, value any) {
	t.Helper()
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	if err := connection.Write(ctx, websocket.MessageText, encoded); err != nil {
		t.Fatal(err)
	}
}

func assertWebSocketMessageType(t *testing.T, ctx context.Context, connection *websocket.Conn, expected string) {
	t.Helper()
	_, encoded, err := connection.Read(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var message struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(encoded, &message); err != nil || message.Type != expected {
		t.Fatalf("message type=%q want=%q err=%v body=%s", message.Type, expected, err, encoded)
	}
}

func assertWebSocketOperationCode(t *testing.T, ctx context.Context, connection *websocket.Conn, id string, expected auth.Code) {
	t.Helper()
	_, encoded, err := connection.Read(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var message struct {
		ID      string          `json:"id"`
		Type    string          `json:"type"`
		Payload graphQLResponse `json:"payload"`
	}
	if err := json.Unmarshal(encoded, &message); err != nil {
		t.Fatal(err)
	}
	if message.ID != id || message.Type != "next" {
		t.Fatalf("unexpected operation message: %s", encoded)
	}
	if expected == "" {
		assertNoErrors(t, message.Payload)
	} else {
		assertCode(t, message.Payload, expected)
	}
	assertWebSocketMessageType(t, ctx, connection, "complete")
}
