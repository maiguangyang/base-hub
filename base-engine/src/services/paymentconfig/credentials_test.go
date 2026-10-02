package paymentconfig

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"math/big"
	"strings"
	"testing"
	"time"

	"base-engine/auth"
	"base-engine/config"
	"base-engine/gen"
	"gorm.io/gorm"
)

func TestPaymentConfigCredentialEncryptionAndRotation(t *testing.T) {
	db := paymentTestDB(t)
	oldKey := []byte(strings.Repeat("a", 32))
	newKey := []byte(strings.Repeat("b", 32))
	keys := config.PaymentConfigSecurity{ActiveKeyID: "old", Keys: map[string][]byte{"old": oldKey}}
	principal := paymentPrincipal()
	store := NewStore(db, keys)
	input := SaveInput{ScopeRef: ScopeRef{Scope: "GLOBAL"}, Channel: "WECHAT", MerchantID: "1234567890", RatePpm: 3800, WechatCredentials: testWechatCredentials(t)}
	view, err := store.Save(context.Background(), principal, input)
	if err != nil {
		t.Fatal(err)
	}
	assertEncryptedPaymentRecord(t, db, input, view)
	input.RecordID, input.Version, input.RatePpm, input.WechatCredentials = view.Channels[0].Own.RecordID, 1, 4100, nil
	view, err = store.Save(context.Background(), principal, input)
	if err != nil || view.Channels[0].Own.Version != 2 || view.Channels[0].Own.RatePpm != 4100 {
		t.Fatalf("fee-only update: %#v %v", view, err)
	}
	assertPaymentKeyRotation(t, db, principal, input, oldKey, newKey)
}

func assertEncryptedPaymentRecord(t *testing.T, db *gorm.DB, input SaveInput, view ScopeView) {
	t.Helper()
	row := readPaymentRow(t, db, "global_payment_configs")
	if row.KeyID == nil || *row.KeyID != "old" || strings.Contains(value(row.CredentialCiphertext), "PRIVATE KEY") {
		t.Fatal("credential not encrypted")
	}
	var logs []gen.AuditLog
	if err := db.Find(&logs).Error; err != nil || len(logs) != 1 {
		t.Fatalf("audit: %v %d", err, len(logs))
	}
	if strings.Contains(value(logs[0].MetadataJSON), "PRIVATE KEY") || strings.Contains(value(logs[0].MetadataJSON), input.WechatCredentials.APIV3Key) {
		t.Fatal("audit leaked credential")
	}
	if view.Channels[0].Own.MerchantMasked != "******7890" || view.Channels[0].Own.RatePpm != 3800 {
		t.Fatalf("redacted status: %#v", view)
	}
}

func assertPaymentKeyRotation(t *testing.T, db *gorm.DB, principal *auth.WorkspacePrincipal, input SaveInput, oldKey, newKey []byte) {
	t.Helper()
	rotated := NewStore(db, config.PaymentConfigSecurity{ActiveKeyID: "new", Keys: map[string][]byte{"new": newKey, "old": oldKey}})
	view, err := rotated.Read(context.Background(), principal, ScopeRef{Scope: "GLOBAL"})
	if err != nil || view.Channels[0].Own.State != "VALID" {
		t.Fatalf("old key unavailable after rotation: %#v %v", view, err)
	}
	missingOld := NewStore(db, config.PaymentConfigSecurity{ActiveKeyID: "new", Keys: map[string][]byte{"new": newKey}})
	view, err = missingOld.Read(context.Background(), principal, ScopeRef{Scope: "GLOBAL"})
	if err != nil || view.Channels[0].Own.State != "ERROR" || view.Channels[0].Own.ReasonCode != "DECRYPTION_FAILED" {
		t.Fatalf("missing old key: %#v %v", view, err)
	}
	input.RecordID, input.Version, input.WechatCredentials = view.Channels[0].Own.RecordID, 2, testWechatCredentials(t)
	if _, err := missingOld.Save(context.Background(), principal, input); err != nil {
		t.Fatalf("replacing lost-key credential: %v", err)
	}
	row := readPaymentRow(t, db, "global_payment_configs")
	if value(row.KeyID) != "new" {
		t.Fatalf("active key not used: %s", value(row.KeyID))
	}
}

func TestPaymentConfigAlipayAndValidation(t *testing.T) {
	private, cert := testCertificate(t)
	block, _ := pem.Decode([]byte(cert))
	parsed, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatal(err)
	}
	publicDER, err := x509.MarshalPKIXPublicKey(parsed.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	valid := AlipayCredentials{AppPrivateKey: private, AlipayPublicKey: string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: publicDER}))}
	if err := validateAlipay(valid); err != nil {
		t.Fatal(err)
	}
	privateBlock, _ := pem.Decode([]byte(private))
	valid.AppPrivateKey = base64.StdEncoding.EncodeToString(privateBlock.Bytes)
	valid.AlipayPublicKey = base64.StdEncoding.EncodeToString(publicDER)
	if err := validateAlipay(valid); err != nil {
		t.Fatalf("raw key pair rejected: %v", err)
	}
	store := NewStore(paymentTestDB(t), config.PaymentConfigSecurity{ActiveKeyID: "test", Keys: map[string][]byte{"test": []byte(strings.Repeat("k", 32))}})
	view, err := store.Save(context.Background(), paymentPrincipal(), SaveInput{ScopeRef: ScopeRef{Scope: "GLOBAL"}, Channel: "ALIPAY", MerchantID: "1234567890123456", Environment: "SANDBOX", AlipayCredentials: &valid})
	if err != nil || view.Channels[1].Own.State != "VALID" {
		t.Fatalf("public-key Alipay config not saved: %#v %v", view, err)
	}
	assertLegacyAndInvalidPaymentCredentials(t, private, cert, valid)
}

