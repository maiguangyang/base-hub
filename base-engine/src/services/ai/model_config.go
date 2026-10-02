package ai

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"

	"base-engine/auth"
	"base-engine/config"
	"base-engine/src/services/audit"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrModelConfigForbidden = errors.New("AI model configuration forbidden")
	ErrModelConfigInvalid   = errors.New("AI model configuration invalid")
	ErrModelConfigInactive  = errors.New("AI model configuration inactive")
	ErrModelConfigConflict  = errors.New("AI model configuration conflict")
)

const modelConfigRecordID = 1

type StoredModelConfig struct {
	ID            int    `gorm:"primaryKey"`
	Name          string `gorm:"type:varchar(128);not null"`
	BaseURL       string `gorm:"type:varchar(2048);not null"`
	KeyCiphertext []byte `gorm:"not null"`
	KeyID         string `gorm:"type:varchar(64);not null"`
	Version       uint64 `gorm:"not null"`
	Status        string `gorm:"type:varchar(16);not null"`
	TestedAt      *time.Time
	UpdatedAt     time.Time
}

type ModelConfigStatus struct {
	Name          string     `json:"modelName"`
	BaseURL       string     `json:"baseUrl"`
	KeyConfigured bool       `json:"keyConfigured"`
	Version       uint64     `json:"version"`
	Status        string     `json:"status"`
	TestedAt      *time.Time `json:"testedAt,omitempty"`
}

type ModelConfigUpdate struct {
	Name            string `json:"modelName"`
	BaseURL         string `json:"baseUrl"`
	APIKey          string `json:"apiKey"`
	ExpectedVersion uint64 `json:"version"`
}

type ModelConfigStore struct {
	db       *gorm.DB
	security config.AIModelSecurityConfig
	probe    func(context.Context, ModelConfig) error
	audit    *audit.Service
}

func NewModelConfigStore(db *gorm.DB, security config.AIModelSecurityConfig, probe func(context.Context, ModelConfig) error) *ModelConfigStore {
	if probe == nil {
		probe = ProbeModelConnection
	}
	return &ModelConfigStore{db: db, security: security, probe: probe, audit: audit.NewService()}
}

func (s *ModelConfigStore) Status(ctx context.Context, principal *auth.WorkspacePrincipal) (ModelConfigStatus, error) {
	if err := authorizeModelConfig(principal, "aiModelConfig:read"); err != nil {
		return ModelConfigStatus{}, err
	}
	row, err := s.load(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ModelConfigStatus{Status: "UNCONFIGURED"}, nil
	}
	if err != nil {
		return ModelConfigStatus{}, err
	}
	return statusFor(row), nil
}

func (s *ModelConfigStore) Save(ctx context.Context, principal *auth.WorkspacePrincipal, input ModelConfigUpdate) (ModelConfigStatus, error) {
	if err := authorizeModelConfig(principal, "aiModelConfig:manage"); err != nil {
		return ModelConfigStatus{}, err
	}
	if err := s.validateInput(input); err != nil {
		return ModelConfigStatus{}, err
	}
	var saved StoredModelConfig
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		saved, err = s.saveTx(tx, input)
		if err != nil {
			return err
		}
		action := "aiModelConfig:save"
		if input.ExpectedVersion > 0 && input.APIKey != "" {
			action = "aiModelConfig:rotate_key"
		}
		return s.auditModelConfig(tx, principal, action, saved.Version, saved.Status, "SUCCESS")
	})
	if err != nil {
		return ModelConfigStatus{}, err
	}
	return statusFor(saved), nil
}

