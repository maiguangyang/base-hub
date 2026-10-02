/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-22
 */
package src

import (
	"context"
	"errors"
	"testing"
	"time"

	"base-engine/gen"
	sessionservice "base-engine/src/services/session"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGeneratedMutationPublishesAfterCommit(t *testing.T) {
	database, resolver, target, now := newMutationSessionFixture(t)
	publisher := &commitAssertPublisher{t: t, db: database}
	ctx := sessionservice.WithPendingEvents(context.Background())
	ctx = gen.EnrichContextWithMutations(ctx, resolver)
	ids, err := sessionservice.NewService(time.Hour).RevokeWorkspaceSessions(ctx, gen.GetTransaction(ctx), target.AccountID, *target.OrganizationID, sessionservice.RevocationCodeAuthorityChanged, now)
	if err != nil {
		t.Fatal(err)
	}
	sessionservice.QueueRevokedSessions(ctx, ids, *target.OrganizationID, now)
	if err := finishGeneratedMutation(ctx, resolver, publisher, nil); err != nil {
		t.Fatal(err)
	}
	if publisher.count != 1 {
		t.Fatalf("published count = %d", publisher.count)
	}
}

func TestGeneratedMutationRollsBackWithoutPublishing(t *testing.T) {
	database, resolver, target, now := newMutationSessionFixture(t)
	publisher := &commitAssertPublisher{t: t, db: database}
	ctx := sessionservice.WithPendingEvents(context.Background())
	ctx = gen.EnrichContextWithMutations(ctx, resolver)
	ids, err := sessionservice.NewService(time.Hour).RevokeWorkspaceSessions(ctx, gen.GetTransaction(ctx), target.AccountID, *target.OrganizationID, sessionservice.RevocationCodeAuthorityChanged, now)
	if err != nil {
		t.Fatal(err)
	}
	sessionservice.QueueRevokedSessions(ctx, ids, *target.OrganizationID, now)
	if err := finishGeneratedMutation(ctx, resolver, publisher, errors.New("write failed")); err == nil {
		t.Fatal("handler failure was lost")
	}
	var stored gen.Session
	if err := database.First(&stored, "id = ?", target.ID).Error; err != nil || stored.RevokedAt != nil || publisher.count != 0 {
		t.Fatalf("rollback state = %#v, published=%d, err=%v", stored, publisher.count, err)
	}
}

func newMutationSessionFixture(t *testing.T) (*gorm.DB, *gen.GeneratedResolver, *gen.Session, time.Time) {
	t.Helper()
	database, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{IgnoreRelationshipsWhenMigrating: true})
	if err != nil {
		t.Fatal(err)
	}
	if err := database.AutoMigrate(&gen.Session{}); err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	organizationID := "hq-org"
	target := &gen.Session{ID: "target-session", AccountID: "target-account", OrganizationID: &organizationID, WorkspaceType: gen.WorkspaceTypeHeadquarters, ExpiresAt: now.Add(time.Hour), LastSeenAt: now}
	if err := database.Create(target).Error; err != nil {
		t.Fatal(err)
	}
	return database, &gen.GeneratedResolver{DB: gen.NewDB(database), Handlers: gen.DefaultResolutionHandlers()}, target, now
}

type commitAssertPublisher struct {
	t     *testing.T
	db    *gorm.DB
	count int
}

func (p *commitAssertPublisher) Subscribe(context.Context, string) (<-chan *gen.SessionEvent, error) {
	return nil, nil
}

func (p *commitAssertPublisher) PublishSession(sessionID string, _ *gen.SessionEvent) {
	p.t.Helper()
	var stored gen.Session
	if err := p.db.First(&stored, "id = ?", sessionID).Error; err != nil || stored.RevokedAt == nil {
		p.t.Fatalf("event published before commit: %#v, %v", stored, err)
	}
	p.count++
}
