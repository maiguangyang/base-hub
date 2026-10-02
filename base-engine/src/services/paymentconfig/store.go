package paymentconfig

import (
	"context"
	"strings"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/config"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"gorm.io/gorm"
)

type Store struct {
	db    *gorm.DB
	keys  config.PaymentConfigSecurity
	audit *audit.Service
}

type configRow struct {
	ID                   string  `gorm:"column:id;primaryKey"`
	Channel              string  `gorm:"column:channel"`
	MerchantID           *string `gorm:"column:merchant_id"`
	Environment          *string `gorm:"column:environment"`
	RatePpm              int     `gorm:"column:rate_ppm"`
	ConfigState          string  `gorm:"column:config_state"`
	Version              uint64  `gorm:"column:version"`
	KeyID                *string `gorm:"column:key_id"`
	CredentialCiphertext *string `gorm:"column:credential_ciphertext"`
	OrganizationID       string  `gorm:"column:organization_id"`
	StoreID              string  `gorm:"column:store_id"`
	IsDelete             *int64  `gorm:"column:is_delete"`
}

type target struct {
	scope, table, column, id string
	organizationID           string
}

func NewStore(db *gorm.DB, keys config.PaymentConfigSecurity) *Store {
	return &Store{db: db, keys: keys, audit: audit.NewService()}
}

func authorize(principal *auth.WorkspacePrincipal, action string) error {
	if principal == nil || principal.WorkspaceType != auth.WorkspaceTypeHeadquarters || !principal.Has(action) {
		return ErrForbidden
	}
	return nil
}

func (s *Store) ready() error {
	if s.keys.ActiveKeyID == "" || len(s.keys.Keys[s.keys.ActiveKeyID]) != 32 {
		return ErrUnavailable
	}
	return nil
}

func validChannel(channel string) bool { return channel == "WECHAT" || channel == "ALIPAY" }

func (s *Store) scopeTarget(ctx context.Context, ref ScopeRef) (target, error) {
	switch ref.Scope {
	case "GLOBAL":
		if ref.OrganizationID != "" || ref.StoreID != "" {
			return target{}, ErrInvalid
		}
		return target{scope: "GLOBAL", table: "global_payment_configs"}, nil
	case "FRANCHISE":
		return s.franchiseTarget(ctx, ref)
	case "STORE":
		return s.storeTarget(ctx, ref)
	default:
		return target{}, ErrInvalid
	}
}

func (s *Store) franchiseTarget(ctx context.Context, ref ScopeRef) (target, error) {
	if ref.OrganizationID == "" || ref.StoreID != "" {
		return target{}, ErrInvalid
	}
	org, err := s.activeFranchise(ctx, ref.OrganizationID)
	if err != nil {
		return target{}, err
	}
	return target{scope: "FRANCHISE", table: "franchise_payment_configs", column: "organization_id", id: org.ID, organizationID: org.ID}, nil
}

func (s *Store) storeTarget(ctx context.Context, ref ScopeRef) (target, error) {
	if ref.StoreID == "" {
		return target{}, ErrInvalid
	}
	var shop gen.Store
	if err := s.db.WithContext(ctx).First(&shop, "id = ?", ref.StoreID).Error; err != nil {
		return target{}, ErrInvalid
	}
	if shop.Lifecycle != gen.StoreLifecycleActive || !activeRow(shop.IsDelete) {
		return target{}, ErrInvalid
	}
	org, err := s.activeFranchise(ctx, shop.OrganizationID)
	if err != nil || (ref.OrganizationID != "" && ref.OrganizationID != org.ID) {
		return target{}, ErrInvalid
	}
	return target{scope: "STORE", table: "store_payment_configs", column: "store_id", id: shop.ID, organizationID: org.ID}, nil
}

func (s *Store) activeFranchise(ctx context.Context, id string) (gen.Organization, error) {
	var org gen.Organization
	if err := s.db.WithContext(ctx).First(&org, "id = ?", id).Error; err != nil {
		return org, ErrInvalid
	}
	if org.Type != gen.OrganizationTypeFranchise || org.Status != gen.OrganizationStatusActive || !activeRow(org.IsDelete) {
		return org, ErrInvalid
	}
	return org, nil
}

