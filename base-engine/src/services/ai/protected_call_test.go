package ai

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestProtectedCallUsesServerDocument(t *testing.T) {
	var received struct {
		Query     string         `json:"query"`
		Variables map[string]any `json:"variables"`
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatal(err)
		}
		w.WriteHeader(http.StatusAccepted)
	})
	response, err := ProtectedCall(t.Context(), handler, FixedRequest{
		Method: http.MethodPost, Path: "/graphql", Document: `query { viewer { account { id } } }`,
		Variables: json.RawMessage(`{"query":"mutation { deleteStores(id:[\"store-a\"]) }"}`),
		Cookie:    "session=test", Origin: "https://admin.example.com",
	})
	if err != nil || response.Status != http.StatusAccepted || received.Query != `query { viewer { account { id } } }` || received.Variables["query"] == nil {
		t.Fatalf("fixed document lost: status=%d query=%q variables=%v err=%v", response.Status, received.Query, received.Variables, err)
	}
}

func TestProtectedCallRejectsAIAndArbitraryPaths(t *testing.T) {
	called := false
	handler := http.HandlerFunc(func(http.ResponseWriter, *http.Request) { called = true })
	for _, path := range []string{"/api/ai/preview", "/api/ai/run", "/api/ai/model-config", "/graphql?query=attack", "/api/system-initialization"} {
		_, err := ProtectedCall(context.Background(), handler, FixedRequest{
			Method: http.MethodPost, Path: path, Document: `query { viewer { account { id } } }`,
			Variables: json.RawMessage(`{}`), Cookie: "session=test", Origin: "https://admin.example.com",
		})
		if !errors.Is(err, ErrProtectedCall) || called {
			t.Fatalf("path %q accepted: called=%t err=%v", path, called, err)
		}
	}
}

func TestProtectedCallForwardsFixedHTTPCommand(t *testing.T) {
	var path, cookie, origin string
	var input map[string]any
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path, cookie, origin = r.URL.Path, r.Header.Get("Cookie"), r.Header.Get("Origin")
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			t.Fatal(err)
		}
		w.WriteHeader(http.StatusAccepted)
	})
	response, err := ProtectedCall(t.Context(), handler, FixedRequest{
		Method: http.MethodPost, Path: "/api/franchise-initial-account", Variables: json.RawMessage(`{"organizationId":"org-a"}`),
		Cookie: "session=test", Origin: "https://admin.example.com",
	})
	if err != nil || response.Status != http.StatusAccepted || path != "/api/franchise-initial-account" ||
		cookie != "session=test" || origin != "https://admin.example.com" || input["organizationId"] != "org-a" {
		t.Fatalf("fixed HTTP command status=%d path=%q cookie=%q origin=%q input=%v err=%v", response.Status, path, cookie, origin, input, err)
	}
}

func TestProtectedCallPaymentReadUsesFixedGETAndQuery(t *testing.T) {
	var method, path, scope, storeID, cookie string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path, scope, storeID, cookie = r.Method, r.URL.Path, r.URL.Query().Get("scope"), r.URL.Query().Get("storeId"), r.Header.Get("Cookie")
		w.WriteHeader(http.StatusOK)
	})
	response, err := ProtectedCall(t.Context(), handler, FixedRequest{Method: http.MethodGet, Path: "/api/payment-config",
		Variables: json.RawMessage(`{"scope":"STORE","storeId":"store-a"}`), Cookie: "session=test", Origin: "https://admin.example.com"})
	if err != nil || response.Status != http.StatusOK || method != http.MethodGet || path != "/api/payment-config" || scope != "STORE" || storeID != "store-a" || cookie != "session=test" {
		t.Fatalf("payment read status=%d method=%s path=%s scope=%s store=%s cookie=%s err=%v", response.Status, method, path, scope, storeID, cookie, err)
	}
}

func TestProtectedCallPaymentWriteOnlyAllowsReviewedRoutes(t *testing.T) {
	called := 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called++; w.WriteHeader(http.StatusOK) })
	for _, path := range []string{"/api/payment-config/state", "/api/payment-config/restore-inheritance"} {
		_, err := ProtectedCall(t.Context(), handler, FixedRequest{Method: http.MethodPost, Path: path,
			Variables: json.RawMessage(`{"scope":"GLOBAL","channel":"WECHAT","recordId":"one","version":1}`), Cookie: "session=test", Origin: "https://admin.example.com"})
		if err != nil {
			t.Fatalf("reviewed route %s rejected: %v", path, err)
		}
	}
	_, err := ProtectedCall(t.Context(), handler, FixedRequest{Method: http.MethodPut, Path: "/api/payment-config",
		Variables: json.RawMessage(`{}`), Cookie: "session=test", Origin: "https://admin.example.com"})
	if !errors.Is(err, ErrProtectedCall) || called != 2 {
		t.Fatalf("credential save exposed: called=%d err=%v", called, err)
	}
}
