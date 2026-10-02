package paymentconfig

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"math/big"
	"strings"
	"time"
)

func validMerchant(channel, id, environment string) bool {
	length := len(id)
	if channel == "WECHAT" {
		if environment != "" || length < 6 || length > 32 {
			return false
		}
	} else if channel == "ALIPAY" {
		if (environment != "PRODUCTION" && environment != "SANDBOX") || length != 16 {
			return false
		}
	} else {
		return false
	}
	return digitsOnly(id)
}

func digitsOnly(id string) bool {
	for _, char := range id {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

func validateAndEncode(input SaveInput) ([]byte, error) {
	switch input.Channel {
	case "WECHAT":
		if input.WechatCredentials == nil || input.AlipayCredentials != nil {
			return nil, ErrInvalid
		}
		if err := validateWechat(*input.WechatCredentials); err != nil {
			return nil, err
		}
		return encodeCredential(input.WechatCredentials)
	case "ALIPAY":
		if input.AlipayCredentials == nil || input.WechatCredentials != nil {
			return nil, ErrInvalid
		}
		if err := validateAlipay(*input.AlipayCredentials); err != nil {
			return nil, err
		}
		return encodeCredential(input.AlipayCredentials)
	default:
		return nil, ErrInvalid
	}
}

func validateStoredCredential(channel string, data []byte) error {
	switch channel {
	case "WECHAT":
		var credential WechatCredentials
		if err := json.Unmarshal(data, &credential); err != nil {
			return ErrInvalid
		}
		return validateWechat(credential)
	case "ALIPAY":
		var credential AlipayCredentials
		if err := json.Unmarshal(data, &credential); err != nil {
			return ErrInvalid
		}
		if credential.AlipayPublicKey != "" {
			return validateAlipay(credential)
		}
		var legacy legacyAlipayCredentials
		if err := json.Unmarshal(data, &legacy); err != nil {
			return ErrInvalid
		}
		return validateLegacyAlipay(legacy)
	default:
		return ErrInvalid
	}
}

func validateWechat(c WechatCredentials) error {
	if len(c.APIV3Key) != 32 {
		return ErrInvalid
	}
	certificate, err := parseCertificate(c.MerchantAPICert)
	if err != nil {
		return err
	}
	private, err := parsePrivateKey(c.MerchantPrivateKey)
	if err != nil || !matchingRSA(private, certificate.PublicKey) {
		return ErrInvalid
	}
	serial := new(big.Int)
	if _, ok := serial.SetString(strings.TrimSpace(c.MerchantAPISerial), 16); !ok || serial.Cmp(certificate.SerialNumber) != 0 {
		return ErrInvalid
	}
	return validateWechatVerification(c)
}

func validateWechatVerification(c WechatCredentials) error {
	switch c.VerificationMode {
	case "PLATFORM_CERT":
		if c.WechatPublicKeyID != "" || c.WechatPublicKey != "" {
			return ErrInvalid
		}
	case "PUBLIC_KEY":
		if !strings.HasPrefix(c.WechatPublicKeyID, "PUB_KEY_ID_") || len(c.WechatPublicKeyID) <= len("PUB_KEY_ID_") {
			return ErrInvalid
		}
		if _, err := parsePublicKey(c.WechatPublicKey); err != nil {
			return err
		}
	default:
		return ErrInvalid
	}
	return nil
}

func validateAlipay(c AlipayCredentials) error {
	if _, err := parseAlipayPrivateKey(c.AppPrivateKey); err != nil {
		return ErrInvalid
	}
	if _, err := parseAlipayPublicKey(c.AlipayPublicKey); err != nil {
		return ErrInvalid
	}
	return nil
}

type legacyAlipayCredentials struct {
	AppPrivateKey    string `json:"appPrivateKey"`
	AppPublicCert    string `json:"appPublicCert"`
	AlipayRootCert   string `json:"alipayRootCert"`
	AlipayPublicCert string `json:"alipayPublicCert"`
}

func validateLegacyAlipay(c legacyAlipayCredentials) error {
	appCert, err := parseCertificate(c.AppPublicCert)
	if err != nil {
		return err
	}
	private, err := parsePrivateKey(c.AppPrivateKey)
	if err != nil || !matchingRSA(private, appCert.PublicKey) {
		return ErrInvalid
	}
	if _, err := parseCertificate(c.AlipayPublicCert); err != nil {
		return err
	}
	return validateCertificateBundle(c.AlipayRootCert)
}

func validateCertificateBundle(raw string) error {
	remaining := []byte(raw)
	count := 0
	for len(strings.TrimSpace(string(remaining))) > 0 {
		block, rest := pem.Decode(remaining)
		if block == nil || block.Type != "CERTIFICATE" {
			return ErrInvalid
		}
		certificate, err := x509.ParseCertificate(block.Bytes)
		if err != nil || !validCertificateDate(certificate) {
			return ErrInvalid
		}
		count++
		remaining = rest
	}
	if count == 0 {
		return ErrInvalid
	}
	return nil
}

func parseCertificate(raw string) (*x509.Certificate, error) {
	block, rest := pem.Decode([]byte(raw))
	if block == nil || block.Type != "CERTIFICATE" || strings.TrimSpace(string(rest)) != "" {
		return nil, ErrInvalid
	}
	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil || !validCertificateDate(certificate) {
		return nil, ErrInvalid
	}
	return certificate, nil
}

func validCertificateDate(c *x509.Certificate) bool {
	if c == nil {
		return false
	}
	now := time.Now()
	return !now.Before(c.NotBefore) && !now.After(c.NotAfter)
}

func parsePrivateKey(raw string) (*rsa.PrivateKey, error) {
	block, rest := pem.Decode([]byte(raw))
	if block == nil || strings.TrimSpace(string(rest)) != "" {
		return nil, ErrInvalid
	}
	if block.Type == "RSA PRIVATE KEY" {
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err == nil {
			return key, nil
		}
	}
	if block.Type == "PRIVATE KEY" {
		value, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err == nil {
			if key, ok := value.(*rsa.PrivateKey); ok {
				return key, nil
			}
		}
	}
	return nil, ErrInvalid
}

func parsePublicKey(raw string) (*rsa.PublicKey, error) {
	block, rest := pem.Decode([]byte(raw))
	if block == nil || strings.TrimSpace(string(rest)) != "" {
		return nil, ErrInvalid
	}
	if block.Type == "PUBLIC KEY" {
		value, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err == nil {
			if key, ok := value.(*rsa.PublicKey); ok {
				return key, nil
			}
		}
	}
	if block.Type == "RSA PUBLIC KEY" {
		key, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err == nil {
			return key, nil
		}
	}
	return nil, ErrInvalid
}

func matchingRSA(private *rsa.PrivateKey, public any) bool {
	if private == nil {
		return false
	}
	key, ok := public.(*rsa.PublicKey)
	return ok && private.PublicKey.E == key.E && private.PublicKey.N.Cmp(key.N) == 0
}