func assertLegacyAndInvalidPaymentCredentials(t *testing.T, private, cert string, valid AlipayCredentials) {
	t.Helper()
	legacyData, err := json.Marshal(map[string]string{"appPrivateKey": private, "appPublicCert": cert, "alipayRootCert": cert, "alipayPublicCert": cert})
	if err != nil {
		t.Fatal(err)
	}
	if err := validateStoredCredential("ALIPAY", legacyData); err != nil {
		t.Fatalf("legacy certificate config rejected: %v", err)
	}
	valid.AlipayPublicKey = "garbage"
	if err := validateAlipay(valid); err == nil {
		t.Fatal("invalid Alipay public key accepted")
	}
	wechat := testWechatCredentials(t)
	if err := validateWechat(*wechat); err != nil {
		t.Fatal(err)
	}
	wechat.APIV3Key = "short"
	if err := validateWechat(*wechat); err == nil {
		t.Fatal("short API V3 key accepted")
	}
}

func TestPaymentConfigExpiredCertificateStopsEnable(t *testing.T) {
	db := paymentTestDB(t)
	keys := config.PaymentConfigSecurity{ActiveKeyID: "test", Keys: map[string][]byte{"test": []byte(strings.Repeat("k", 32))}}
	service := NewStore(db, keys)
	principal := paymentPrincipal()
	ref := ScopeRef{Scope: "GLOBAL"}
	view, err := service.Save(context.Background(), principal, SaveInput{ScopeRef: ref, Channel: "WECHAT", MerchantID: "1234567890", WechatCredentials: testWechatCredentials(t)})
	if err != nil {
		t.Fatal(err)
	}
	private, cert := testCertificateWithValidity(t, time.Now().Add(-2*time.Hour), time.Now().Add(-time.Hour))
	expired := WechatCredentials{MerchantAPISerial: "2A", MerchantAPICert: cert, APIV3Key: strings.Repeat("k", 32), MerchantPrivateKey: private, VerificationMode: "PLATFORM_CERT"}
	data, err := json.Marshal(expired)
	if err != nil {
		t.Fatal(err)
	}
	keyID, cipher, err := encryptCredentials(keys, data)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Table("global_payment_configs").Where("id = ?", view.Channels[0].Own.RecordID).Updates(map[string]any{"key_id": keyID, "credential_ciphertext": cipher, "config_state": "DISABLED"}).Error; err != nil {
		t.Fatal(err)
	}
	view, err = service.Read(context.Background(), principal, ref)
	if err != nil || view.Channels[0].Own.State != "ERROR" || view.Channels[0].Own.ReasonCode != "CREDENTIALS_INVALID" {
		t.Fatalf("expired certificate status: %#v %v", view, err)
	}
	_, err = service.SetState(context.Background(), principal, StateInput{ScopeRef: ref, Channel: "WECHAT", RecordID: view.Channels[0].Own.RecordID, Version: view.Channels[0].Own.Version, State: "VALID"})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expired certificate enabled: %v", err)
	}
}

func paymentPrincipal() *auth.WorkspacePrincipal {
	return &auth.WorkspacePrincipal{AccountID: "hq", WorkspaceType: auth.WorkspaceTypeHeadquarters,
		Permissions: map[string]struct{}{"paymentConfig:read": {}, "paymentConfig:manage": {}}}
}

func readPaymentRow(t *testing.T, db *gorm.DB, table string) configRow {
	t.Helper()
	var row configRow
	if err := db.Table(table).Take(&row).Error; err != nil {
		t.Fatal(err)
	}
	return row
}

func testWechatCredentials(t *testing.T) *WechatCredentials {
	t.Helper()
	private, cert := testCertificate(t)
	return &WechatCredentials{MerchantAPISerial: "2A", MerchantAPICert: cert, APIV3Key: strings.Repeat("k", 32), MerchantPrivateKey: private, VerificationMode: "PLATFORM_CERT"}
}

func testCertificate(t *testing.T) (string, string) {
	t.Helper()
	return testCertificateWithValidity(t, time.Now().Add(-time.Hour), time.Now().Add(time.Hour))
}

func testCertificateWithValidity(t *testing.T, start, end time.Time) (string, string) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	template := &x509.Certificate{SerialNumber: big.NewInt(42), Subject: pkix.Name{CommonName: "test"},
		NotBefore: start, NotAfter: end, KeyUsage: x509.KeyUsageDigitalSignature}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	private := string(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}))
	cert := string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}))
	return private, cert
}