func (s *ModelConfigStore) saveTx(tx *gorm.DB, input ModelConfigUpdate) (StoredModelConfig, error) {
	row := StoredModelConfig{ID: modelConfigRecordID}
	read := tx
	if input.ExpectedVersion > 0 {
		read = read.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	err := read.First(&row, modelConfigRecordID).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return row, err
	}
	if row.Version != input.ExpectedVersion {
		return row, ErrModelConfigConflict
	}
	key, err := s.keyForUpdate(row, input.APIKey)
	if err != nil {
		return row, err
	}
	if unchangedModelConfigSave(row, input, s.security.ActiveKeyID) {
		return row, nil
	}
	ciphertext, err := encryptModelKey(s.security, key)
	if err != nil {
		return row, err
	}
	row.Name, row.BaseURL = input.Name, input.BaseURL
	row.KeyCiphertext, row.KeyID = ciphertext, s.security.ActiveKeyID
	row.Version, row.Status, row.TestedAt = row.Version+1, "DRAFT", nil
	if err := saveModelConfigRow(tx, row, input.ExpectedVersion); err != nil {
		return row, err
	}
	return row, nil
}

func unchangedModelConfigSave(row StoredModelConfig, input ModelConfigUpdate, activeKeyID string) bool {
	return row.Version > 0 && input.APIKey == "" && row.Name == input.Name &&
		row.BaseURL == input.BaseURL && row.KeyID == activeKeyID
}

func (s *ModelConfigStore) keyForUpdate(row StoredModelConfig, supplied string) (string, error) {
	if supplied != "" {
		return supplied, nil
	}
	if row.Version == 0 {
		return "", ErrModelConfigInvalid
	}
	return decryptModelKey(s.security, row.KeyID, row.KeyCiphertext)
}

func saveModelConfigRow(tx *gorm.DB, row StoredModelConfig, previous uint64) error {
	if previous == 0 {
		result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&row)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return ErrModelConfigConflict
		}
		return nil
	}
	result := tx.Model(&StoredModelConfig{}).Where("id = ? AND version = ?", modelConfigRecordID, previous).
		Updates(map[string]any{
			"name": row.Name, "base_url": row.BaseURL, "key_ciphertext": row.KeyCiphertext,
			"key_id": row.KeyID, "version": row.Version, "status": row.Status, "tested_at": nil,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return ErrModelConfigConflict
	}
	return nil
}

func (s *ModelConfigStore) Active(ctx context.Context) (ModelConfig, error) {
	model, _, err := s.ActiveWithVersion(ctx)
	return model, err
}

func (s *ModelConfigStore) ActiveWithVersion(ctx context.Context) (ModelConfig, uint64, error) {
	row, err := s.load(ctx)
	if err != nil || row.Status != "ACTIVE" {
		return ModelConfig{}, 0, ErrModelConfigInactive
	}
	key, err := decryptModelKey(s.security, row.KeyID, row.KeyCiphertext)
	if err != nil {
		return ModelConfig{}, 0, err
	}
	return ModelConfig{Name: row.Name, BaseURL: row.BaseURL, APIKey: key}, row.Version, nil
}

func (s *ModelConfigStore) load(ctx context.Context) (StoredModelConfig, error) {
	var row StoredModelConfig
	err := s.db.WithContext(ctx).First(&row, modelConfigRecordID).Error
	return row, err
}

func (s *ModelConfigStore) validateInput(input ModelConfigUpdate) error {
	if strings.TrimSpace(input.Name) == "" || len(input.Name) > 128 || len(input.BaseURL) > 2048 {
		return ErrModelConfigInvalid
	}
	return s.validateURL(input.BaseURL)
}

func (s *ModelConfigStore) validateURL(raw string) error {
	parsed, err := url.Parse(raw)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" || parsed.User != nil ||
		parsed.RawQuery != "" || parsed.Fragment != "" {
		return ErrModelConfigInvalid
	}
	return nil
}

func authorizeModelConfig(principal *auth.WorkspacePrincipal, action string) error {
	if principal == nil || principal.WorkspaceType != auth.WorkspaceTypeHeadquarters || !principal.Has(action) {
		return ErrModelConfigForbidden
	}
	return nil
}

func statusFor(row StoredModelConfig) ModelConfigStatus {
	return ModelConfigStatus{
		Name: row.Name, BaseURL: row.BaseURL, KeyConfigured: len(row.KeyCiphertext) > 0,
		Version: row.Version, Status: row.Status, TestedAt: row.TestedAt,
	}
}
