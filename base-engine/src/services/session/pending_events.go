/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-22
 */
package session

import (
	"context"
	"time"

	"base-engine/gen"
)

type pendingEventsKey struct{}

type pendingSessionEvent struct {
	sessionID string
	event     *gen.SessionEvent
}

type pendingSessionEvents struct {
	events []pendingSessionEvent
}

// WithPendingEvents attaches a request-owned queue for post-commit publication.
func WithPendingEvents(ctx context.Context) context.Context {
	return context.WithValue(ctx, pendingEventsKey{}, &pendingSessionEvents{})
}

// QueueRevokedSessions appends terminal workspace events without publishing them.
func QueueRevokedSessions(ctx context.Context, sessionIDs []string, organizationID string, occurredAt time.Time) {
	queue, _ := ctx.Value(pendingEventsKey{}).(*pendingSessionEvents)
	if queue == nil {
		return
	}
	for _, sessionID := range sessionIDs {
		organization := organizationID
		queue.events = append(queue.events, pendingSessionEvent{sessionID: sessionID, event: &gen.SessionEvent{
			Code: gen.SessionEventCodeSessionRevoked, SessionID: sessionID,
			OrganizationID: &organization, OccurredAt: occurredAt,
		}})
	}
}

// PublishPending publishes queued events after the owning transaction commits.
func PublishPending(ctx context.Context, publisher Publisher) {
	queue, _ := ctx.Value(pendingEventsKey{}).(*pendingSessionEvents)
	if queue == nil || publisher == nil {
		return
	}
	for _, pending := range queue.events {
		publisher.PublishSession(pending.sessionID, pending.event)
	}
	queue.events = nil
}
