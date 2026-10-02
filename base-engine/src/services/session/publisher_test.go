/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package session

import (
	"context"
	"testing"
	"time"

	"base-engine/gen"
)

// TestPublisherBroadcastsToEveryTabAndCoalesces 验证同一 Session 多标签广播且慢订阅者不阻塞。
func TestPublisherBroadcastsToEveryTabAndCoalesces(t *testing.T) {
	publisher := NewPublisher()
	first, cancelFirst := subscribeForTest(t, publisher, "session-1")
	defer cancelFirst()
	second, cancelSecond := subscribeForTest(t, publisher, "session-1")
	defer cancelSecond()
	event := &gen.SessionEvent{Code: gen.SessionEventCodeSessionRevoked, SessionID: "session-1", OccurredAt: time.Now()}
	publisher.PublishSession("session-1", event)
	publisher.PublishSession("session-1", event)
	assertEvent(t, first, gen.SessionEventCodeSessionRevoked)
	assertEvent(t, second, gen.SessionEventCodeSessionRevoked)
}

// TestPublisherCancellationClosesOnlySubscriber 验证取消只注销当前标签页。
func TestPublisherCancellationClosesOnlySubscriber(t *testing.T) {
	publisher := NewPublisher()
	first, cancelFirst := subscribeForTest(t, publisher, "session-1")
	second, cancelSecond := subscribeForTest(t, publisher, "session-1")
	defer cancelSecond()
	cancelFirst()
	select {
	case _, open := <-first:
		if open {
			t.Fatal("cancelled subscriber remained open")
		}
	case <-time.After(time.Second):
		t.Fatal("cancelled subscriber was not closed")
	}
	event := &gen.SessionEvent{Code: gen.SessionEventCodeSessionRevoked, SessionID: "session-1", OccurredAt: time.Now()}
	publisher.PublishSession("session-1", event)
	assertEvent(t, second, gen.SessionEventCodeSessionRevoked)
}

func subscribeForTest(t *testing.T, publisher *LocalPublisher, sessionID string) (<-chan *gen.SessionEvent, context.CancelFunc) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	channel, err := publisher.Subscribe(ctx, sessionID)
	if err != nil {
		cancel()
		t.Fatal(err)
	}
	return channel, cancel
}

func assertEvent(t *testing.T, channel <-chan *gen.SessionEvent, code gen.SessionEventCode) {
	t.Helper()
	select {
	case event := <-channel:
		if event == nil || event.Code != code {
			t.Fatalf("event = %#v", event)
		}
	case <-time.After(time.Second):
		t.Fatal("event timeout")
	}
}
