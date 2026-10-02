package ai

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"base-engine/auth"
	"base-engine/config"
	"base-engine/gen"
	"gorm.io/gorm"
)

func TestAIModelConfigRequiresVersionedProbeBeforeActivation(t *testing.T) {
	db := modelConfigTestDB(t)
	security := config.AIModelSecurityConfig{ActiveKeyID: "v1", Keys: map[string][]byte{"v1": bytes.Repeat([]byte{7}, 32)}}
	probed := 0
	store := NewModelConfigStore(db, security, func(_ context.Context, cfg ModelConfig) error {
		probed++
		if cfg.APIKey != "secret" || cfg.Name != "example" {
			t.Fatalf("probe config = %+v", cfg)
		}
		return nil
	})
	principal := modelConfigPrincipal("aiModelConfig:manage")
	saved, err := store.Save(t.Context(), principal, ModelConfigUpdate{Name: "example", BaseURL: "https://model.example/v1", APIKey: "secret"})
	if err != nil {
		t.Fatal(err)
	}
	active := assertModelConfigActivationFlow(t, store, principal, saved.Version, &probed)
	assertModelConfigAuditActions(t, db)
	assertModelConfigDeactivation(t, store, principal, active.Version)
	updated, err := store.Save(t.Context(), principal, ModelConfigUpdate{Name: "changed", BaseURL: "https://model.example/v1", ExpectedVersion: active.Version})
	if err != nil || updated.Status != "DRAFT" {
		t.Fatalf("update = %+v, %v", updated, err)
	}
	if _, err := store.Active(t.Context()); !errors.Is(err, ErrModelConfigInactive) {
		t.Fatalf("old config remained active: %v", err)
	}
	if _, err := store.Activate(t.Context(), principal, active.Version); !errors.Is(err, ErrModelConfigConflict) {
		t.Fatalf("stale activation = %v", err)
	}
}

func assertModelConfigDeactivation(t *testing.T, store *ModelConfigStore, principal *auth.WorkspacePrincipal, version uint64) {
	t.Helper()
	deactivated, err := store.Deactivate(t.Context(), principal, version)
	if err != nil || deactivated.Status != "DRAFT" || !errors.Is(store.CurrentVersionActive(t.Context(), version), ErrModelConfigInactive) {
		t.Fatalf("deactivation = %+v, %v", deactivated, err)
	}
}

func assertModelConfigAuditActions(t *testing.T, db *gorm.DB) {
	t.Helper()
	var records []gen.AuditLog
	if err := db.Where("resource_type = ?", "aiModelConfig").Order("rowid").Find(&records).Error; err != nil {
		t.Fatal(err)
	}
	if len(records) != 3 || records[0].Action != "aiModelConfig:save" || records[1].Action != "aiModelConfig:probe" || records[2].Action != "aiModelConfig:activate" {
		t.Fatalf("audit actions = %+v", records)
	}
}

func assertModelConfigActivationFlow(t *testing.T, store *ModelConfigStore, principal *auth.WorkspacePrincipal, version uint64, probed *int) ModelConfigStatus {
	t.Helper()
	if _, err := store.Activate(t.Context(), principal, version); !errors.Is(err, ErrModelConfigConflict) {
		t.Fatalf("activation before probe = %v", err)
	}
	tested, err := store.Probe(t.Context(), principal, version)
	if err != nil || tested.Status != "TESTED" || *probed != 1 {
		t.Fatalf("probe = %+v, %v; calls=%d", tested, err, *probed)
	}
	active, err := store.Activate(t.Context(), principal, version)
	if err != nil || active.Status != "ACTIVE" {
		t.Fatalf("activation = %+v, %v", active, err)
	}
	if _, err := store.Active(t.Context()); err != nil {
		t.Fatalf("active config unavailable: %v", err)
	}
	return active
}

func TestAIModelConfigFailedProbeRemainsDraft(t *testing.T) {
	db := modelConfigTestDB(t)
	security := config.AIModelSecurityConfig{ActiveKeyID: "v1", Keys: map[string][]byte{"v1": bytes.Repeat([]byte{7}, 32)}}
	store := NewModelConfigStore(db, security, func(context.Context, ModelConfig) error { return ErrModelUnavailable })
	principal := modelConfigPrincipal("aiModelConfig:read", "aiModelConfig:manage")
	saved, err := store.Save(t.Context(), principal, ModelConfigUpdate{Name: "example", BaseURL: "https://model.example/v1", APIKey: "secret"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.Probe(t.Context(), principal, saved.Version); !errors.Is(err, ErrModelUnavailable) {
		t.Fatalf("probe error = %v", err)
	}
	status, err := store.Status(t.Context(), principal)
	if err != nil || status.Status != "DRAFT" {
		t.Fatalf("failed probe state = %+v, %v", status, err)
	}
}
