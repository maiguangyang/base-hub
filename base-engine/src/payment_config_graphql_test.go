package src_test

import (
	"net/http"
	"strings"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
)

func TestPaymentConfigGeneratedGraphQLIsolation(t *testing.T) {
	db := openResolverTestDB(t)
	cfg := resolverTestConfig()
	seedResolverPermission(t, db, "payment-org-read", "organization:read", gen.PermissionScopeSystem)
	seedResolverPermission(t, db, "payment-config-read", "paymentConfig:read", gen.PermissionScopeSystem)
	hq := seedResolverWorkspace(t, db, cfg, "hq", gen.OrganizationTypeHeadquarters, gen.RoleKindHqSuperAdmin)
	franchise := seedResolverWorkspace(t, db, cfg, "franchise", gen.OrganizationTypeFranchise, gen.RoleKindFranchiseOwner)
	secret := "payment-secret-must-not-appear"
	merchant := "1234567890"
	rows := []any{
		&gen.GlobalPaymentConfig{ID: "global-payment", Channel: "WECHAT", ConfigState: "VALID", Version: 1, CredentialCiphertext: &secret},
		&gen.FranchisePaymentConfig{ID: "franchise-payment", OrganizationID: "franchise-organization", Channel: "WECHAT", ConfigState: "VALID", Version: 1, CredentialCiphertext: &secret, MerchantID: &merchant},
		&gen.Store{ID: "franchise-store", Code: "TEST", Name: "Test", OrganizationID: "franchise-organization", Lifecycle: gen.StoreLifecycleActive},
		&gen.StorePaymentConfig{ID: "store-payment", StoreID: "franchise-store", Channel: "WECHAT", ConfigState: "VALID", Version: 1, CredentialCiphertext: &secret, MerchantID: &merchant},
	}
	for _, row := range rows {
		if err := db.Create(row).Error; err != nil {
			t.Fatal(err)
		}
	}
	handler := governanceHandler(db, cfg)
	queries := []string{
		`{ globalPaymentConfig(id:"global-payment") { credentialCiphertext keyId merchantId } }`,
		`{ franchisePaymentConfig(id:"franchise-payment") { credentialCiphertext } }`,
		`{ storePaymentConfig(id:"store-payment") { credentialCiphertext } }`,
		`{ organization(id:"franchise-organization") { paymentConfigs { credentialCiphertext } } }`,
		`{ organization(id:"franchise-organization") { paymentConfigsIds } }`,
		`{ store(id:"franchise-store") { paymentConfigsIds } }`,
		`{ organizations(filter:{paymentConfigs:{merchantId:"1234567890"}}) { total } }`,
		`{ stores(filter:{paymentConfigs:{merchantId:"1234567890"}}) { total } }`,
		`{ organizations(sort:[{paymentConfigs:{merchantId:ASC}}]) { total } }`,
	}
	for _, cookie := range []struct {
		name  string
		value *http.Cookie
	}{{"hq", hq}, {"franchise", franchise}} {
		for _, query := range queries {
			response := executeResolverGraphQL(t, handler, cookie.value, query)
			if strings.Contains(response.Body.String(), secret) || strings.Contains(response.Body.String(), "franchise-payment") || strings.Contains(response.Body.String(), "store-payment") {
				t.Fatalf("%s leaked payment config: %s", cookie.name, response.Body.String())
			}
			assertGraphQLErrorCode(t, response, auth.CodePermissionDenied)
		}
	}
}
