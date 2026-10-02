package config

import (
	"bytes"
	"encoding/base64"
	"errors"
	"testing"
)

func TestAIModelSecurityConfigParsesVersionedKeys(t *testing.T) {
	oldKey := base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{7}, 32))
	newKey := base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{8}, 32))
	t.Setenv("AI_MODEL_ENCRYPTION_KEYS", "v2:"+newKey+",v1:"+oldKey)
	got, err := LoadAIModelSecurityConfig()
	if err != nil || got.ActiveKeyID != "v2" || len(got.Keys) != 2 {
		t.Fatalf("parsed security config = %+v, %v", got, err)
	}
}

func TestAIModelSecurityConfigRejectsMalformedKey(t *testing.T) {
	t.Setenv("AI_MODEL_ENCRYPTION_KEYS", "v1:invalid")
	if _, err := LoadAIModelSecurityConfig(); !errors.Is(err, ErrAIModelSecurityConfig) {
		t.Fatalf("malformed key error = %v", err)
	}
}

func TestAIModelSecurityConfigDoesNotRequireOriginPolicy(t *testing.T) {
	key := base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{7}, 32))
	t.Setenv("AI_MODEL_ENCRYPTION_KEYS", "v1:"+key)
	t.Setenv("AI_MODEL_ALLOWED_ORIGINS", "http://zsgw.sjdistributor.com:4000")
	got, err := LoadAIModelSecurityConfig()
	if err != nil || got.ActiveKeyID != "v1" {
		t.Fatalf("self-hosted model security config = %+v, %v", got, err)
	}
}
