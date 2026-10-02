package paymentconfig

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"base-engine/auth"
	"base-engine/config"
	"base-engine/gen"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestPaymentConfigDisabledInheritanceAndABA(t *testing.T) {
	db := paymentTestDB(t)
	keys := config.PaymentConfigSecurity{ActiveKeyID: "one", Keys: map[string][]byte{"one": []byte("12345678901234567890123456789012")}}
	store := NewStore(db, keys)
	principal := &auth.WorkspacePrincipal{AccountID: "hq", WorkspaceType: auth.WorkspaceTypeHeadquarters, Permissions: map[string]struct{}{"paymentConfig:read": {}, "paymentConfig:manage": {}}}
	ctx := context.Background()
	global := ScopeRef{Scope: "GLOBAL"}
	franchise := ScopeRef{Scope: "FRANCHISE", OrganizationID: "franchise"}
	shop := ScopeRef{Scope: "STORE", StoreID: "shop"}
	assertPaymentDisabledInheritance(t, store, principal, ctx, global, franchise, shop)
	assertPaymentABA(t, store, principal, ctx, shop)
}

func TestPaymentConfigScopeAndRateValidation(t *testing.T) {
	db := paymentTestDB(t)
	service := NewStore(db, config.PaymentConfigSecurity{ActiveKeyID: "one", Keys: map[string][]byte{"one": []byte("12345678901234567890123456789012")}})
	principal := paymentPrincipal()
	ctx := context.Background()
	if err := db.Create(&gen.Organization{ID: "headquarters", Code: "HQ", Name: "HQ", Type: gen.OrganizationTypeHeadquarters, Status: gen.OrganizationStatusActive}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&gen.Store{ID: "direct", Code: "DIRECT", Name: "Direct", OrganizationID: "headquarters", Lifecycle: gen.StoreLifecycleActive}).Error; err != nil {
		t.Fatal(err)
	}
	for _, ref := range []ScopeRef{{Scope: "STORE", StoreID: "direct"}, {Scope: "STORE", StoreID: "shop", OrganizationID: "headquarters"}, {Scope: "FRANCHISE", OrganizationID: "headquarters"}, {Scope: "STORE", StoreID: "missing"}} {
		if _, err := service.Read(ctx, principal, ref); !errors.Is(err, ErrInvalid) {
			t.Fatalf("invalid scope %#v: %v", ref, err)
		}
	}
	if _, err := service.Read(ctx, principal, ScopeRef{Scope: "GLOBAL"}); err != nil {
		t.Fatalf("global configuration unavailable for direct stores: %v", err)
	}
	for _, rate := range []int{-1, 1000001} {
		if _, err := service.Save(ctx, principal, SaveInput{ScopeRef: ScopeRef{Scope: "GLOBAL"}, Channel: "WECHAT", RatePpm: rate}); !errors.Is(err, ErrInvalid) {
			t.Fatalf("invalid rate %d: %v", rate, err)
		}
	}
}

func TestPaymentConfigReadRequiresActiveStore(t *testing.T) {
	db := paymentTestDB(t)
	service := NewStore(db, config.PaymentConfigSecurity{ActiveKeyID: "one", Keys: map[string][]byte{"one": []byte("12345678901234567890123456789012")}})
	for _, lifecycle := range []gen.StoreLifecycle{gen.StoreLifecycleDraft, gen.StoreLifecyclePendingApproval, gen.StoreLifecycleRejected, gen.StoreLifecycleActive} {
		if err := db.Model(&gen.Store{}).Where("id = ?", "shop").Update("lifecycle", lifecycle).Error; err != nil {
			t.Fatal(err)
		}
		_, err := service.Read(context.Background(), paymentPrincipal(), ScopeRef{Scope: "STORE", StoreID: "shop"})
		if lifecycle == gen.StoreLifecycleActive {
			if err != nil {
				t.Fatalf("active store read failed: %v", err)
			}
		} else if !errors.Is(err, ErrInvalid) {
			t.Fatalf("%s store error = %v, want invalid scope", lifecycle, err)
		}
	}
}

