package paymentconfig

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"strings"
)

func parseAlipayPrivateKey(raw string) (*rsa.PrivateKey, error) {
	if strings.Contains(raw, "-----BEGIN") {
		return parsePrivateKey(raw)
	}
	der, err := base64.StdEncoding.DecodeString(strings.Join(strings.Fields(raw), ""))
	if err != nil {
		return nil, ErrInvalid
	}
	if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
		return key, nil
	}
	value, err := x509.ParsePKCS8PrivateKey(der)
	if err != nil {
		return nil, ErrInvalid
	}
	key, ok := value.(*rsa.PrivateKey)
	if !ok {
		return nil, ErrInvalid
	}
	return key, nil
}

func parseAlipayPublicKey(raw string) (*rsa.PublicKey, error) {
	if strings.Contains(raw, "-----BEGIN") {
		return parsePublicKey(raw)
	}
	der, err := base64.StdEncoding.DecodeString(strings.Join(strings.Fields(raw), ""))
	if err != nil {
		return nil, ErrInvalid
	}
	if value, err := x509.ParsePKIXPublicKey(der); err == nil {
		if key, ok := value.(*rsa.PublicKey); ok {
			return key, nil
		}
	}
	key, err := x509.ParsePKCS1PublicKey(der)
	if err != nil {
		return nil, ErrInvalid
	}
	return key, nil
}
