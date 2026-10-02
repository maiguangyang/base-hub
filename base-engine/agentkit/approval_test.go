package agentkit

import (
	"encoding/json"
	"testing"
)

func TestApprovalGuardExecutesOnlyApprovedArgumentsInOrder(t *testing.T) {
	guard, err := NewApprovalGuard([]ApprovedStep{
		{ToolID: "first", Arguments: json.RawMessage(`{"id":"a"}`), MaxCalls: 1},
		{ToolID: "second", Arguments: json.RawMessage(`{"id":"b"}`), MaxCalls: 1},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := guard.Reserve("second", json.RawMessage(`{"id":"b"}`)); err == nil {
		t.Fatal("out-of-order step accepted")
	}
	if err := guard.Reserve("first", json.RawMessage(`{"id":"other"}`)); err == nil {
		t.Fatal("changed arguments accepted")
	}
	if err := guard.Reserve("first", json.RawMessage(`{"id":"a"}`)); err != nil {
		t.Fatal(err)
	}
	if err := guard.RecordSuccess("first"); err != nil {
		t.Fatal(err)
	}
	if err := guard.Complete(); err == nil {
		t.Fatal("incomplete plan accepted")
	}
	if err := guard.Reserve("second", json.RawMessage(`{"id":"b"}`)); err != nil {
		t.Fatal(err)
	}
	if err := guard.RecordSuccess("second"); err != nil {
		t.Fatal(err)
	}
	if err := guard.Complete(); err != nil {
		t.Fatal(err)
	}
}

func TestApprovalGuardDistinguishesLargeJSONIntegers(t *testing.T) {
	guard, err := NewApprovalGuard([]ApprovedStep{{
		ToolID: "write", Arguments: json.RawMessage(`{"input":{"record":9007199254740992}}`), MaxCalls: 1,
	}})
	if err != nil {
		t.Fatal(err)
	}
	if err := guard.Reserve("write", json.RawMessage(`{"input":{"record":9007199254740993}}`)); err == nil {
		t.Fatal("a distinct large integer matched the approved arguments")
	}
	if err := guard.Reserve("write", json.RawMessage(`{"input":{"record":9007199254740992}}`)); err != nil {
		t.Fatalf("approved integer was rejected: %v", err)
	}
}
