package config

import (
	"encoding/base64"
	"strings"
	"testing"
)

func TestPaymentConfigSecurityKeyring(t *testing.T) {
	first := base64.StdEncoding.EncodeToString([]byte(strings.Repeat("a", 32)))
	second := base64.StdEncoding.EncodeToString([]byte(strings.Repeat("b", 32)))
	t.Setenv("PAYMENT_CONFIG_ENCRYPTION_KEYS", "current:"+first+",old:"+second)
	keys, err := LoadPaymentConfigSecurity()
	if err != nil || keys.ActiveKeyID != "current" || len(keys.Keys) != 2 {
		t.Fatalf("keyring: %#v %v", keys, err)
	}
	for _, value := range []string{"", "bad", "a:short", "a:" + first + ",a:" + second} {
		t.Setenv("PAYMENT_CONFIG_ENCRYPTION_KEYS", value)
		if _, err := LoadPaymentConfigSecurity(); err == nil {
			t.Fatalf("accepted malformed keyring %q", value)
		}
	}
}
