package agentkit

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

type storedApproval struct {
	plan       json.RawMessage
	binding    string
	expires    time.Time
	digest     string
	prompt     string
	guard      *ApprovalGuard
	timer      *time.Timer
	issued     bool
	used       bool
	validating bool
}

// ApprovalStore keeps one-time plans in this process. Callers provide authoritative identity bindings.
type ApprovalStore[T any] struct {
	mu    sync.Mutex
	key   []byte
	plans map[string]*storedApproval
	now   func() time.Time
}

func NewApprovalStore[T any](now func() time.Time) (*ApprovalStore[T], error) {
	if now == nil {
		return nil, errors.New("MISSING_CLOCK")
	}
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return &ApprovalStore[T]{key: key, plans: map[string]*storedApproval{}, now: now}, nil
}

func (store *ApprovalStore[T]) Register(id, binding string, expires time.Time, steps []ApprovedStep, plan T, prompt string) error {
	if store == nil || id == "" || binding == "" || !store.now().Before(expires) {
		return errors.New("INVALID_PLAN")
	}
	guard, err := NewApprovalGuard(steps)
	if err != nil {
		return err
	}
	encoded, err := json.Marshal(plan)
	if err != nil {
		return err
	}
	digestInput, err := json.Marshal(struct {
		Plan    json.RawMessage
		Binding string
		Expires time.Time
		Steps   []ApprovedStep
	}{Plan: encoded, Binding: binding, Expires: expires, Steps: steps})
	if err != nil {
		return err
	}
	digest := sha256.Sum256(digestInput)
	store.mu.Lock()
	defer store.mu.Unlock()
	if store.plans[id] != nil {
		return errors.New("PLAN_ID_COLLISION")
	}
	state := &storedApproval{plan: encoded, binding: binding, expires: expires, digest: base64.RawURLEncoding.EncodeToString(digest[:]), prompt: prompt, guard: guard}
	store.plans[id] = state
	state.timer = time.AfterFunc(expires.Sub(store.now()), func() { store.expireUnused(id) })
	return nil
}

func (store *ApprovalStore[T]) Issue(id string, plan T) (string, error) {
	encoded, err := json.Marshal(plan)
	if err != nil {
		return "", err
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	state := store.plans[id]
	if state == nil || state.issued || !store.now().Before(state.expires) || !bytes.Equal(encoded, state.plan) {
		return "", errors.New("PLAN_NOT_VALIDATED")
	}
	payload, err := json.Marshal(approvalToken{ID: id, Expires: state.expires.Unix(), Digest: state.digest})
	if err != nil {
		return "", err
	}
	encodedToken := base64.RawURLEncoding.EncodeToString(payload)
	state.issued = true
	return encodedToken + "." + base64.RawURLEncoding.EncodeToString(store.sign(encodedToken)), nil
}

func (store *ApprovalStore[T]) Consume(token, binding string, validate func(T) bool) (T, error) {
	var zero T
	payload, err := store.parseToken(token)
	if err != nil {
		return zero, err
	}
	store.mu.Lock()
	state := store.plans[payload.ID]
	if !validConsumeState(state, payload, store.now()) {
		store.mu.Unlock()
		return zero, errors.New("TOKEN_EXPIRED_OR_USED")
	}
	if binding != state.binding {
		store.mu.Unlock()
		return zero, errors.New("TOKEN_BINDING_MISMATCH")
	}
	var plan T
	if err := json.Unmarshal(state.plan, &plan); err != nil {
		store.mu.Unlock()
		return zero, err
	}
	if validate == nil {
		state.used = true
		store.mu.Unlock()
		return plan, nil
	}
	state.validating = true
	store.mu.Unlock()
	approved := validate(plan)
	store.mu.Lock()
	defer store.mu.Unlock()
	if store.plans[payload.ID] != state {
		return zero, errors.New("TOKEN_EXPIRED_OR_USED")
	}
	state.validating = false
	if !validConsumeState(state, payload, store.now()) {
		return zero, errors.New("TOKEN_EXPIRED_OR_USED")
	}
	if !approved {
		return zero, errors.New("TOKEN_BINDING_MISMATCH")
	}
	state.used = true
	return plan, nil
}

func validConsumeState(state *storedApproval, payload approvalToken, now time.Time) bool {
	return state != nil && state.issued && !state.used && !state.validating && now.Before(state.expires) &&
		state.expires.Unix() == payload.Expires && state.digest == payload.Digest
}

func (store *ApprovalStore[T]) Prompt(id string) string {
	store.mu.Lock()
	defer store.mu.Unlock()
	if state := store.plans[id]; state != nil && state.used {
		return state.prompt
	}
	return ""
}

func (store *ApprovalStore[T]) UsedPlan(id string) (T, bool) {
	var plan T
	store.mu.Lock()
	defer store.mu.Unlock()
	state := store.plans[id]
	if state == nil || !state.used || json.Unmarshal(state.plan, &plan) != nil {
		return plan, false
	}
	return plan, true
}

func (store *ApprovalStore[T]) AuthorizeStep(id, toolID string, args json.RawMessage) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	state := store.plans[id]
	if state == nil || !state.used {
		return errors.New("PLAN_NOT_ACTIVE")
	}
	return state.guard.Reserve(toolID, args)
}

func (store *ApprovalStore[T]) RecordSuccess(id, toolID string) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	state := store.plans[id]
	if state == nil || !state.used {
		return errors.New("STEP_OUTCOME_MISMATCH")
	}
	return state.guard.RecordSuccess(toolID)
}

func (store *ApprovalStore[T]) Complete(id string) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	state := store.plans[id]
	if state == nil || !state.used {
		return errors.New("PLAN_NOT_ACTIVE")
	}
	return state.guard.Complete()
}

func (store *ApprovalStore[T]) Cancel(id string) {
	store.mu.Lock()
	defer store.mu.Unlock()
	if state := store.plans[id]; state != nil {
		state.prompt = ""
		if state.timer != nil {
			state.timer.Stop()
		}
		delete(store.plans, id)
	}
}

func (store *ApprovalStore[T]) Has(id string) bool {
	store.mu.Lock()
	defer store.mu.Unlock()
	return store.plans[id] != nil
}

func (store *ApprovalStore[T]) SweepExpired() {
	store.mu.Lock()
	defer store.mu.Unlock()
	for id, state := range store.plans {
		if state.used || store.now().Before(state.expires) {
			continue
		}
		state.prompt = ""
		if state.timer != nil {
			state.timer.Stop()
		}
		delete(store.plans, id)
	}
}

func (store *ApprovalStore[T]) expireUnused(id string) {
	store.mu.Lock()
	defer store.mu.Unlock()
	state := store.plans[id]
	if state == nil || state.used {
		return
	}
	if remaining := state.expires.Sub(store.now()); remaining > 0 {
		state.timer = time.AfterFunc(remaining, func() { store.expireUnused(id) })
		return
	}
	state.prompt = ""
	delete(store.plans, id)
}
