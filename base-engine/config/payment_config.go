package config

import (
	"encoding/base64"
	"errors"
	"os"
	"strings"
)

var ErrPaymentConfigSecurity = errors.New("invalid payment configuration keyring")

type PaymentConfigSecurity struct {
	ActiveKeyID string
	Keys        map[string][]byte
}

func LoadPaymentConfigSecurity() (PaymentConfigSecurity, error) {
	result := PaymentConfigSecurity{Keys: make(map[string][]byte)}
	for _, entry := range strings.Split(os.Getenv("PAYMENT_CONFIG_ENCRYPTION_KEYS"), ",") {
		id, encoded, ok := strings.Cut(strings.TrimSpace(entry), ":")
		key, err := base64.StdEncoding.DecodeString(encoded)
		if !ok || id == "" || len(id) > 64 || err != nil || len(key) != 32 {
			return PaymentConfigSecurity{}, ErrPaymentConfigSecurity
		}
		if _, duplicate := result.Keys[id]; duplicate {
			return PaymentConfigSecurity{}, ErrPaymentConfigSecurity
		}
		if result.ActiveKeyID == "" {
			result.ActiveKeyID = id
		}
		result.Keys[id] = key
	}
	if result.ActiveKeyID == "" {
		return PaymentConfigSecurity{}, ErrPaymentConfigSecurity
	}
	return result, nil
}
