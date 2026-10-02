package ai

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"

	"base-engine/config"
)

var ErrModelConfigUnavailable = errors.New("AI model configuration unavailable")

func encryptModelKey(security config.AIModelSecurityConfig, plaintext string) ([]byte, error) {
	key := security.Keys[security.ActiveKeyID]
	if len(key) != 32 {
		return nil, ErrModelConfigUnavailable
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, ErrModelConfigUnavailable
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, ErrModelConfigUnavailable
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, ErrModelConfigUnavailable
	}
	return gcm.Seal(nonce, nonce, []byte(plaintext), nil), nil
}

func decryptModelKey(security config.AIModelSecurityConfig, keyID string, ciphertext []byte) (string, error) {
	key := security.Keys[keyID]
	if len(key) != 32 {
		return "", ErrModelConfigUnavailable
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", ErrModelConfigUnavailable
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil || len(ciphertext) < gcm.NonceSize() {
		return "", ErrModelConfigUnavailable
	}
	value, err := gcm.Open(nil, ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():], nil)
	if err != nil {
		return "", ErrModelConfigUnavailable
	}
	return string(value), nil
}
