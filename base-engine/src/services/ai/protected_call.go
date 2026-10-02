package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

var ErrProtectedCall = errors.New("invalid protected call")

type FixedRequest struct {
	Method, Path, Document string
	Variables              json.RawMessage
	Cookie, Origin         string
}

type FixedResponse struct {
	Status int
	Body   []byte
}

func ProtectedCall(ctx context.Context, handler http.Handler, fixed FixedRequest) (FixedResponse, error) {
	if !validFixedTransport(ctx, handler, fixed) || !fixedPathAllowed(fixed.Method, fixed.Path) {
		return FixedResponse{}, ErrProtectedCall
	}
	body, err := fixedRequestBody(fixed)
	if err != nil {
		return FixedResponse{}, err
	}
	path, err := fixedRequestPath(fixed)
	if err != nil {
		return FixedResponse{}, err
	}
	request, err := http.NewRequestWithContext(ctx, fixed.Method, path, bytes.NewReader(body))
	if err != nil {
		return FixedResponse{}, ErrProtectedCall
	}
	request.Header.Set("Cookie", fixed.Cookie)
	request.Header.Set("Origin", fixed.Origin)
	if fixed.Method != http.MethodGet {
		request.Header.Set("Content-Type", "application/json")
	}
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return FixedResponse{Status: response.Code, Body: bytes.Clone(response.Body.Bytes())}, nil
}

func validFixedTransport(ctx context.Context, handler http.Handler, fixed FixedRequest) bool {
	return ctx != nil && handler != nil && (fixed.Method == http.MethodPost || fixed.Method == http.MethodGet) &&
		fixed.Cookie != "" && fixed.Origin != ""
}

func fixedPathAllowed(method, path string) bool {
	if strings.HasPrefix(path, "/api/ai/") {
		return false
	}
	if method == http.MethodGet {
		return path == "/api/payment-config"
	}
	return path == "/graphql" || path == "/api/franchise-initial-account" ||
		path == "/api/payment-config/state" || path == "/api/payment-config/restore-inheritance"
}

func fixedRequestBody(fixed FixedRequest) ([]byte, error) {
	if len(fixed.Variables) == 0 || !json.Valid(fixed.Variables) {
		return nil, ErrProtectedCall
	}
	if fixed.Method == http.MethodGet {
		if fixed.Document != "" {
			return nil, ErrProtectedCall
		}
		return nil, nil
	}
	if fixed.Path == "/graphql" {
		if strings.TrimSpace(fixed.Document) == "" {
			return nil, ErrProtectedCall
		}
		return json.Marshal(struct {
			Query     string          `json:"query"`
			Variables json.RawMessage `json:"variables"`
		}{Query: fixed.Document, Variables: fixed.Variables})
	}
	if fixed.Document != "" {
		return nil, ErrProtectedCall
	}
	return bytes.Clone(fixed.Variables), nil
}

func fixedRequestPath(fixed FixedRequest) (string, error) {
	if fixed.Method != http.MethodGet {
		return fixed.Path, nil
	}
	var values map[string]string
	if err := json.Unmarshal(fixed.Variables, &values); err != nil || values == nil {
		return "", ErrProtectedCall
	}
	query := url.Values{}
	for key, value := range values {
		if key != "scope" && key != "organizationId" && key != "storeId" {
			return "", ErrProtectedCall
		}
		query.Set(key, value)
	}
	return fixed.Path + "?" + query.Encode(), nil
}
