package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"regexp"
	"sort"
	"strings"

	"base-engine/auth"
)

func (s *ApprovalService) validateStep(proposed ApprovedStep, sequence int, principal *auth.WorkspacePrincipal, attestation *UserAttestation, operatorMessage string) (ApprovedStep, []string, error) {
	spec, exists := s.catalog.Lookup(proposed.ToolID)
	if !exists || !validProposedStep(proposed, spec, sequence) {
		return ApprovedStep{}, nil, errors.New("INVALID_APPROVED_STEP")
	}
	if !specAllowed(spec, principal, PhaseRun, map[string]struct{}{spec.ID: {}}) {
		return ApprovedStep{}, nil, errors.New("STEP_NOT_AUTHORIZED")
	}
	args, err := decodeStepArguments(proposed.Arguments, spec.ArgumentFields)
	if err != nil {
		return ApprovedStep{}, nil, err
	}
	if err := validatePlanArguments(proposed.Arguments, spec); err != nil {
		return ApprovedStep{}, nil, errors.New("INVALID_STEP_ARGUMENTS")
	}
	if err := validateEvidenceSource(args, spec, operatorMessage); err != nil {
		return ApprovedStep{}, nil, err
	}
	required, err := bindUserAttestation(args, spec, attestation)
	if err != nil || len(required) > 0 {
		return ApprovedStep{}, required, err
	}
	canonical, err := json.Marshal(args)
	if err != nil {
		return ApprovedStep{}, nil, err
	}
	proposed.Arguments = canonical
	return validateStepTargets(proposed, spec, principal)
}

var customerPointReferencePattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]{3,63}$`)
var customerEvidencePattern = regexp.MustCompile(`(?i)(?:凭证编号|证明编号|工单号|proof reference|\bproof\b|work order)\s*[:：为]?\s*([A-Za-z0-9][A-Za-z0-9_-]{0,127})`)
var customerBasisPattern = regexp.MustCompile(`(?i)(?:处理依据代码|依据代码|basis code)\s*[:：为]?\s*([A-Za-z0-9][A-Za-z0-9_-]{0,127})`)

func validateEvidenceSource(args map[string]any, spec ToolSpec, operatorMessage string) error {
	for _, path := range spec.EvidencePaths {
		value, err := lookupArgumentPath(args, path)
		proof, ok := value.(string)
		if err != nil || !ok || strings.TrimSpace(proof) == "" {
			return errors.New("EVIDENCE_NOT_FROM_OPERATOR")
		}
		if !evidenceProofMatchesOperator(spec.ID, path, proof, operatorMessage) {
			return errors.New("EVIDENCE_NOT_FROM_OPERATOR")
		}
	}
	return nil
}

func evidenceProofMatchesOperator(toolID, path, proof, operatorMessage string) bool {
	if !isCustomerEvidenceMutation(toolID) {
		return strings.Contains(operatorMessage, proof)
	}
	if isCustomerPointMutation(toolID) && path == "note" && !validCustomerPointReference(proof) {
		return false
	}
	return operatorHasCustomerEvidence(operatorMessage, proof, path)
}

func isCustomerPointMutation(toolID string) bool {
	return toolID == "HqGrantCustomerPoints" || toolID == "HqReverseCustomerPoints"
}

func isCustomerEvidenceMutation(toolID string) bool {
	switch toolID {
	case "HqRequestCustomerCancellation", "HqCompleteCustomerCancellation", "HqGrantCustomerPoints",
		"HqReverseCustomerPoints", "HqCorrectCustomerPoints":
		return true
	default:
		return false
	}
}

func validCustomerPointReference(value string) bool {
	return customerPointReferencePattern.MatchString(value) && strings.ContainsAny(value, "0123456789_-")
}

func operatorHasCustomerEvidence(message, proof, path string) bool {
	pattern := customerEvidencePattern
	if path == "basisCode" {
		pattern = customerBasisPattern
	}
	for _, match := range pattern.FindAllStringSubmatch(message, -1) {
		if match[1] == proof {
			return true
		}
	}
	return false
}

func validatePlanArguments(raw json.RawMessage, spec ToolSpec) error {
	if spec.InputSchema == nil {
		return errors.New("MISSING_TOOL_INPUT_SCHEMA")
	}
	var args map[string]any
	if err := json.Unmarshal(raw, &args); err != nil {
		return err
	}
	return spec.InputSchema.Validate(args)
}

func validProposedStep(step ApprovedStep, spec ToolSpec, sequence int) bool {
	return spec.Mode == ModeWrite && spec.OperationID == step.OperationID && step.Sequence == sequence &&
		step.MaxCalls >= 1 && step.MaxCalls <= 10
}

func bindUserAttestation(args map[string]any, spec ToolSpec, attestation *UserAttestation) ([]string, error) {
	if !spec.RequiredAttestation {
		return nil, nil
	}
	if _, exists := args["evidenceReference"]; exists {
		return nil, errors.New("MODEL_ATTESTATION_FORBIDDEN")
	}
	if _, exists := args["attested"]; exists {
		return nil, errors.New("MODEL_ATTESTATION_FORBIDDEN")
	}
	if attestation == nil || !attestation.Attested || strings.TrimSpace(attestation.EvidenceReference) == "" {
		return []string{"evidenceReference", "attested"}, nil
	}
	args["evidenceReference"] = attestation.EvidenceReference
	args["attested"] = true
	return nil, nil
}

func decodeStepArguments(raw json.RawMessage, allowed []string) (map[string]any, error) {
	if len(allowed) == 0 || !json.Valid(raw) {
		return nil, errors.New("INVALID_STEP_ARGUMENTS")
	}
	var object map[string]any
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	if err := decoder.Decode(&object); err != nil || object == nil {
		return nil, errors.New("INVALID_STEP_ARGUMENTS")
	}
	fields := map[string]struct{}{}
	for _, key := range allowed {
		fields[key] = struct{}{}
	}
	for key := range object {
		if _, ok := fields[key]; !ok {
			return nil, errors.New("UNAPPROVED_ARGUMENT_FIELD")
		}
	}
	return object, nil
}

func validateStepTargets(step ApprovedStep, spec ToolSpec, principal *auth.WorkspacePrincipal) (ApprovedStep, []string, error) {
	var arguments map[string]any
	if err := json.Unmarshal(step.Arguments, &arguments); err != nil {
		return ApprovedStep{}, nil, err
	}
	derived, err := extractArgumentIDs(arguments, targetFieldPaths(spec))
	if err != nil {
		return ApprovedStep{}, nil, err
	}
	if spec.WriteKind == WriteCreate && len(derived) == 0 && principal.OrganizationID != nil {
		derived = []string{*principal.OrganizationID}
	}
	if spec.WriteKind == WriteCreate && len(derived) == 0 && len(spec.ScopeFields) == 0 {
		derived = []string{"workspace:" + string(principal.WorkspaceType)}
	}
	if len(derived) == 0 {
		return ApprovedStep{}, nil, errors.New("MISSING_STEP_TARGET")
	}
	return bindDerivedTargets(step, spec, derived)
}

func targetFieldPaths(spec ToolSpec) []string {
	if spec.WriteKind == WriteCreate {
		return spec.ScopeFields
	}
	return spec.TargetFields
}

func bindDerivedTargets(step ApprovedStep, spec ToolSpec, derived []string) (ApprovedStep, []string, error) {
	if spec.WriteKind == WriteCreate {
		if len(step.TargetIDs) > 0 || len(step.ScopeIDs) > 0 && !sameIDs(step.ScopeIDs, derived) {
			return ApprovedStep{}, nil, errors.New("SCOPE_MISMATCH")
		}
		step.ScopeIDs = derived
		return step, nil, nil
	}
	if spec.WriteKind != WriteExisting || len(step.ScopeIDs) > 0 || len(step.TargetIDs) > 0 && !sameIDs(step.TargetIDs, derived) {
		return ApprovedStep{}, nil, errors.New("TARGET_MISMATCH")
	}
	step.TargetIDs = derived
	return step, nil, nil
}

func extractArgumentIDs(arguments map[string]any, paths []string) ([]string, error) {
	var ids []string
	for _, path := range paths {
		value, err := lookupArgumentPath(arguments, path)
		if err != nil {
			return nil, err
		}
		ids, err = appendTargetIDs(ids, value)
		if err != nil {
			return nil, err
		}
	}
	sort.Strings(ids)
	for index := 1; index < len(ids); index++ {
		if ids[index] == ids[index-1] {
			return nil, errors.New("DUPLICATE_STEP_TARGET")
		}
	}
	return ids, nil
}

func lookupArgumentPath(arguments map[string]any, path string) (any, error) {
	var value any = arguments
	for _, key := range strings.Split(path, ".") {
		object, ok := value.(map[string]any)
		if !ok {
			return nil, errors.New("MISSING_STEP_TARGET")
		}
		value = object[key]
	}
	return value, nil
}

func appendTargetIDs(ids []string, value any) ([]string, error) {
	switch typed := value.(type) {
	case string:
		if typed == "" {
			return nil, errors.New("INVALID_STEP_TARGET")
		}
		return append(ids, typed), nil
	case []any:
		for _, item := range typed {
			itemID, ok := item.(string)
			if !ok || itemID == "" {
				return nil, errors.New("INVALID_STEP_TARGET")
			}
			ids = append(ids, itemID)
		}
		return ids, nil
	default:
		return nil, errors.New("INVALID_STEP_TARGET")
	}
}

func sameIDs(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	copyLeft := append([]string(nil), left...)
	sort.Strings(copyLeft)
	for index := range copyLeft {
		if copyLeft[index] != right[index] {
			return false
		}
	}
	return true
}
