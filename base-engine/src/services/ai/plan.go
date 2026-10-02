package ai

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"base-engine/agentkit"
	"base-engine/auth"
)

type ApprovedStep struct {
	OperationID, ToolID string
	TargetIDs, ScopeIDs []string
	Arguments           json.RawMessage
	MaxCalls, Sequence  int
}

type PlanDraft struct {
	Steps           []ApprovedStep
	promptHash      string
	prompt          string
	operatorMessage string
	attachmentID    string
}

func (draft *PlanDraft) BindPromptHash(hash string) { draft.promptHash = hash }

func (draft *PlanDraft) BindPrompt(prompt string) {
	hash := sha256.Sum256([]byte(prompt))
	draft.promptHash = hex.EncodeToString(hash[:])
	draft.prompt = prompt
}

func (draft *PlanDraft) BindOperatorMessage(message string) { draft.operatorMessage = message }
func (draft *PlanDraft) BindAttachmentID(id string)         { draft.attachmentID = id }

type UserAttestation struct {
	EvidenceReference string `json:"evidenceReference"`
	Attested          bool   `json:"attested"`
}

type Plan struct {
	ID, AccountID, SessionID, WorkspaceType, OrganizationID, PromptHash string
	AttachmentID                                                        string
	Steps                                                               []ApprovedStep
	ExpiresAt                                                           time.Time
}

type ApprovalService struct {
	catalog *Catalog
	store   *agentkit.ApprovalStore[Plan]
	now     func() time.Time
}

func NewApprovalService(catalog *Catalog) (*ApprovalService, error) {
	if catalog == nil {
		return nil, errors.New("MISSING_CATALOG")
	}
	service := &ApprovalService{catalog: catalog, now: time.Now}
	store, err := agentkit.NewApprovalStore[Plan](func() time.Time { return service.now() })
	if err != nil {
		return nil, err
	}
	service.store = store
	return service, nil
}

func (s *ApprovalService) ValidateDraft(draft PlanDraft, principal *auth.WorkspacePrincipal, attestation *UserAttestation) (Plan, []string, error) {
	if !validDraftContext(s, draft, principal) {
		return Plan{}, nil, errors.New("INVALID_PLAN_DRAFT")
	}
	plan := newPlan(draft, principal, s.now())
	if plan.ID == "" {
		return Plan{}, nil, errors.New("PLAN_ID_UNAVAILABLE")
	}
	totalCalls := 0
	for index, proposed := range draft.Steps {
		step, required, err := s.validateDraftStep(proposed, len(draft.Steps), index+1, principal, attestation, draft.operatorMessage)
		if err != nil {
			return Plan{}, nil, err
		}
		if len(required) > 0 {
			return Plan{}, required, nil
		}
		if err := validateImageAttachmentStep(step, draft.attachmentID); err != nil {
			return Plan{}, nil, err
		}
		if duplicatePlanStep(plan.Steps, step) {
			return Plan{}, nil, errors.New("DUPLICATE_PLAN_STEP")
		}
		plan.Steps = append(plan.Steps, step)
		totalCalls += step.MaxCalls
		if totalCalls > 100 {
			return Plan{}, nil, errors.New("PLAN_CALL_LIMIT")
		}
	}
	if err := s.registerValidatedPlan(plan, draft.prompt); err != nil {
		return Plan{}, nil, err
	}
	return plan, nil, nil
}

func (s *ApprovalService) validateDraftStep(proposed ApprovedStep, stepCount, sequence int, principal *auth.WorkspacePrincipal, attestation *UserAttestation, operatorMessage string) (ApprovedStep, []string, error) {
	if s.requiresSingleCall(proposed) && (stepCount != 1 || proposed.MaxCalls != 1) {
		return ApprovedStep{}, nil, errors.New("SINGLE_CALL_APPROVAL_REQUIRED")
	}
	return s.validateStep(proposed, sequence, principal, attestation, operatorMessage)
}

func (s *ApprovalService) requiresSingleCall(step ApprovedStep) bool {
	spec, ok := s.catalog.Lookup(step.ToolID)
	return ok && spec.SingleCallApproval
}

func duplicatePlanStep(existing []ApprovedStep, candidate ApprovedStep) bool {
	for _, step := range existing {
		if step.ToolID == candidate.ToolID && bytes.Equal(step.Arguments, candidate.Arguments) {
			return true
		}
	}
	return false
}

func (s *ApprovalService) registerValidatedPlan(plan Plan, prompt string) error {
	steps := make([]agentkit.ApprovedStep, 0, len(plan.Steps))
	for _, step := range plan.Steps {
		steps = append(steps, agentkit.ApprovedStep{ToolID: step.ToolID, Arguments: step.Arguments, MaxCalls: step.MaxCalls})
	}
	return s.store.Register(plan.ID, planBinding(plan), plan.ExpiresAt, steps, plan, prompt)
}

func (s *ApprovalService) expireUnusedPlan(planID string) {
	s.store.SweepExpired()
}

func (s *ApprovalService) promptFor(planID string) string {
	return s.store.Prompt(planID)
}

func validDraftContext(service *ApprovalService, draft PlanDraft, principal *auth.WorkspacePrincipal) bool {
	return service != nil && principal != nil && principal.AccountID != "" && principal.SessionID != "" &&
		draft.promptHash != "" && len(draft.Steps) > 0 && len(draft.Steps) <= 100
}

func newPlan(draft PlanDraft, principal *auth.WorkspacePrincipal, now time.Time) Plan {
	plan := Plan{ID: randomPlanID(), AccountID: principal.AccountID, SessionID: principal.SessionID,
		WorkspaceType: string(principal.WorkspaceType), PromptHash: draft.promptHash, AttachmentID: draft.attachmentID, ExpiresAt: now.Add(5 * time.Minute)}
	if principal.OrganizationID != nil {
		plan.OrganizationID = *principal.OrganizationID
	}
	return plan
}

func randomPlanID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}
