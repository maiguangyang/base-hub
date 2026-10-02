package agentkit

import (
	"encoding/json"
	"testing"
	"time"
)

type storeTestPlan struct{ Request string }
type storeFixture struct {
	store   *ApprovalStore[storeTestPlan]
	now     time.Time
	expires time.Time
	plan    storeTestPlan
	token   string
}

func newStoreFixture(t *testing.T) *storeFixture {
	t.Helper()
	fixture := &storeFixture{now: time.Unix(1_800_000_000, 0), plan: storeTestPlan{Request: "create"}}
	store, err := NewApprovalStore[storeTestPlan](func() time.Time { return fixture.now })
	if err != nil {
		t.Fatal(err)
	}
	fixture.store = store
	fixture.expires = fixture.now.Add(5 * time.Minute)
	steps := []ApprovedStep{{ToolID: "write", Arguments: json.RawMessage(`{"id":"a"}`), MaxCalls: 1}}
	if err := store.Register("plan-1", "user/session", fixture.expires, steps, fixture.plan, "private prompt"); err != nil {
		t.Fatal(err)
	}
	fixture.token, err = store.Issue("plan-1", fixture.plan)
	if err != nil {
		t.Fatal(err)
	}
	return fixture
}

func TestApprovalStoreBindsAndConsumesTokenOnce(t *testing.T) {
	fixture := newStoreFixture(t)
	if _, err := fixture.store.Consume(fixture.token, "other/session", nil); err == nil {
		t.Fatal("token accepted for another principal")
	}
	got, err := fixture.store.Consume(fixture.token, "user/session", nil)
	if err != nil || got != fixture.plan || fixture.store.Prompt("plan-1") != "private prompt" {
		t.Fatalf("consume plan=%+v prompt=%q err=%v", got, fixture.store.Prompt("plan-1"), err)
	}
	if _, err := fixture.store.Consume(fixture.token, "user/session", nil); err == nil {
		t.Fatal("token replay accepted")
	}
}

func TestApprovalStoreKeepsConsumedRunActive(t *testing.T) {
	fixture := newStoreFixture(t)
	if _, err := fixture.store.Consume(fixture.token, "user/session", nil); err != nil {
		t.Fatal(err)
	}
	fixture.now = fixture.expires.Add(time.Second)
	if err := fixture.store.AuthorizeStep("plan-1", "write", json.RawMessage(`{"id":"a"}`)); err != nil {
		t.Fatalf("consumed plan expired during run: %v", err)
	}
	if err := fixture.store.RecordSuccess("plan-1", "write"); err != nil {
		t.Fatal(err)
	}
	if err := fixture.store.Complete("plan-1"); err != nil {
		t.Fatal(err)
	}
	fixture.store.Cancel("plan-1")
	if fixture.store.Has("plan-1") {
		t.Fatal("cancelled plan remains stored")
	}
}

func TestApprovalStoreConsumeValidatorCanReadStore(t *testing.T) {
	fixture := newStoreFixture(t)
	plan, err := fixture.store.Consume(fixture.token, "user/session", func(plan storeTestPlan) bool {
		return plan == fixture.plan && fixture.store.Has("plan-1")
	})
	if err != nil || plan != fixture.plan {
		t.Fatalf("consume plan=%+v err=%v", plan, err)
	}
}

func TestApprovalStoreConsumeReservesTokenDuringValidation(t *testing.T) {
	fixture := newStoreFixture(t)
	entered := make(chan struct{})
	release := make(chan struct{})
	result := make(chan error, 1)
	go func() {
		_, err := fixture.store.Consume(fixture.token, "user/session", func(storeTestPlan) bool {
			close(entered)
			<-release
			return true
		})
		result <- err
	}()
	<-entered
	if _, err := fixture.store.Consume(fixture.token, "user/session", nil); err == nil {
		close(release)
		t.Fatal("concurrent consume accepted a token under validation")
	}
	close(release)
	if err := <-result; err != nil {
		t.Fatalf("first consume failed: %v", err)
	}
}

func TestApprovalStoreConsumeCanRetryAfterRejectedValidation(t *testing.T) {
	fixture := newStoreFixture(t)
	if _, err := fixture.store.Consume(fixture.token, "user/session", func(storeTestPlan) bool { return false }); err == nil {
		t.Fatal("rejected validation consumed the token")
	}
	if _, err := fixture.store.Consume(fixture.token, "user/session", nil); err != nil {
		t.Fatalf("valid retry was rejected: %v", err)
	}
}

func TestApprovalStoreConsumeRejectsPlanCancelledDuringValidation(t *testing.T) {
	fixture := newStoreFixture(t)
	entered := make(chan struct{})
	release := make(chan struct{})
	result := make(chan error, 1)
	go func() {
		_, err := fixture.store.Consume(fixture.token, "user/session", func(storeTestPlan) bool {
			close(entered)
			<-release
			return true
		})
		result <- err
	}()
	<-entered
	fixture.store.Cancel("plan-1")
	close(release)
	if err := <-result; err == nil {
		t.Fatal("cancelled plan was consumed after validation")
	}
}
