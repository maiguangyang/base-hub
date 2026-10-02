package paymentconfig

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"

	"base-engine/config"
)

func encryptCredentials(keys config.PaymentConfigSecurity, plaintext []byte) (string, string, error) {
	key := keys.Keys[keys.ActiveKeyID]
	if len(key) != 32 || keys.ActiveKeyID == "" {
		return "", "", ErrUnavailable
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", ErrUnavailable
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", ErrUnavailable
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", "", ErrUnavailable
	}
	return keys.ActiveKeyID, base64.StdEncoding.EncodeToString(aead.Seal(nonce, nonce, plaintext, nil)), nil
}

func decryptCredentials(keys config.PaymentConfigSecurity, keyID, ciphertext string) ([]byte, error) {
	key := keys.Keys[keyID]
	if len(key) != 32 {
		return nil, ErrUnavailable
	}
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, ErrUnavailable
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, ErrUnavailable
	}
	aead, err := cipher.NewGCM(block)
	if err != nil || len(data) < aead.NonceSize() {
		return nil, ErrUnavailable
	}
	result, err := aead.Open(nil, data[:aead.NonceSize()], data[aead.NonceSize():], nil)
	if err != nil {
		return nil, ErrUnavailable
	}
	return result, nil
}
