package paymentconfig

import (
	"context"
	"base-engine/auth"
	"base-engine/src/services/audit"
	"gorm.io/gorm"
)

func (s *Store) Save(ctx context.Context, principal *auth.WorkspacePrincipal, input SaveInput) (ScopeView, error) {
	t, err := s.mutationTarget(ctx, principal, input.ScopeRef, input.Channel)
	if err != nil {
		return ScopeView{}, err
	}
	if input.RatePpm < 0 || input.RatePpm > 1000000 {
		return ScopeView{}, ErrInvalid
	}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error { return s.saveTx(tx, principal, t, input) })
	if err != nil {
		return ScopeView{}, paymentDBError(err)
	}
	return s.Read(ctx, principal, input.ScopeRef)
}

func (s *Store) saveTx(tx *gorm.DB, principal *auth.WorkspacePrincipal, t target, input SaveInput) error {
	current, found, err := s.load(tx, t, input.Channel)
	if err != nil {
		return err
	}
	if err := checkVersion(current, found, input.RecordID, input.Version); err != nil {
		return err
	}
	row, err := s.prepareSave(current, input)
	if err != nil {
		return err
	}
	if !found {
		row.ID, row.Channel, row.Version = newID(), input.Channel, 1
	}
	if err := writeConfigRow(tx, t, row, found, input.RecordID, input.Version); err != nil {
		return err
	}
	return s.writeAudit(tx, principal, t, row, "paymentConfig:save", current.Version, row.Version)
}

func (s *Store) prepareSave(current configRow, input SaveInput) (configRow, error) {
	row := current
	if input.MerchantID == "" {
		input.MerchantID = value(current.MerchantID)
	}
	if input.Environment == "" {
		input.Environment = value(current.Environment)
	}
	if !validMerchant(input.Channel, input.MerchantID, input.Environment) {
		return row, ErrInvalid
	}
	keyID, ciphertext, err := s.saveCredential(current, input)
	if err != nil {
		return row, err
	}
	row.KeyID, row.CredentialCiphertext = keyID, ciphertext
	row.MerchantID, row.Environment = &input.MerchantID, &input.Environment
	row.RatePpm, row.Version = input.RatePpm, current.Version+1
	if current.Version == 0 {
		row.ConfigState = "VALID"
	}
	return row, nil
}

func (s *Store) saveCredential(current configRow, input SaveInput) (*string, *string, error) {
	if input.WechatCredentials != nil || input.AlipayCredentials != nil {
		data, err := validateAndEncode(input)
		if err != nil {
			return nil, nil, err
		}
		keyID, ciphertext, err := encryptCredentials(s.keys, data)
		if err != nil {
			return nil, nil, err
		}
		return &keyID, &ciphertext, nil
	}
	if current.CredentialCiphertext == nil || current.KeyID == nil || input.MerchantID != value(current.MerchantID) || input.Environment != value(current.Environment) {
		return nil, nil, ErrInvalid
	}
	if err := s.validateRetainedCredential(current); err != nil {
		return nil, nil, err
	}
	return current.KeyID, current.CredentialCiphertext, nil
}

func (s *Store) validateRetainedCredential(row configRow) error {
	data, err := decryptCredentials(s.keys, *row.KeyID, *row.CredentialCiphertext)
	if err != nil {
		return ErrUnavailable
	}
	return validateStoredCredential(row.Channel, data)
}

func (s *Store) SetState(ctx context.Context, principal *auth.WorkspacePrincipal, input StateInput) (ScopeView, error) {
	t, err := s.mutationTarget(ctx, principal, input.ScopeRef, input.Channel)
	if err != nil {
		return ScopeView{}, err
	}
	if input.State != "VALID" && input.State != "DISABLED" {
		return ScopeView{}, ErrInvalid
	}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error { return s.stateTx(tx, principal, t, input) })
	if err != nil {
		return ScopeView{}, paymentDBError(err)
	}
	return s.Read(ctx, principal, input.ScopeRef)
}

func (s *Store) stateTx(tx *gorm.DB, principal *auth.WorkspacePrincipal, t target, input StateInput) error {
	row, found, err := s.load(tx, t, input.Channel)
	if err != nil {
		return err
	}
	if err := checkVersion(row, found, input.RecordID, input.Version); err != nil {
		return err
	}
	oldVersion := row.Version
	if !found {
		if input.State != "DISABLED" {
			return ErrInvalid
		}
		row = configRow{ID: newID(), Channel: input.Channel, Version: 1, ConfigState: "DISABLED"}
	} else {
		if input.State == "VALID" && !s.canEnable(row) {
			return ErrInvalid
		}
		row.ConfigState, row.Version = input.State, row.Version+1
	}
	if err := writeConfigRow(tx, t, row, found, input.RecordID, input.Version); err != nil {
		return err
	}
	return s.writeAudit(tx, principal, t, row, "paymentConfig:state", oldVersion, row.Version)
}

