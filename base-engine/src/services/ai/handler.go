package ai

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"base-engine/auth"
	"base-engine/config"
)

type previewRequest struct {
	Prompt       string           `json:"prompt"`
	History      []chatTurn       `json:"history,omitempty"`
	Attestation  *UserAttestation `json:"attestation,omitempty"`
	AttachmentID string           `json:"attachmentId,omitempty"`
}

type chatTurn struct {
	Role string `json:"role"`
	Text string `json:"text"`
}
type runRequest struct {
	PreviewToken string `json:"previewToken"`
}

func (s *Service) RegisterRoutes(router *mux.Router, security config.SecurityConfig) {
	router.HandleFunc("/api/ai/preview", func(w http.ResponseWriter, r *http.Request) { s.handlePreview(w, r, security) }).Methods(http.MethodPost).Name(HTTPRouteName("http.ai.preview", previewRequest{}, AIEvent{}))
	router.HandleFunc("/api/ai/run", func(w http.ResponseWriter, r *http.Request) { s.handleRun(w, r, security) }).Methods(http.MethodPost).Name(HTTPRouteName("http.ai.run", runRequest{}, AIEvent{}))
}

func (s *Service) handlePreview(w http.ResponseWriter, r *http.Request, security config.SecurityConfig) {
	principal, token, err := s.authorizeAIRequest(r, security)
	if err != nil {
		writeAIHTTPError(w, err)
		return
	}
	var input previewRequest
	if err := decodeAIRequest(w, r, &input); err != nil {
		writeAIHTTPError(w, err)
		return
	}
	input.Prompt = strings.TrimSpace(input.Prompt)
	if input.Prompt == "" || len(input.Prompt) > 16000 || !validChatHistory(input.History) {
		writeAIHTTPError(w, errAIBadRequest)
		return
	}
	release, err := s.acquire(principal.AccountID, PhasePreview)
	if err != nil {
		writeAIHTTPError(w, err)
		return
	}
	s.servePhase(w, r, principal, token, PhasePreview, release, func(ctx context.Context, sink *eventSink) error {
		return s.runPreview(ctx, principal, token, r.Header.Get("Origin"), input, sink)
	})
}

func (s *Service) handleRun(w http.ResponseWriter, r *http.Request, security config.SecurityConfig) {
	principal, token, err := s.authorizeAIRequest(r, security)
	if err != nil {
		writeAIHTTPError(w, err)
		return
	}
	var input runRequest
	if err := decodeAIRequest(w, r, &input); err != nil {
		writeAIHTTPError(w, err)
		return
	}
	if input.PreviewToken == "" {
		writeAIHTTPError(w, errAIBadRequest)
		return
	}
	release, err := s.acquire(principal.AccountID, PhaseRun)
	if err != nil {
		writeAIHTTPError(w, err)
		return
	}
	plan, err := s.approval.Consume(input.PreviewToken, principal)
	if err != nil {
		release()
		writeAIHTTPError(w, err)
		return
	}
	s.servePhase(w, r, principal, token, PhaseRun, release, func(ctx context.Context, sink *eventSink) error {
		return s.runApproved(ctx, principal, token, r.Header.Get("Origin"), plan, sink)
	})
}

var errAIBadRequest = errors.New("INVALID_AI_REQUEST")

func (s *Service) authorizeAIRequest(r *http.Request, security config.SecurityConfig) (*auth.WorkspacePrincipal, string, error) {
	if !s.Ready() {
		return nil, "", errors.New("AI_NOT_READY")
	}
	origin := strings.TrimSpace(r.Header.Get("Origin"))
	if origin == "" || !security.OriginAllowed(origin) {
		return nil, "", errors.New("ORIGIN_NOT_ALLOWED")
	}
	if r.Header.Get("Content-Type") != "application/json" {
		return nil, "", errors.New("UNSUPPORTED_MEDIA_TYPE")
	}
	cookie, err := r.Cookie(auth.SessionCookieName)
	if err != nil || cookie.Value == "" {
		return nil, "", errors.New("AUTH_REQUIRED")
	}
	principal, err := s.config.ResolvePrincipal(r.Context(), cookie.Value, s.now())
	if err != nil || principal == nil || !adminWorkspace(principal.WorkspaceType) {
		return nil, "", errors.New("AUTH_REQUIRED")
	}
	return principal, cookie.Value, nil
}

func adminWorkspace(workspace auth.WorkspaceType) bool {
	return workspace == auth.WorkspaceTypeHeadquarters || workspace == auth.WorkspaceTypeFranchise
}

func writeAIHTTPError(w http.ResponseWriter, err error) {
	status := http.StatusBadRequest
	switch err.Error() {
	case "AI_NOT_READY":
		status = http.StatusServiceUnavailable
	case "ORIGIN_NOT_ALLOWED":
		status = http.StatusForbidden
	case "UNSUPPORTED_MEDIA_TYPE":
		status = http.StatusUnsupportedMediaType
	case "AUTH_REQUIRED":
		status = http.StatusUnauthorized
	case "AI_RUN_ALREADY_ACTIVE":
		status = http.StatusConflict
	case "AI_PREVIEW_RATE_LIMIT", "AI_RUN_RATE_LIMIT":
		status = http.StatusTooManyRequests
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"code": err.Error()})
}
