package ai

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"

	"base-engine/auth"
	"base-engine/src/services/audit"
	"google.golang.org/adk/v2/agent"
)

func (s *Service) FixedToolRuntime() FixedToolRuntime {
	return FixedToolRuntime{Call: s.callFixedTool, DeliverSecret: s.deliverSecret}
}

func (s *Service) callFixedTool(ctx agent.Context, proposed ToolSpec, variables json.RawMessage) (FixedResponse, error) {
	if err := ctx.Err(); err != nil {
		return FixedResponse{}, err
	}
	state, principal, err := s.fixedToolPrincipal(ctx)
	if err != nil {
		return FixedResponse{}, err
	}
	spec, ok := s.config.Catalog.Lookup(proposed.ID)
	if !ok || !sameFixedSpec(spec, proposed) || !specAllowed(spec, principal, state.phase, state.approvedIDs) {
		return FixedResponse{}, errors.New("AI_TOOL_NOT_AUTHORIZED")
	}
	if err := state.checkModelCurrent(ctx); err != nil {
		return FixedResponse{}, err
	}
	arguments, err := state.preparedFixedArguments(spec, variables)
	if err != nil {
		return FixedResponse{}, err
	}
	if err := state.validateFixedImageAttachment(ctx, principal, spec, arguments); err != nil {
		return FixedResponse{}, err
	}
	path, document, method := fixedToolTransport(spec)
	protectedContext := audit.WithAIInvocation(ctx, state.runID, spec.ID)
	return ProtectedCall(protectedContext, s.protectedHandler, FixedRequest{Method: method, Path: path, Document: document,
		Variables: arguments, Cookie: auth.SessionCookieName + "=" + state.token, Origin: state.origin})
}

func (s *Service) fixedToolPrincipal(ctx agent.Context) (*runState, *auth.WorkspacePrincipal, error) {
	state, err := stateFromContext(ctx)
	if err != nil || state.service != s {
		return nil, nil, errors.New("INVALID_TOOL_RUN")
	}
	principal, err := state.checkPrincipal(ctx)
	return state, principal, err
}

func (state *runState) validateFixedImageAttachment(ctx agent.Context, principal *auth.WorkspacePrincipal, spec ToolSpec, arguments json.RawMessage) error {
	if spec.ID != productMainImageToolID {
		return nil
	}
	if err := validateImageAttachmentStep(ApprovedStep{ToolID: spec.ID, Arguments: arguments}, state.attachmentID); err != nil {
		return err
	}
	if state.service.config.ValidateImageAttachment == nil {
		return errors.New("IMAGE_ATTACHMENT_UNAVAILABLE")
	}
	return state.service.config.ValidateImageAttachment(ctx, principal, state.attachmentID)
}

func (state *runState) preparedFixedArguments(spec ToolSpec, variables json.RawMessage) (json.RawMessage, error) {
	arguments, err := state.checkedToolArguments(spec, variables)
	if err != nil {
		return nil, err
	}
	return injectTrustedFields(spec, state.planID, arguments)
}

func fixedToolTransport(spec ToolSpec) (path, document, method string) {
	if spec.Path != "" {
		return spec.Path, "", spec.Method
	}
	return "/graphql", spec.Document, http.MethodPost
}

func (state *runState) checkedToolArguments(spec ToolSpec, variables json.RawMessage) (json.RawMessage, error) {
	arguments, err := state.approvedToolArguments(spec, variables)
	if err != nil {
		return nil, err
	}
	if err := validateReadArguments(spec, arguments); err != nil {
		return nil, err
	}
	return arguments, nil
}

func sameFixedSpec(left, right ToolSpec) bool {
	left.InputSchema = nil
	right.InputSchema = nil
	return reflect.DeepEqual(left, right)
}

func (state *runState) approvedToolArguments(spec ToolSpec, variables json.RawMessage) (json.RawMessage, error) {
	if spec.Mode != ModeWrite {
		return variables, nil
	}
	if state.phase != PhaseRun {
		return nil, errors.New("WRITE_IN_PREVIEW")
	}
	bound, err := state.service.approval.bindAttestationArguments(state.planID, spec, variables)
	if err != nil {
		return nil, err
	}
	if err := state.service.approval.AuthorizeStep(state.planID, spec.ID, bound); err != nil {
		return nil, err
	}
	return bound, nil
}

func (s *Service) deliverSecret(ctx agent.Context, payload SecretPayload) error {
	state, err := stateFromContext(ctx)
	if err != nil || state.service != s || state.phase != PhaseRun {
		return errors.New("INVALID_SECRET_DELIVERY")
	}
	return state.sink.Emit("secret", map[string]string{"toolId": payload.ToolID, "targetAccountId": payload.TargetAccountID, "value": payload.Value})
}
