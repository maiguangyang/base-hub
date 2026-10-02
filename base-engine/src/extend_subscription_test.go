/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package src_test

import (
	"context"
	"testing"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	enginesrc "base-engine/src"
	sessionservice "base-engine/src/services/session"
)

// TestSessionEventsResolverUsesAuthenticatedSession 验证订阅键只能来自当前权威 Principal。
func TestSessionEventsResolverUsesAuthenticatedSession(t *testing.T) {
	database := openResolverTestDB(t)
	publisher := sessionservice.NewPublisher()
	dependencies := enginesrc.NewDependencies(database, resolverTestConfig(), publisher)
	resolver := enginesrc.NewResolver(gen.NewDB(database), &gen.EventController{}, dependencies).Subscription()
	if _, err := resolver.SessionEvents(context.Background()); auth.ErrorCode(err) != auth.CodeAuthRequired {
		t.Fatalf("missing principal error = %v", err)
	}
	ctx, cancel := context.WithCancel(auth.WithPrincipal(context.Background(), &auth.WorkspacePrincipal{SessionID: "session-current"}))
	channel, err := resolver.SessionEvents(ctx)
	if err != nil {
		t.Fatal(err)
	}
	event := &gen.SessionEvent{Code: gen.SessionEventCodeOrganizationSuspended, SessionID: "session-current", OccurredAt: time.Now()}
	publisher.PublishSession("session-other", event)
	assertNoSessionEvent(t, channel)
	publisher.PublishSession("session-current", event)
	assertSessionEvent(t, channel, event)
	cancel()
	assertSubscriptionClosed(t, channel)
}

func assertNoSessionEvent(t *testing.T, channel <-chan *gen.SessionEvent) {
	t.Helper()
	select {
	case received := <-channel:
		t.Fatalf("received another session event: %#v", received)
	case <-time.After(20 * time.Millisecond):
	}
}

func assertSessionEvent(t *testing.T, channel <-chan *gen.SessionEvent, expected *gen.SessionEvent) {
	t.Helper()
	select {
	case received := <-channel:
		if received != expected {
			t.Fatalf("event = %#v", received)
		}
	case <-time.After(time.Second):
		t.Fatal("session event timeout")
	}
}

func assertSubscriptionClosed(t *testing.T, channel <-chan *gen.SessionEvent) {
	t.Helper()
	select {
	case _, open := <-channel:
		if open {
			t.Fatal("cancelled subscription remained open")
		}
	case <-time.After(time.Second):
		t.Fatal("cancelled subscription was not closed")
	}
}

// TestLegacyWebSocketSubscriptionIsDisabled 验证生成器遗留订阅入口不能成为业务通道。
func TestLegacyWebSocketSubscriptionIsDisabled(t *testing.T) {
	database := openResolverTestDB(t)
	dependencies := enginesrc.NewDependencies(database, resolverTestConfig(), sessionservice.NewPublisher())
	configuration := enginesrc.New(gen.NewDB(database), &gen.EventController{}, dependencies)
	resolver := configuration.Resolvers.Subscription()
	if _, err := resolver.WebSocket(context.Background()); auth.ErrorCode(err) != auth.CodeSessionPubSubRequired {
		t.Fatalf("legacy WebSocket error = %v", err)
	}
}
