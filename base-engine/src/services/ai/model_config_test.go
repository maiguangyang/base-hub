package ai

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"base-engine/auth"
	"base-engine/config"
	"base-engine/gen"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestAIModelConfigEncryptedAndRedacted(t *testing.T) {
	db := modelConfigTestDB(t)
	security := config.AIModelSecurityConfig{
		ActiveKeyID: "v1", Keys: map[string][]byte{"v1": bytes.Repeat([]byte{7}, 32)},
	}
	store := NewModelConfigStore(db, security, nil)
	principal := modelConfigPrincipal("aiModelConfig:read", "aiModelConfig:manage")
	status, err := store.Save(t.Context(), principal, ModelConfigUpdate{
		Name: "example", BaseURL: "https://model.example/v1", APIKey: "private-key",
	})
	if err != nil {
		t.Fatal(err)
	}
	if status.Status != "DRAFT" || !status.KeyConfigured || status.Version != 1 {
		t.Fatalf("saved status = %+v", status)
	}
	assertModelConfigEncrypted(t, db)
	read, err := store.Status(t.Context(), principal)
	if err != nil || read.KeyConfigured != true || read.Name != "example" {
		t.Fatalf("read status = %+v, %v", read, err)
	}
	if _, err := store.Active(t.Context()); !errors.Is(err, ErrModelConfigInactive) {
		t.Fatalf("draft config was active: %v", err)
	}
}

func TestAIModelConfigAcceptsExternalSelfHostedHTTPWithoutOriginPolicy(t *testing.T) {
	store := NewModelConfigStore(modelConfigTestDB(t), config.AIModelSecurityConfig{
		ActiveKeyID: "v1", Keys: map[string][]byte{"v1": bytes.Repeat([]byte{7}, 32)},
	}, nil)
	status, err := store.Save(t.Context(), modelConfigPrincipal("aiModelConfig:manage"), ModelConfigUpdate{
		Name: "metis-coder-max", BaseURL: "http://zsgw.sjdistributor.com:4000/v1", APIKey: "private-key",
	})
	if err != nil || status.Status != "DRAFT" {
		t.Fatalf("self-hosted HTTP config = %+v, %v", status, err)
	}
}

func assertModelConfigEncrypted(t *testing.T, db *gorm.DB) {
	t.Helper()
	var row StoredModelConfig
	if err := db.First(&row).Error; err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(row.KeyCiphertext), "private-key") || len(row.KeyCiphertext) == 0 {
		t.Fatalf("secret was not encrypted at rest")
	}
}

func TestAIModelConfigRejectsWrongWorkspaceAndInvalidURL(t *testing.T) {
	db := modelConfigTestDB(t)
	security := config.AIModelSecurityConfig{
		ActiveKeyID: "v1", Keys: map[string][]byte{"v1": bytes.Repeat([]byte{7}, 32)},
	}
	store := NewModelConfigStore(db, security, nil)
	principal := modelConfigPrincipal("aiModelConfig:manage")
	principal.WorkspaceType = auth.WorkspaceTypeFranchise
	_, err := store.Save(t.Context(), principal, ModelConfigUpdate{Name: "example", BaseURL: "https://model.example/v1", APIKey: "key"})
	if !errors.Is(err, ErrModelConfigForbidden) {
		t.Fatalf("franchise error = %v", err)
	}
	principal.WorkspaceType = auth.WorkspaceTypeHeadquarters
	_, err = store.Save(t.Context(), principal, ModelConfigUpdate{Name: "example", BaseURL: "file:///private/etc/passwd", APIKey: "key"})
	if !errors.Is(err, ErrModelConfigInvalid) {
		t.Fatalf("URL error = %v", err)
	}
}

func TestAIModelConfigRequiresEncryptionKey(t *testing.T) {
	store := NewModelConfigStore(modelConfigTestDB(t), config.AIModelSecurityConfig{}, nil)
	_, err := store.Save(t.Context(), modelConfigPrincipal("aiModelConfig:manage"), ModelConfigUpdate{
		Name: "example", BaseURL: "https://model.example/v1", APIKey: "private-key",
	})
	if !errors.Is(err, ErrModelConfigUnavailable) {
		t.Fatalf("missing encryption key error = %v", err)
	}
}

func modelConfigTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{IgnoreRelationshipsWhenMigrating: true})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&StoredModelConfig{}, &gen.AuditLog{}); err != nil {
		t.Fatal(err)
	}
	return db
}

func modelConfigPrincipal(actions ...string) *auth.WorkspacePrincipal {
	permissions := make(map[string]struct{}, len(actions))
	for _, action := range actions {
		permissions[action] = struct{}{}
	}
	return &auth.WorkspacePrincipal{
		AccountID: "hq-user", SessionID: "session-1",
		WorkspaceType: auth.WorkspaceTypeHeadquarters, Permissions: permissions,
	}
}