func (s *Store) canEnable(row configRow) bool {
	if row.MerchantID == nil || row.KeyID == nil || row.CredentialCiphertext == nil {
		return false
	}
	data, err := decryptCredentials(s.keys, *row.KeyID, *row.CredentialCiphertext)
	return err == nil && validMerchant(row.Channel, value(row.MerchantID), value(row.Environment)) && validateStoredCredential(row.Channel, data) == nil
}

func (s *Store) RestoreInheritance(ctx context.Context, principal *auth.WorkspacePrincipal, input ResetInput) (ScopeView, error) {
	t, err := s.mutationTarget(ctx, principal, input.ScopeRef, input.Channel)
	if err != nil {
		return ScopeView{}, err
	}
	if t.scope == "GLOBAL" || input.RecordID == "" || input.Version == 0 {
		return ScopeView{}, ErrInvalid
	}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error { return s.restoreTx(tx, principal, t, input) })
	if err != nil {
		return ScopeView{}, paymentDBError(err)
	}
	return s.Read(ctx, principal, input.ScopeRef)
}

func (s *Store) restoreTx(tx *gorm.DB, principal *auth.WorkspacePrincipal, t target, input ResetInput) error {
	row, found, err := s.load(tx, t, input.Channel)
	if err != nil {
		return err
	}
	if err := checkVersion(row, found, input.RecordID, input.Version); err != nil {
		return err
	}
	result := tx.Unscoped().Table(t.table).Where("id = ? AND version = ? AND channel = ?", input.RecordID, input.Version, input.Channel).Delete(&configRow{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return ErrConflict
	}
	return s.writeAudit(tx, principal, t, row, "paymentConfig:restore_inheritance", row.Version, 0)
}

func (s *Store) mutationTarget(ctx context.Context, principal *auth.WorkspacePrincipal, ref ScopeRef, channel string) (target, error) {
	if err := authorize(principal, "paymentConfig:manage"); err != nil {
		return target{}, err
	}
	if err := authorize(principal, "paymentConfig:read"); err != nil {
		return target{}, err
	}
	if err := s.ready(); err != nil {
		return target{}, err
	}
	if !validChannel(channel) {
		return target{}, ErrInvalid
	}
	return s.scopeTarget(ctx, ref)
}

func checkVersion(row configRow, found bool, id string, version uint64) error {
	if !found {
		if id != "" || version != 0 {
			return ErrConflict
		}
		return nil
	}
	if id != row.ID || version != row.Version {
		return ErrConflict
	}
	return nil
}

func writeConfigRow(tx *gorm.DB, t target, row configRow, update bool, id string, previous uint64) error {
	fields := map[string]any{"channel": row.Channel, "merchant_id": row.MerchantID, "environment": row.Environment,
		"rate_ppm": row.RatePpm, "config_state": row.ConfigState, "version": row.Version,
		"key_id": row.KeyID, "credential_ciphertext": row.CredentialCiphertext}
	if !update {
		fields["id"] = row.ID
		if t.column != "" {
			fields[t.column] = t.id
		}
		return tx.Table(t.table).Create(fields).Error
	}
	result := tx.Table(t.table).Where("id = ? AND version = ? AND channel = ?", id, previous, row.Channel).Updates(fields)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return ErrConflict
	}
	return nil
}

func (s *Store) writeAudit(tx *gorm.DB, p *auth.WorkspacePrincipal, t target, row configRow, action string, oldVersion, newVersion uint64) error {
	var session *string
	if p.SessionID != "" {
		session = &p.SessionID
	}
	var org, shop *string
	if t.organizationID != "" {
		org = &t.organizationID
	}
	if t.scope == "STORE" {
		shop = &t.id
	}
	return s.audit.Write(tx, audit.Entry{ActorAccountID: p.AccountID, SessionID: session, OrganizationID: org, StoreID: shop,
		Action: action, ResourceType: "paymentConfig", ResourceID: row.ID, ResultCode: "SUCCESS",
		Metadata: audit.MetadataForPrincipal(p, audit.Metadata{PaymentConfigOldVersion: &oldVersion, PaymentConfigNewVersion: &newVersion,
			PaymentChannel: row.Channel, PaymentScope: t.scope, TargetStatus: row.ConfigState})})
}
