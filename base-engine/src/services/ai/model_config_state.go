package ai

import (
	"context"
	"time"

	"base-engine/auth"
	"gorm.io/gorm"
)

func (s *ModelConfigStore) Probe(ctx context.Context, principal *auth.WorkspacePrincipal, version uint64) (ModelConfigStatus, error) {
	if err := authorizeModelConfig(principal, "aiModelConfig:manage"); err != nil {
		return ModelConfigStatus{}, err
	}
	row, err := s.load(ctx)
	if err != nil || row.Version != version || row.Status != "DRAFT" {
		return ModelConfigStatus{}, ErrModelConfigConflict
	}
	key, err := decryptModelKey(s.security, row.KeyID, row.KeyCiphertext)
	if err != nil {
		return ModelConfigStatus{}, err
	}
	if err := s.probe(ctx, ModelConfig{Name: row.Name, BaseURL: row.BaseURL, APIKey: key}); err != nil {
		if auditErr := s.auditModelConfig(s.db.WithContext(ctx), principal, "aiModelConfig:probe", version, "DRAFT", "MODEL_PROBE_FAILED"); auditErr != nil {
			return ModelConfigStatus{}, auditErr
		}
		return ModelConfigStatus{}, err
	}
	now := time.Now().UTC()
	return s.changeStatus(ctx, principal, version, "DRAFT", "TESTED", &now, "aiModelConfig:probe")
}

func (s *ModelConfigStore) Activate(ctx context.Context, principal *auth.WorkspacePrincipal, version uint64) (ModelConfigStatus, error) {
	if err := authorizeModelConfig(principal, "aiModelConfig:manage"); err != nil {
		return ModelConfigStatus{}, err
	}
	return s.changeStatus(ctx, principal, version, "TESTED", "ACTIVE", nil, "aiModelConfig:activate")
}

func (s *ModelConfigStore) Deactivate(ctx context.Context, principal *auth.WorkspacePrincipal, version uint64) (ModelConfigStatus, error) {
	if err := authorizeModelConfig(principal, "aiModelConfig:manage"); err != nil {
		return ModelConfigStatus{}, err
	}
	return s.changeStatus(ctx, principal, version, "ACTIVE", "DRAFT", nil, "aiModelConfig:deactivate")
}

func (s *ModelConfigStore) changeStatus(ctx context.Context, principal *auth.WorkspacePrincipal, version uint64, from, to string, testedAt *time.Time, action string) (ModelConfigStatus, error) {
	updates := map[string]any{"status": to}
	if testedAt != nil || to == "DRAFT" {
		updates["tested_at"] = testedAt
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&StoredModelConfig{}).
			Where("id = ? AND version = ? AND status = ?", modelConfigRecordID, version, from).Updates(updates)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return ErrModelConfigConflict
		}
		return s.auditModelConfig(tx, principal, action, version, to, "SUCCESS")
	})
	if err != nil {
		return ModelConfigStatus{}, err
	}
	row, err := s.load(ctx)
	if err != nil {
		return ModelConfigStatus{}, err
	}
	return statusFor(row), nil
}

func (s *ModelConfigStore) CurrentVersionActive(ctx context.Context, version uint64) error {
	var row StoredModelConfig
	err := s.db.WithContext(ctx).Select("version", "status").First(&row, modelConfigRecordID).Error
	if err != nil || row.Version != version || row.Status != "ACTIVE" {
		return ErrModelConfigInactive
	}
	return nil
}
