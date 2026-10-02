package ai

import (
	"bytes"
	"encoding/json"
	"errors"
)

func (s *ApprovalService) bindAttestationArguments(planID string, spec ToolSpec, modelArgs json.RawMessage) (json.RawMessage, error) {
	if !spec.RequiredAttestation {
		return modelArgs, nil
	}
	plan, active := s.store.UsedPlan(planID)
	if !active {
		return nil, errors.New("PLAN_NOT_ACTIVE")
	}
	for _, step := range plan.Steps {
		if step.ToolID != spec.ID {
			continue
		}
		matches, err := matchesAttestedStep(step.Arguments, modelArgs)
		if err != nil {
			return nil, err
		}
		if matches {
			return step.Arguments, nil
		}
	}
	return nil, errors.New("ATTESTATION_ARGUMENT_MISMATCH")
}

func matchesAttestedStep(expectedRaw, modelArgs json.RawMessage) (bool, error) {
	var expected map[string]any
	if err := json.Unmarshal(expectedRaw, &expected); err != nil {
		return false, err
	}
	delete(expected, "evidenceReference")
	delete(expected, "attested")
	stripped, err := json.Marshal(expected)
	if err != nil {
		return false, err
	}
	var actual any
	if err := json.Unmarshal(modelArgs, &actual); err != nil {
		return false, err
	}
	canonical, err := json.Marshal(actual)
	return bytes.Equal(stripped, canonical), err
}

func (s *ApprovalService) AuthorizeStep(planID, toolID string, args json.RawMessage) error {
	if s == nil || planID == "" || toolID == "" {
		return errors.New("INVALID_STEP_REQUEST")
	}
	return s.store.AuthorizeStep(planID, toolID, args)
}

func (s *ApprovalService) Complete(planID string) error {
	return s.store.Complete(planID)
}

func (s *ApprovalService) RecordSuccess(planID, toolID string) error {
	return s.store.RecordSuccess(planID, toolID)
}

func (s *ApprovalService) Cancel(planID string) {
	s.store.Cancel(planID)
}
