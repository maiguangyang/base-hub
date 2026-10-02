/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-22
 */
package session

import (
	"context"
	"testing"
	"time"

	"base-engine/gen"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRevokeWorkspaceSessions(t *testing.T) {
	database, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{IgnoreRelationshipsWhenMigrating: true})
	if err != nil {
		t.Fatal(err)
	}
	if err := database.AutoMigrate(&gen.Session{}); err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	hqID, otherID := "hq-org", "other-org"
	sessions := []gen.Session{
		{ID: "target", AccountID: "account", OrganizationID: &hqID, WorkspaceType: gen.WorkspaceTypeHeadquarters, ExpiresAt: now.Add(time.Hour)},
		{ID: "other-workspace", AccountID: "account", OrganizationID: &otherID, WorkspaceType: gen.WorkspaceTypeHeadquarters, ExpiresAt: now.Add(time.Hour)},
		{ID: "other-account", AccountID: "other", OrganizationID: &hqID, WorkspaceType: gen.WorkspaceTypeHeadquarters, ExpiresAt: now.Add(time.Hour)},
		{ID: "discovery", AccountID: "account", WorkspaceType: gen.WorkspaceTypeDiscovery, ExpiresAt: now.Add(time.Hour)},
	}
	for index := range sessions {
		if err := database.Create(&sessions[index]).Error; err != nil {
			t.Fatal(err)
		}
	}
	ids, err := NewService(time.Hour).RevokeWorkspaceSessions(context.Background(), database, "account", hqID, "AUTHORITY_CHANGED", now)
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 1 || ids[0] != "target" {
		t.Fatalf("revoked IDs = %#v", ids)
	}
	assertSessionsRemainActive(t, database, "other-workspace", "other-account", "discovery")
}

func assertSessionsRemainActive(t *testing.T, database *gorm.DB, sessionIDs ...string) {
	t.Helper()
	for _, id := range sessionIDs {
		var stored gen.Session
		if err := database.First(&stored, "id = ?", id).Error; err != nil || stored.RevokedAt != nil {
			t.Fatalf("unrelated session %s changed: %#v, %v", id, stored, err)
		}
	}
}

func TestPendingSessionEventsPublishOnlyWhenFinished(t *testing.T) {
	publisher := &recordingPublisher{}
	organizationID := "hq-org"
	now := time.Now()
	ctx := WithPendingEvents(context.Background())
	QueueRevokedSessions(ctx, []string{"session-1", "session-2"}, organizationID, now)
	if len(publisher.events) != 0 {
		t.Fatalf("events published before finish: %#v", publisher.events)
	}
	PublishPending(ctx, publisher)
	if len(publisher.events) != 2 {
		t.Fatalf("published events = %#v", publisher.events)
	}
	for _, event := range publisher.events {
		if event.event.Code != gen.SessionEventCodeSessionRevoked || event.event.OrganizationID == nil || *event.event.OrganizationID != organizationID {
			t.Fatalf("event = %#v", event.event)
		}
	}
}

type publishedSessionEvent struct {
	sessionID string
	event     *gen.SessionEvent
}

type recordingPublisher struct{ events []publishedSessionEvent }

func (p *recordingPublisher) Subscribe(context.Context, string) (<-chan *gen.SessionEvent, error) {
	return nil, nil
}

func (p *recordingPublisher) PublishSession(sessionID string, event *gen.SessionEvent) {
	p.events = append(p.events, publishedSessionEvent{sessionID: sessionID, event: event})
}