func activeRow(flag *int64) bool { return flag == nil || *flag == 1 }

func (s *Store) load(tx *gorm.DB, t target, channel string) (configRow, bool, error) {
	var row configRow
	query := tx.Table(t.table).Where("channel = ? AND (is_delete IS NULL OR is_delete = 1)", channel)
	if t.column != "" {
		query = query.Where(t.column+" = ?", t.id)
	}
	result := query.Limit(1).Find(&row)
	return row, result.RowsAffected == 1, result.Error
}

func (s *Store) Read(ctx context.Context, principal *auth.WorkspacePrincipal, ref ScopeRef) (ScopeView, error) {
	if err := authorize(principal, "paymentConfig:read"); err != nil {
		return ScopeView{}, err
	}
	if err := s.ready(); err != nil {
		return ScopeView{}, err
	}
	t, err := s.scopeTarget(ctx, ref)
	if err != nil {
		return ScopeView{}, err
	}
	view := ScopeView{Channels: make([]ChannelView, 0, 2)}
	for _, channel := range []string{"WECHAT", "ALIPAY"} {
		own, effective, err := s.resolveChannel(ctx, t, channel)
		if err != nil {
			return ScopeView{}, err
		}
		view.Channels = append(view.Channels, ChannelView{Own: own, Effective: effective})
	}
	return view, nil
}

func (s *Store) resolveChannel(ctx context.Context, t target, channel string) (ConfigStatus, ConfigStatus, error) {
	chain := []target{t}
	if t.scope == "STORE" {
		chain = append(chain, target{scope: "FRANCHISE", table: "franchise_payment_configs", column: "organization_id", id: t.organizationID})
	}
	if t.scope != "GLOBAL" {
		chain = append(chain, target{scope: "GLOBAL", table: "global_payment_configs"})
	}
	own := ConfigStatus{Scope: t.scope, Channel: channel, State: "UNCONFIGURED"}
	for i, current := range chain {
		row, found, err := s.load(s.db.WithContext(ctx), current, channel)
		if err != nil {
			return own, ConfigStatus{}, err
		}
		if !found {
			continue
		}
		status := s.status(current, row)
		if i == 0 {
			own = status
		}
		return own, status, nil
	}
	return own, ConfigStatus{Scope: t.scope, Channel: channel, State: "UNCONFIGURED"}, nil
}

func (s *Store) status(t target, row configRow) ConfigStatus {
	status := ConfigStatus{Scope: t.scope, Channel: row.Channel, SourceScope: t.scope, SourceID: t.id, RecordID: row.ID,
		MerchantMasked: maskMerchant(row.MerchantID), RatePpm: row.RatePpm, State: row.ConfigState,
		Version: row.Version, CredentialsConfigured: row.CredentialCiphertext != nil && *row.CredentialCiphertext != ""}
	if reason := s.credentialError(row); reason != "" {
		status.State, status.ReasonCode = "ERROR", reason
	}
	return status
}

func (s *Store) credentialError(row configRow) string {
	if row.CredentialCiphertext == nil || *row.CredentialCiphertext == "" {
		if row.ConfigState == "VALID" {
			return "CREDENTIALS_MISSING"
		}
		return ""
	}
	if row.KeyID == nil {
		return "DECRYPTION_FAILED"
	}
	data, err := decryptCredentials(s.keys, *row.KeyID, *row.CredentialCiphertext)
	if err != nil {
		return "DECRYPTION_FAILED"
	}
	if !validMerchant(row.Channel, value(row.MerchantID), value(row.Environment)) || validateStoredCredential(row.Channel, data) != nil {
		return "CREDENTIALS_INVALID"
	}
	return ""
}

func maskMerchant(value *string) string {
	if value == nil || *value == "" {
		return ""
	}
	if len(*value) <= 4 {
		return strings.Repeat("*", len(*value))
	}
	return strings.Repeat("*", len(*value)-4) + (*value)[len(*value)-4:]
}

func newID() string { return uuid.Must(uuid.NewV4()).String() }
