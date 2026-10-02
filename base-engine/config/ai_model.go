package config

import (
	"encoding/base64"
	"errors"
	"os"
	"strings"
)

var ErrAIModelSecurityConfig = errors.New("invalid AI model security configuration")

type AIModelSecurityConfig struct {
	ActiveKeyID string
	Keys        map[string][]byte
}

func LoadAIModelSecurityConfig() (AIModelSecurityConfig, error) {
	result := AIModelSecurityConfig{Keys: make(map[string][]byte)}
	for _, entry := range strings.Split(os.Getenv("AI_MODEL_ENCRYPTION_KEYS"), ",") {
		if strings.TrimSpace(entry) == "" {
			continue
		}
		keyID, key, err := parseAIModelKey(entry)
		if err != nil {
			return AIModelSecurityConfig{}, err
		}
		if _, exists := result.Keys[keyID]; exists {
			return AIModelSecurityConfig{}, ErrAIModelSecurityConfig
		}
		if result.ActiveKeyID == "" {
			result.ActiveKeyID = keyID
		}
		result.Keys[keyID] = key
	}
	return result, nil
}

func parseAIModelKey(entry string) (string, []byte, error) {
	keyID, encoded, found := strings.Cut(strings.TrimSpace(entry), ":")
	key, err := base64.StdEncoding.DecodeString(encoded)
	if !found || keyID == "" || err != nil || len(key) != 32 {
		return "", nil, ErrAIModelSecurityConfig
	}
	return keyID, key, nil
}
