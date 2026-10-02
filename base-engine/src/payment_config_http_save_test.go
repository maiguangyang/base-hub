package src_test

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"base-engine/auth"
	"base-engine/config"
	"base-engine/gen"
	enginesrc "base-engine/src"
	"base-engine/src/services/paymentconfig"
)

func TestPaymentConfigHTTPPutRedactsAndAudits(t *testing.T) {
	db := openResolverTestDB(t)
	router := mux.NewRouter()
	keys := config.PaymentConfigSecurity{ActiveKeyID: "test", Keys: map[string][]byte{"test": bytes.Repeat([]byte{7}, 32)}}
	enginesrc.RegisterPaymentConfigRoutes(router, paymentconfig.NewStore(db, keys), config.SecurityConfig{AllowedOrigins: map[string]struct{}{"https://admin.example": {}}})
	hq := &auth.WorkspacePrincipal{AccountID: "test", WorkspaceType: auth.WorkspaceTypeHeadquarters, Permissions: map[string]struct{}{"paymentConfig:read": {}, "paymentConfig:manage": {}}}
	private, cert := paymentHTTPCertificate(t)
	input := paymentconfig.SaveInput{ScopeRef: paymentconfig.ScopeRef{Scope: "GLOBAL"}, Channel: "WECHAT", MerchantID: "1234567890", RatePpm: 3800,
		WechatCredentials: &paymentconfig.WechatCredentials{MerchantAPISerial: "2A", MerchantAPICert: cert, APIV3Key: strings.Repeat("z", 32), MerchantPrivateKey: private, VerificationMode: "PLATFORM_CERT"}}
	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	response := paymentHTTPCall(router, hq, http.MethodPut, "/api/payment-config", string(body), "https://admin.example")
	if response.Code != 200 {
		t.Fatalf("save: %d %s", response.Code, response.Body.String())
	}
	assertPaymentHTTPSecretsAbsent(t, response.Body.String(), private, input.WechatCredentials.APIV3Key, input.MerchantID)
	var view paymentconfig.ScopeView
	if err := json.Unmarshal(response.Body.Bytes(), &view); err != nil {
		t.Fatal(err)
	}
	input.WechatCredentials, input.MerchantID, input.RecordID, input.Version, input.RatePpm = nil, "", view.Channels[0].Own.RecordID, 1, 4000
	body, err = json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	response = paymentHTTPCall(router, hq, http.MethodPut, "/api/payment-config", string(body), "https://admin.example")
	if response.Code != 200 || !strings.Contains(response.Body.String(), `"ratePpm":4000`) {
		t.Fatalf("fee update: %d %s", response.Code, response.Body.String())
	}
	var auditLogs []gen.AuditLog
	if err := db.Find(&auditLogs).Error; err != nil || len(auditLogs) != 2 {
		t.Fatalf("audit records: %v %d", err, len(auditLogs))
	}
	for _, audit := range auditLogs {
		assertPaymentHTTPSecretsAbsent(t, valueOrEmpty(audit.MetadataJSON), private, strings.Repeat("z", 32), "1234567890")
	}
}

func paymentHTTPCertificate(t *testing.T) (string, string) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	certificate := &x509.Certificate{SerialNumber: big.NewInt(42), Subject: pkix.Name{CommonName: "payment-test"}, NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(time.Hour)}
	der, err := x509.CreateCertificate(rand.Reader, certificate, certificate, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	private := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	cert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	return string(private), string(cert)
}

func assertPaymentHTTPSecretsAbsent(t *testing.T, body string, secrets ...string) {
	t.Helper()
	for _, secret := range secrets {
		if strings.Contains(body, secret) {
			t.Fatal("payment response or audit leaked a secret")
		}
	}
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