func TestPaymentConfigFeeUpdateKeepsDisabledState(t *testing.T) {
	db := paymentTestDB(t)
	service := NewStore(db, config.PaymentConfigSecurity{ActiveKeyID: "one", Keys: map[string][]byte{"one": []byte("12345678901234567890123456789012")}})
	principal := paymentPrincipal()
	ref := ScopeRef{Scope: "GLOBAL"}
	view, err := service.Save(context.Background(), principal, SaveInput{ScopeRef: ref, Channel: "WECHAT", MerchantID: "1234567890", RatePpm: 3800, WechatCredentials: testWechatCredentials(t)})
	if err != nil {
		t.Fatal(err)
	}
	active := view.Channels[0].Own
	view, err = service.SetState(context.Background(), principal, StateInput{ScopeRef: ref, Channel: "WECHAT", RecordID: active.RecordID, Version: active.Version, State: "DISABLED"})
	if err != nil {
		t.Fatal(err)
	}
	disabled := view.Channels[0].Own
	view, err = service.Save(context.Background(), principal, SaveInput{ScopeRef: ref, Channel: "WECHAT", RecordID: disabled.RecordID, Version: disabled.Version, MerchantID: "1234567890", RatePpm: 4100})
	if err != nil {
		t.Fatal(err)
	}
	got := view.Channels[0].Own
	if got.State != "DISABLED" || got.RatePpm != 4100 || got.Version != disabled.Version+1 {
		t.Fatalf("fee update changed disabled state: %#v", got)
	}
}

func TestPaymentConfigManageWithoutReadDoesNotCommit(t *testing.T) {
	db := paymentTestDB(t)
	service := NewStore(db, config.PaymentConfigSecurity{ActiveKeyID: "one", Keys: map[string][]byte{"one": []byte("12345678901234567890123456789012")}})
	principal := &auth.WorkspacePrincipal{AccountID: "hq", WorkspaceType: auth.WorkspaceTypeHeadquarters, Permissions: map[string]struct{}{"paymentConfig:manage": {}}}
	_, err := service.SetState(context.Background(), principal, StateInput{ScopeRef: ScopeRef{Scope: "GLOBAL"}, Channel: "WECHAT", State: "DISABLED"})
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("manage-only mutation error = %v, want forbidden", err)
	}
	var count int64
	if err := db.Table("global_payment_configs").Count(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("manage-only mutation committed %d records", count)
	}
}

func TestPaymentConfigAuditRecordsOldAndNewVersions(t *testing.T) {
	db := paymentTestDB(t)
	service := NewStore(db, config.PaymentConfigSecurity{ActiveKeyID: "one", Keys: map[string][]byte{"one": []byte("12345678901234567890123456789012")}})
	principal := paymentPrincipal()
	ref := ScopeRef{Scope: "FRANCHISE", OrganizationID: "franchise"}
	view, err := service.SetState(context.Background(), principal, StateInput{ScopeRef: ref, Channel: "WECHAT", State: "DISABLED"})
	if err != nil {
		t.Fatal(err)
	}
	one := view.Channels[0].Own
	view, err = service.SetState(context.Background(), principal, StateInput{ScopeRef: ref, Channel: "WECHAT", RecordID: one.RecordID, Version: one.Version, State: "DISABLED"})
	if err != nil {
		t.Fatal(err)
	}
	two := view.Channels[0].Own
	view, err = service.Save(context.Background(), principal, SaveInput{ScopeRef: ref, Channel: "WECHAT", RecordID: two.RecordID, Version: two.Version,
		MerchantID: "1234567890", RatePpm: 3800, WechatCredentials: testWechatCredentials(t)})
	if err != nil {
		t.Fatal(err)
	}
	three := view.Channels[0].Own
	if _, err := service.RestoreInheritance(context.Background(), principal, ResetInput{ScopeRef: ref, Channel: "WECHAT", RecordID: three.RecordID, Version: three.Version}); err != nil {
		t.Fatal(err)
	}
	assertPaymentAuditTransitions(t, db)
}

func assertPaymentAuditTransitions(t *testing.T, db *gorm.DB) {
	t.Helper()
	var logs []gen.AuditLog
	if err := db.Where("resource_type = ?", "paymentConfig").Find(&logs).Error; err != nil {
		t.Fatal(err)
	}
	want := map[string]int{"paymentConfig:state:0:1": 1, "paymentConfig:state:1:2": 1,
		"paymentConfig:save:2:3": 1, "paymentConfig:restore_inheritance:3:0": 1}
	got := make(map[string]int)
	for _, log := range logs {
		var metadata map[string]any
		if err := json.Unmarshal([]byte(value(log.MetadataJSON)), &metadata); err != nil {
			t.Fatal(err)
		}
		old, oldOK := metadata["paymentConfigOldVersion"].(float64)
		next, nextOK := metadata["paymentConfigNewVersion"].(float64)
		if !oldOK || !nextOK {
			t.Fatalf("%s missing version transition: %#v", log.Action, metadata)
		}
		got[fmt.Sprintf("%s:%d:%d", log.Action, int(old), int(next))]++
	}
	if len(got) != len(want) {
		t.Fatalf("audit transitions = %#v, want %#v", got, want)
	}
	for transition, count := range want {
		if got[transition] != count {
			t.Fatalf("audit transitions = %#v, want %#v", got, want)
		}
	}
}

