package agentkit

import (
	"bytes"
	"encoding/json"
	"errors"
	"sync"
)

type ApprovedStep struct {
	ToolID    string
	Arguments json.RawMessage
	MaxCalls  int
}

// ApprovalGuard enforces the exact arguments, sequence, and call limits of a plan.
type ApprovalGuard struct {
	mu        sync.Mutex
	steps     []ApprovedStep
	counts    []int
	successes []int
	current   int
	pending   bool
}

func NewApprovalGuard(steps []ApprovedStep) (*ApprovalGuard, error) {
	if len(steps) == 0 || len(steps) > 100 {
		return nil, errors.New("INVALID_APPROVED_STEPS")
	}
	guard := &ApprovalGuard{counts: make([]int, len(steps)), successes: make([]int, len(steps))}
	for _, step := range steps {
		if step.ToolID == "" || step.MaxCalls < 1 {
			return nil, errors.New("INVALID_APPROVED_STEP")
		}
		canonical, err := canonicalArguments(step.Arguments)
		if err != nil {
			return nil, err
		}
		step.Arguments = canonical
		guard.steps = append(guard.steps, step)
	}
	return guard, nil
}

func canonicalArguments(args json.RawMessage) ([]byte, error) {
	var value map[string]any
	if !json.Valid(args) {
		return nil, errors.New("INVALID_STEP_ARGUMENTS")
	}
	decoder := json.NewDecoder(bytes.NewReader(args))
	decoder.UseNumber()
	if err := decoder.Decode(&value); err != nil || value == nil {
		return nil, errors.New("INVALID_STEP_ARGUMENTS")
	}
	return json.Marshal(value)
}

func (guard *ApprovalGuard) Reserve(toolID string, args json.RawMessage) error {
	canonical, err := canonicalArguments(args)
	if err != nil {
		return err
	}
	guard.mu.Lock()
	defer guard.mu.Unlock()
	if guard.pending {
		return errors.New("STEP_PENDING")
	}
	index := guard.current
	if index >= len(guard.steps) {
		return errors.New("PLAN_EXHAUSTED")
	}
	if !stepMatches(guard.steps[index], toolID, canonical) {
		if guard.counts[index] == 0 || index+1 >= len(guard.steps) || !stepMatches(guard.steps[index+1], toolID, canonical) {
			return errors.New("STEP_ARGUMENTS_OR_ORDER_MISMATCH")
		}
		index++
	}
	if guard.counts[index] >= guard.steps[index].MaxCalls {
		return errors.New("STEP_CALL_LIMIT")
	}
	guard.current = index
	guard.counts[index]++
	guard.pending = true
	return nil
}

func stepMatches(step ApprovedStep, toolID string, canonical []byte) bool {
	return step.ToolID == toolID && bytes.Equal(step.Arguments, canonical)
}

func (guard *ApprovalGuard) RecordSuccess(toolID string) error {
	guard.mu.Lock()
	defer guard.mu.Unlock()
	if !guard.pending || guard.steps[guard.current].ToolID != toolID {
		return errors.New("STEP_OUTCOME_MISMATCH")
	}
	guard.successes[guard.current]++
	guard.pending = false
	return nil
}

func (guard *ApprovalGuard) Complete() error {
	guard.mu.Lock()
	defer guard.mu.Unlock()
	for _, successes := range guard.successes {
		if successes == 0 {
			return errors.New("APPROVED_STEP_NOT_EXECUTED")
		}
	}
	if guard.pending {
		return errors.New("APPROVED_STEP_PENDING")
	}
	return nil
}
