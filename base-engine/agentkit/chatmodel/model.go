package chatmodel

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"iter"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"google.golang.org/adk/v2/model"
)

var (
	ErrModelConfig      = errors.New("invalid AI model configuration")
	ErrModelProtocol    = errors.New("AI model protocol error")
	ErrModelUnavailable = errors.New("AI model unavailable")
	errStreamStopped    = errors.New("stream consumer stopped")
)

type ModelConfig struct {
	Name    string
	BaseURL string
	APIKey  string
}

type chatModel struct {
	name     string
	apiKey   string
	endpoint string
	client   *http.Client
}

type upstreamHTTPError struct{ statusCode int }

func (e *upstreamHTTPError) Error() string {
	return fmt.Sprintf("%s: upstream HTTP %d", ErrModelUnavailable, e.statusCode)
}
func (e *upstreamHTTPError) Unwrap() error           { return ErrModelUnavailable }
func (e *upstreamHTTPError) UpstreamStatusCode() int { return e.statusCode }

func NewModel(ctx context.Context, cfg ModelConfig) (model.LLM, error) {
	return NewModelWithClient(ctx, cfg, nil)
}

// NewModelWithClient uses the supplied HTTP transport while enforcing the same redirect policy.
func NewModelWithClient(_ context.Context, cfg ModelConfig, supplied *http.Client) (model.LLM, error) {
	parsed, err := validateModelConfig(cfg)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	if supplied != nil {
		*client = *supplied
	}
	if client.Timeout == 0 {
		client.Timeout = 2 * time.Minute
	}
	client.CheckRedirect = sameOriginRedirect(parsed)
	return &chatModel{
		name: cfg.Name, apiKey: cfg.APIKey,
		endpoint: strings.TrimRight(cfg.BaseURL, "/") + "/chat/completions",
		client:   client,
	}, nil
}

func newModelWithClient(ctx context.Context, cfg ModelConfig, supplied *http.Client) (model.LLM, error) {
	return NewModelWithClient(ctx, cfg, supplied)
}

func validateModelConfig(cfg ModelConfig) (*url.URL, error) {
	for _, field := range []struct{ name, value string }{
		{"name", cfg.Name}, {"base_url", cfg.BaseURL}, {"api_key", cfg.APIKey},
	} {
		if strings.TrimSpace(field.value) == "" {
			return nil, fmt.Errorf("%w: %s is required", ErrModelConfig, field.name)
		}
	}
	return validateModelURL(cfg.BaseURL)
}

func validateModelURL(raw string) (*url.URL, error) {
	parsed, err := url.Parse(raw)
	if err != nil || parsed == nil || parsed.Host == "" || parsed.User != nil ||
		parsed.RawQuery != "" || parsed.Fragment != "" {
		return nil, fmt.Errorf("%w: base_url", ErrModelConfig)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, fmt.Errorf("%w: base_url", ErrModelConfig)
	}
	return parsed, nil
}

func sameOriginRedirect(origin *url.URL) func(*http.Request, []*http.Request) error {
	return func(next *http.Request, via []*http.Request) error {
		if len(via) >= 10 || normalizedOrigin(next.URL) != normalizedOrigin(origin) {
			return ErrModelProtocol
		}
		return nil
	}
}

func normalizedOrigin(value *url.URL) string {
	port := value.Port()
	if port == "" && value.Scheme == "http" {
		port = "80"
	}
	if port == "" && value.Scheme == "https" {
		port = "443"
	}
	return strings.ToLower(value.Scheme) + "://" + net.JoinHostPort(strings.ToLower(value.Hostname()), port)
}

func (m *chatModel) Name() string { return m.name }

func (m *chatModel) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		wire, declared, err := buildChatRequest(m.name, req, stream)
		if err != nil {
			yield(nil, err)
			return
		}
		response, err := m.doRequest(ctx, wire)
		if err != nil {
			yield(nil, err)
			return
		}
		defer response.Body.Close()
		if response.StatusCode < 200 || response.StatusCode >= 300 {
			_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, 64*1024))
			yield(nil, &upstreamHTTPError{statusCode: response.StatusCode})
			return
		}
		if stream {
			err = readChatStream(ctx, response.Body, declared, yield)
		} else {
			err = readChatCompletion(response.Body, declared, yield)
		}
		if err != nil && !errors.Is(err, errStreamStopped) {
			yield(nil, err)
		}
	}
}

func (m *chatModel) doRequest(ctx context.Context, wire chatRequest) (*http.Response, error) {
	body, err := json.Marshal(wire)
	if err != nil {
		return nil, ErrModelProtocol
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, m.endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, ErrModelProtocol
	}
	req.Header.Set("Authorization", "Bearer "+m.apiKey)
	req.Header.Set("Content-Type", "application/json")
	if wire.Stream {
		req.Header.Set("Accept", "text/event-stream")
	}
	response, err := m.client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, ErrModelUnavailable
	}
	return response, nil
}
