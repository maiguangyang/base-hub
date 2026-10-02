/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package audit

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"base-engine/gen"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestAuditAIInvocationMetadataUsesServerContext(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{IgnoreRelationshipsWhenMigrating: true})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&gen.AuditLog{}); err != nil {
		t.Fatal(err)
	}
	ctx := WithAIInvocation(context.Background(), "run-1", "HqCreateRole")
	if err := NewService().Write(db.WithContext(ctx), Entry{ActorAccountID: "account-1", Action: "role:create", ResourceType: "role", ResourceID: "role-1", ResultCode: "SUCCESS"}); err != nil {
		t.Fatal(err)
	}
	var record gen.AuditLog
	if err := db.First(&record).Error; err != nil {
		t.Fatal(err)
	}
	if record.ActorAccountID == nil || *record.ActorAccountID != "account-1" {
		t.Fatal("actor changed")
	}
	var metadata map[string]any
	if err := json.Unmarshal([]byte(*record.MetadataJSON), &metadata); err != nil {
		t.Fatal(err)
	}
	if metadata["runId"] != "run-1" || metadata["toolId"] != "HqCreateRole" {
		t.Fatalf("AI metadata = %+v", metadata)
	}
}

// TestWriteUsesAllowlistedMetadata 验证审计只序列化显式结构字段。
func TestWriteUsesAllowlistedMetadata(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{IgnoreRelationshipsWhenMigrating: true})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&gen.AuditLog{}); err != nil {
		t.Fatal(err)
	}
	service := NewService()
	entry := Entry{ActorAccountID: "account-1", Action: "store:submit", ResourceType: "store", ResourceID: "store-1", ResultCode: "SUCCESS", Metadata: Metadata{
		ActorMembershipID: "membership-1", RequestID: "request-1", ReasonCode: "READY", TargetStatus: "PENDING_APPROVAL",
		RoleIDs: []string{"role-1"}, PermissionIDs: []string{"permission-1"},
	}}
	if err := service.Write(db, entry); err != nil {
		t.Fatal(err)
	}
	var record gen.AuditLog
	if err := db.First(&record).Error; err != nil {
		t.Fatal(err)
	}
	assertSafeMetadata(t, record.MetadataJSON)
	var metadata map[string]any
	if err := json.Unmarshal([]byte(*record.MetadataJSON), &metadata); err != nil {
		t.Fatal(err)
	}
	assertMetadataString(t, metadata, "actorMembershipId", "membership-1")
	assertMetadataString(t, metadata, "requestId", "request-1")
	assertMetadataList(t, metadata, "roleIds", "role-1")
	assertMetadataList(t, metadata, "permissionIds", "permission-1")
}

func assertSafeMetadata(t *testing.T, metadataJSON *string) {
	t.Helper()
	if metadataJSON == nil || strings.Contains(*metadataJSON, "password") || strings.Contains(*metadataJSON, "token") {
		t.Fatalf("unsafe metadata: %#v", metadataJSON)
	}
}

func assertMetadataString(t *testing.T, metadata map[string]any, key, expected string) {
	t.Helper()
	if metadata[key] != expected {
		t.Fatalf("metadata %s = %#v, want %q", key, metadata[key], expected)
	}
}

func assertMetadataList(t *testing.T, metadata map[string]any, key, expected string) {
	t.Helper()
	values, ok := metadata[key].([]any)
	if !ok || len(values) != 1 || values[0] != expected {
		t.Fatalf("metadata %s = %#v, want [%q]", key, metadata[key], expected)
	}
}

// TestWriteAllowsAnonymousActor 验证认证失败等匿名事件不会制造伪造账号外键。
func TestWriteAllowsAnonymousActor(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{IgnoreRelationshipsWhenMigrating: true})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&gen.AuditLog{}); err != nil {
		t.Fatal(err)
	}
	if err := NewService().Write(db, Entry{Action: "authentication:login", ResourceType: "account", ResultCode: "INVALID_CREDENTIALS"}); err != nil {
		t.Fatal(err)
	}
	var record gen.AuditLog
	if err := db.First(&record).Error; err != nil {
		t.Fatal(err)
	}
	if record.ActorAccountID != nil || record.ResourceID != nil {
		t.Fatalf("anonymous audit contains identifiers: %#v", record)
	}
}
