package ai

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"base-engine/auth"
	"google.golang.org/adk/v2/model"
	"google.golang.org/adk/v2/tool"
)

type ModelProvider func(context.Context) (model.LLM, uint64, error)

type ServiceConfig struct {
	Model                                                   model.LLM
	ModelProvider                                           ModelProvider
	Catalog                                                 *Catalog
	Prompt                                                  string
	ResolvePrincipal                                        func(context.Context, string, time.Time) (*auth.WorkspacePrincipal, error)
	ValidateImageAttachment                                 func(context.Context, *auth.WorkspacePrincipal, string) error
	ModelContextTokens, InputTokenBudget, OutputTokenBudget int
}

type Service struct {
	config           ServiceConfig
	approval         *ApprovalService
	mu               sync.Mutex
	protectedHandler http.Handler
	tools            []tool.Tool
	toolsSet         bool
	activeAccounts   map[string]bool
	rateByAccount    map[string]*accountRate
	now              func() time.Time
}

type accountRate struct {
	hour           int64
	previews, runs int
}

func NewService(config ServiceConfig) (*Service, error) {
	if !validServiceConfig(config) {
		return nil, errors.New("INVALID_AI_SERVICE_CONFIG")
	}
	approval, err := NewApprovalService(config.Catalog)
	if err != nil {
		return nil, err
	}
	return &Service{config: config, approval: approval, activeAccounts: map[string]bool{}, rateByAccount: map[string]*accountRate{}, now: time.Now}, nil
}

func validServiceConfig(config ServiceConfig) bool {
	return (config.Model != nil || config.ModelProvider != nil) && config.Catalog != nil && config.Prompt != "" &&
		config.ResolvePrincipal != nil && config.ModelContextTokens > 0 && config.InputTokenBudget > 0 &&
		config.OutputTokenBudget > 0
}

func (s *Service) SetProtectedHandler(handler http.Handler) error {
	if s == nil || handler == nil {
		return errors.New("INVALID_PROTECTED_HANDLER")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.protectedHandler != nil {
		return errors.New("PROTECTED_HANDLER_ALREADY_SET")
	}
	s.protectedHandler = handler
	return nil
}

func (s *Service) SetTools(items []tool.Tool) error {
	if s == nil {
		return errors.New("INVALID_AI_SERVICE")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.toolsSet {
		return errors.New("TOOLS_ALREADY_SET")
	}
	if len(items) != len(s.config.Catalog.specs) {
		return errors.New("INCOMPLETE_TOOL_REGISTRATION")
	}
	seen := make(map[string]struct{}, len(items))
	for _, item := range items {
		if item == nil {
			return errors.New("NIL_TOOL")
		}
		name := item.Name()
		if _, ok := s.config.Catalog.Lookup(name); !ok {
			return errors.New("UNREGISTERED_TOOL")
		}
		if _, exists := seen[name]; exists {
			return errors.New("DUPLICATE_TOOL_REGISTRATION")
		}
		seen[name] = struct{}{}
	}
	s.tools = append([]tool.Tool(nil), items...)
	s.toolsSet = true
	return nil
}

func (s *Service) Ready() bool {
	if s == nil {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.protectedHandler != nil && s.toolsSet
}
