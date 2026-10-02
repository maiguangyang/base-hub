package integration_test

import (
	"strings"
	"testing"
)

func TestContractGapOutputShapeRejectsExtraFields(t *testing.T) {
	actual := map[string]any{"id": "account-1", "phone": "010", "displayName": "A", "email": nil, "status": "ACTIVE", "secret": "leaked"}
	if sameGapShape(actual, gapExpectedShape("HqAccount", "account")) {
		t.Fatal("extra account field accepted")
	}
}

func gapExpectedShape(toolID, root string) any {
	account := gapFields("id phone displayName email status")
	invitation := gapFields("id membershipId invitedByAccountId expiresAt acceptedAt revokedAt createdAt")
	switch root {
	case "accounts":
		return gapPage(account)
	case "membershipInvitations":
		return gapPage(invitation)
	default:
		return gapEntityShape(toolID, root, account, invitation)
	}
}

func gapEntityShape(toolID, root string, account, invitation map[string]any) any {
	switch root {
	case "account":
		return account
	case "auditLog":
		return gapFields("id action resourceType resourceId resultCode actorAccountId organizationId storeId createdAt")
	case "membershipInvitation":
		return invitation
	case "operatorMembership":
		return gapMembershipShape()
	case "operatorRole":
		return gapRoleShape()
	case "permission":
		return gapFields("id name action module scope")
	case "store":
		return gapStoreShape(toolID)
	default:
		return nil
	}
}

func gapMembershipShape() map[string]any {
	shape := gapFields("id status storeAccessMode accountId organizationId roles stores")
	shape["roles"] = []any{gapFields("id name kind")}
	shape["stores"] = []any{gapFields("id name lifecycle")}
	return shape
}

func gapRoleShape() map[string]any {
	shape := gapFields("id name kind organizationId permissions")
	shape["permissions"] = []any{gapFields("id name action module scope")}
	return shape
}

func gapStoreShape(toolID string) map[string]any {
	fields := "id code name lifecycle organizationId contactPhone managerName managerPhone province city district address"
	if toolID != "HqDirectStore" {
		fields += " businessHours businessStatus"
	}
	return gapFields(fields)
}

func gapFields(names string) map[string]any {
	fields := map[string]any{}
	for _, name := range strings.Fields(names) {
		fields[name] = nil
	}
	return fields
}

func gapPage(item any) map[string]any {
	shape := gapFields("data total current_page per_page total_page")
	shape["data"] = []any{item}
	return shape
}

func sameGapShape(actual, expected any) bool {
	switch shape := expected.(type) {
	case map[string]any:
		return sameGapObject(actual, shape)
	case []any:
		return sameGapItems(actual, shape[0])
	}
	return true
}

func sameGapObject(actual any, shape map[string]any) bool {
	object, ok := actual.(map[string]any)
	if !ok || len(object) != len(shape) {
		return false
	}
	for name, child := range shape {
		value, exists := object[name]
		if !exists || !sameGapShape(value, child) {
			return false
		}
	}
	return true
}

func sameGapItems(actual, itemShape any) bool {
	items, ok := actual.([]any)
	if !ok {
		return false
	}
	for _, item := range items {
		if !sameGapShape(item, itemShape) {
			return false
		}
	}
	return true
}
