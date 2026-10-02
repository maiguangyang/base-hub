package ai

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"base-engine/auth"
	"base-engine/config"
	"base-engine/gen"
	"gorm.io/gorm"
)

func TestAIModelConfigRetainsKeyAndRejectsStaleUpdate(t *testing.T) {
	db := modelConfigTestDB(t)
	security := config.AIModelSecurityConfig{
		ActiveKeyID: "v1", Keys: map[string][]byte{"v1": bytes.Repeat([]byte{7}, 32)},
	}
	store := NewModelConfigStore(db, security, nil)
	principal := modelConfigPrincipal("aiModelConfig:manage")
	first, err := store.Save(t.Context(), principal, ModelConfigUpdate{Name: "old", BaseURL: "https://model.example/v1", APIKey: "secret"})
	if err != nil {
		t.Fatal(err)
	}
	second, err := store.Save(t.Context(), principal, ModelConfigUpdate{Name: "new", BaseURL: "https://model.example/v1", ExpectedVersion: first.Version})
	if err != nil || second.Version != 2 {
		t.Fatalf("update = %+v, %v", second, err)
	}
	var row StoredModelConfig
	if err := db.First(&row).Error; err != nil {
		t.Fatal(err)
	}
	key, err := decryptModelKey(security, row.KeyID, row.KeyCiphertext)
	if err != nil || key != "secret" {
		t.Fatalf("retained key invalid: %v", err)
	}
	_, err = store.Save(t.Context(), principal, ModelConfigUpdate{Name: "stale", BaseURL: "https://model.example/v1", ExpectedVersion: first.Version})
	if !errors.Is(err, ErrModelConfigConflict) {
		t.Fatalf("stale update = %v", err)
	}
	assertModelConfigKeyRotationAudit(t, db, store, principal, second.Version)
}

func assertModelConfigKeyRotationAudit(t *testing.T, db *gorm.DB, store *ModelConfigStore, principal *auth.WorkspacePrincipal, version uint64) {
	t.Helper()
	if _, err := store.Save(t.Context(), principal, ModelConfigUpdate{Name: "new", BaseURL: "https://model.example/v1", APIKey: "replacement", ExpectedVersion: version}); err != nil {
		t.Fatal(err)
	}
	var audit gen.AuditLog
	if err := db.Where("action = ?", "aiModelConfig:rotate_key").First(&audit).Error; err != nil {
		t.Fatal(err)
	}
	if audit.MetadataJSON == nil || strings.Contains(*audit.MetadataJSON, "replacement") {
		t.Fatalf("unsafe key rotation audit = %#v", audit.MetadataJSON)
	}
}

func TestAIModelConfigRewrapsOnMasterKeyRotation(t *testing.T) {
	db := modelConfigTestDB(t)
	v1 := bytes.Repeat([]byte{7}, 32)
	v2 := bytes.Repeat([]byte{8}, 32)
	security := config.AIModelSecurityConfig{ActiveKeyID: "v1", Keys: map[string][]byte{"v1": v1}}
	principal := modelConfigPrincipal("aiModelConfig:manage")
	first, err := NewModelConfigStore(db, security, nil).Save(t.Context(), principal, ModelConfigUpdate{Name: "old", BaseURL: "https://model.example/v1", APIKey: "secret"})
	if err != nil {
		t.Fatal(err)
	}
	security.ActiveKeyID = "v2"
	security.Keys["v2"] = v2
	rotated, err := NewModelConfigStore(db, security, nil).Save(t.Context(), principal, ModelConfigUpdate{Name: "old", BaseURL: "https://model.example/v1", ExpectedVersion: first.Version})
	if err != nil {
		t.Fatal(err)
	}
	if rotated.Version != first.Version+1 || rotated.Status != "DRAFT" {
		t.Fatalf("rotated status = %+v", rotated)
	}
	var row StoredModelConfig
	if err := db.First(&row).Error; err != nil {
		t.Fatal(err)
	}
	if row.KeyID != "v2" {
		t.Fatalf("key ID = %s", row.KeyID)
	}
	delete(security.Keys, "v1")
	key, err := decryptModelKey(security, row.KeyID, row.KeyCiphertext)
	if err != nil || key != "secret" {
		t.Fatalf("rotated key invalid: %v", err)
	}
}

func TestAIModelConfigUnchangedSavePreservesActiveStatus(t *testing.T) {
	db := modelConfigTestDB(t)
	security := config.AIModelSecurityConfig{ActiveKeyID: "v1", Keys: map[string][]byte{"v1": bytes.Repeat([]byte{7}, 32)}}
	store := NewModelConfigStore(db, security, func(context.Context, ModelConfig) error { return nil })
	principal := modelConfigPrincipal("aiModelConfig:manage")
	saved, err := store.Save(t.Context(), principal, ModelConfigUpdate{Name: "example", BaseURL: "https://model.example/v1", APIKey: "secret"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.Probe(t.Context(), principal, saved.Version); err != nil {
		t.Fatal(err)
	}
	active, err := store.Activate(t.Context(), principal, saved.Version)
	if err != nil {
		t.Fatal(err)
	}
	unchanged, err := store.Save(t.Context(), principal, ModelConfigUpdate{Name: "example", BaseURL: "https://model.example/v1", ExpectedVersion: active.Version})
	if err != nil || unchanged.Status != "ACTIVE" || unchanged.Version != active.Version {
		t.Fatalf("unchanged save = %+v, %v", unchanged, err)
	}
}