func assertPaymentDisabledInheritance(t *testing.T, store *Store, principal *auth.WorkspacePrincipal, ctx context.Context, global, franchise, shop ScopeRef) {
	t.Helper()
	if _, err := store.SetState(ctx, principal, StateInput{ScopeRef: global, Channel: "WECHAT", State: "DISABLED"}); err != nil {
		t.Fatal(err)
	}
	view, err := store.Read(ctx, principal, shop)
	if err != nil || view.Channels[0].Effective.SourceScope != "GLOBAL" || view.Channels[0].Effective.State != "DISABLED" {
		t.Fatalf("global fallback: %#v %v", view, err)
	}
	if _, err := store.SetState(ctx, principal, StateInput{ScopeRef: franchise, Channel: "WECHAT", State: "DISABLED"}); err != nil {
		t.Fatal(err)
	}
	view, err = store.Read(ctx, principal, shop)
	if err != nil || view.Channels[0].Effective.SourceScope != "FRANCHISE" {
		t.Fatalf("franchise override: %#v %v", view, err)
	}
	view, err = store.Read(ctx, principal, shop)
	if err != nil || view.Channels[1].Effective.State != "UNCONFIGURED" {
		t.Fatalf("channel isolation: %#v %v", view, err)
	}
}

func assertPaymentABA(t *testing.T, store *Store, principal *auth.WorkspacePrincipal, ctx context.Context, shop ScopeRef) {
	t.Helper()
	view, err := store.SetState(ctx, principal, StateInput{ScopeRef: shop, Channel: "WECHAT", State: "DISABLED"})
	if err != nil {
		t.Fatal(err)
	}
	old := view.Channels[0].Own
	if _, err := store.RestoreInheritance(ctx, principal, ResetInput{ScopeRef: shop, Channel: "WECHAT", RecordID: old.RecordID, Version: old.Version}); err != nil {
		t.Fatal(err)
	}
	view, err = store.SetState(ctx, principal, StateInput{ScopeRef: shop, Channel: "WECHAT", State: "DISABLED"})
	if err != nil {
		t.Fatal(err)
	}
	if old.RecordID == view.Channels[0].Own.RecordID {
		t.Fatal("recreated override reused record ID")
	}
	if _, err := store.RestoreInheritance(ctx, principal, ResetInput{ScopeRef: shop, Channel: "WECHAT", RecordID: old.RecordID, Version: old.Version}); !errors.Is(err, ErrConflict) {
		t.Fatalf("stale restore: %v", err)
	}
	if _, err := store.Save(ctx, principal, SaveInput{ScopeRef: shop, Channel: "WECHAT", RecordID: old.RecordID, Version: old.Version}); !errors.Is(err, ErrConflict) {
		t.Fatalf("stale save: %v", err)
	}
}

func paymentTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true, IgnoreRelationshipsWhenMigrating: true})
	if err != nil {
		t.Fatal(err)
	}
	for _, model := range []any{&gen.Organization{}, &gen.Store{}, &gen.GlobalPaymentConfig{}, &gen.FranchisePaymentConfig{}, &gen.StorePaymentConfig{}, &gen.AuditLog{}} {
		if err := db.AutoMigrate(model); err != nil {
			t.Fatal(err)
		}
		for _, column := range []string{"is_delete", "weight", "state", "deleted_by", "updated_by", "created_by"} {
			db.Exec("DROP INDEX IF EXISTS `" + column + "`")
		}
	}
	if err := db.Create(&gen.Organization{ID: "franchise", Code: "FRAN", Name: "Franchise", Type: gen.OrganizationTypeFranchise, Status: gen.OrganizationStatusActive}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&gen.Store{ID: "shop", Code: "SHOP", Name: "Shop", OrganizationID: "franchise", Lifecycle: gen.StoreLifecycleActive}).Error; err != nil {
		t.Fatal(err)
	}
	return db
}
