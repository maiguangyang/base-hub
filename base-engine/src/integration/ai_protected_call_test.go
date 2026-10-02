package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/ai"
)

func TestAIProtectedCallParity(t *testing.T) {
	fixture := newSecurityFixture(t)
	document := `query { organizations { total data { id } } }`
	request := ai.FixedRequest{
		Method: http.MethodPost, Path: "/graphql", Document: document,
		Variables: json.RawMessage(`{}`),
		Cookie:    fixture.cookies["session-hq"].String(), Origin: "https://admin.example.com",
	}
	directBody, err := json.Marshal(map[string]any{"query": document, "variables": json.RawMessage(`{}`)})
	if err != nil {
		t.Fatal(err)
	}
	directRequest := httptest.NewRequest(http.MethodPost, "/graphql", bytes.NewReader(directBody))
	directRequest.Header.Set("Content-Type", "application/json")
	directRequest.Header.Set("Origin", request.Origin)
	directRequest.Header.Set("Cookie", request.Cookie)
	direct := httptest.NewRecorder()
	fixture.handler.ServeHTTP(direct, directRequest)
	protected, err := ai.ProtectedCall(context.Background(), fixture.handler, request)
	if err != nil {
		t.Fatal(err)
	}
	if protected.Status != direct.Code || string(protected.Body) != direct.Body.String() {
		t.Fatalf("protected response = %d %s; direct = %d %s", protected.Status, protected.Body, direct.Code, direct.Body.String())
	}
}

func TestAIProtectedCallParityAcrossWorkspaceAndWrites(t *testing.T) {
	cases := []parityScenario{
		{name: "franchise scoped read", session: "session-a-1", document: `query { stores { total data { id } } }`},
		{name: "cross tenant denied", session: "session-a-1", document: `query { store(id: "store-b") { id } }`},
		{name: "franchise store update", session: "session-a-1", document: `mutation { updateStore(id: "store-a-draft", input: {name: "Changed"}) { id name } }`, auditAction: "store:update", storeID: "store-a-draft"},
		{name: "headquarters suspension", session: "session-hq", document: `mutation { suspendOrganization(input: {organizationId: "org-a", reasonCode: "CONTRACT_ENDED"}) { id status } }`, auditAction: "organization:suspend", organizationID: "org-a"},
		{name: "revoked session", session: "session-a-1", document: `query { stores { total } }`, revoke: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) { checkParityScenario(t, tc) })
	}
}

type parityScenario struct {
	name, session, document, auditAction, storeID, organizationID string
	revoke                                                        bool
}

func checkParityScenario(t *testing.T, tc parityScenario) {
	t.Helper()
	var direct, protected paritySnapshot
	t.Run("ordinary", func(t *testing.T) {
		direct = runParityCall(t, tc.session, tc.document, tc.auditAction, tc.storeID, tc.organizationID, tc.revoke, false)
	})
	t.Run("protected", func(t *testing.T) {
		protected = runParityCall(t, tc.session, tc.document, tc.auditAction, tc.storeID, tc.organizationID, tc.revoke, true)
	})
	if !reflect.DeepEqual(protected, direct) {
		t.Fatalf("protected = %#v, ordinary = %#v", protected, direct)
	}
	assertParityWrite(t, tc, direct)
	assertParityAccess(t, tc, direct)
}

func assertParityWrite(t *testing.T, tc parityScenario, direct paritySnapshot) {
	t.Helper()
	if tc.auditAction != "" && (direct.AuditCount != 1 || direct.AuditAccountID != "account-shared" || direct.AuditSessionID != tc.session) {
		t.Fatalf("audit identity = %#v", direct)
	}
	if tc.storeID != "" && direct.StoreName != "Changed" {
		t.Fatalf("store update = %#v", direct)
	}
	if tc.organizationID != "" && direct.Organization != string(gen.OrganizationStatusSuspended) {
		t.Fatalf("organization governance = %#v", direct)
	}
}

func assertParityAccess(t *testing.T, tc parityScenario, direct paritySnapshot) {
	t.Helper()
	if tc.revoke && direct.ErrorCode != string(auth.CodeSessionRevoked) {
		t.Fatalf("revoked session code = %s", direct.ErrorCode)
	}
	if tc.name == "cross tenant denied" && direct.ErrorCode != string(auth.CodePermissionDenied) {
		t.Fatalf("cross-tenant code = %s", direct.ErrorCode)
	}
	if tc.name == "franchise scoped read" {
		data, _ := json.Marshal(direct.Data)
		if bytes.Contains(data, []byte("store-b")) || bytes.Contains(data, []byte("store-hq")) {
			t.Fatalf("franchise data escaped scope: %s", data)
		}
	}
}

type paritySnapshot struct {
	Status         int
	Data           any
	ErrorCode      string
	StoreName      string
	Organization   string
	AuditCount     int
	AuditAccountID string
	AuditSessionID string
}

func runParityCall(t *testing.T, sessionID, document, action, storeID, organizationID string, revoke, protected bool) paritySnapshot {
	t.Helper()
	fixture := newSecurityFixture(t)
	if storeID == "store-a-draft" {
		fixture.createAll([]gen.Store{{ID: storeID, Code: "A-DRAFT", Name: "Draft", Lifecycle: gen.StoreLifecycleDraft, OrganizationID: "org-a"}})
	}
	if revoke {
		if err := fixture.db.Model(&gen.Session{}).Where("id = ?", sessionID).Update("revoked_at", time.Now()).Error; err != nil {
			t.Fatal(err)
		}
	}
	fixed := ai.FixedRequest{Method: http.MethodPost, Path: "/graphql", Document: document,
		Variables: json.RawMessage(`{}`), Cookie: fixture.cookies[sessionID].String(), Origin: "https://admin.example.com"}
	var response ai.FixedResponse
	if protected {
		var err error
		response, err = ai.ProtectedCall(t.Context(), fixture.handler, fixed)
		if err != nil {
			t.Fatal(err)
		}
	} else {
		body, err := json.Marshal(map[string]any{"query": document, "variables": json.RawMessage(`{}`)})
		if err != nil {
			t.Fatal(err)
		}
		request := httptest.NewRequest(http.MethodPost, "/graphql", bytes.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Origin", fixed.Origin)
		request.Header.Set("Cookie", fixed.Cookie)
		recorder := httptest.NewRecorder()
		fixture.handler.ServeHTTP(recorder, request)
		response = ai.FixedResponse{Status: recorder.Code, Body: recorder.Body.Bytes()}
	}
	return snapshotParityResult(t, fixture, response, action, storeID, organizationID)
}

func snapshotParityResult(t *testing.T, fixture *securityFixture, response ai.FixedResponse, action, storeID, organizationID string) paritySnapshot {
	t.Helper()
	var result struct {
		Data   any `json:"data"`
		Errors []struct {
			Extensions map[string]any `json:"extensions"`
		} `json:"errors"`
	}
	if err := json.Unmarshal(response.Body, &result); err != nil {
		t.Fatal(err)
	}
	snapshot := paritySnapshot{Status: response.Status, Data: result.Data}
	if len(result.Errors) > 0 {
		snapshot.ErrorCode, _ = result.Errors[0].Extensions["code"].(string)
	}
	if storeID != "" {
		var store gen.Store
		if err := fixture.db.First(&store, "id = ?", storeID).Error; err != nil {
			t.Fatal(err)
		}
		snapshot.StoreName = store.Name
	}
	if organizationID != "" {
		var organization gen.Organization
		if err := fixture.db.First(&organization, "id = ?", organizationID).Error; err != nil {
			t.Fatal(err)
		}
		snapshot.Organization = string(organization.Status)
	}
	if action != "" {
		snapshotParityAudit(t, fixture, action, response.Body, &snapshot)
	}
	return snapshot
}

func snapshotParityAudit(t *testing.T, fixture *securityFixture, action string, body []byte, snapshot *paritySnapshot) {
	t.Helper()
	var audits []gen.AuditLog
	if err := fixture.db.Where("action = ?", action).Find(&audits).Error; err != nil {
		t.Fatal(err)
	}
	snapshot.AuditCount = len(audits)
	if len(audits) != 1 {
		t.Fatalf("%s audit count = %d, response = %s", action, len(audits), body)
	}
	if audits[0].ActorAccountID != nil {
		snapshot.AuditAccountID = *audits[0].ActorAccountID
	}
	if audits[0].SessionID != nil {
		snapshot.AuditSessionID = *audits[0].SessionID
	}
}
