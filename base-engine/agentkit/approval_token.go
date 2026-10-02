package agentkit

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)

type approvalToken struct {
	ID      string `json:"id"`
	Expires int64  `json:"exp"`
	Digest  string `json:"digest"`
}

func (store *ApprovalStore[T]) parseToken(token string) (approvalToken, error) {
	if store == nil || len(token) > 4096 {
		return approvalToken{}, errors.New("INVALID_TOKEN")
	}
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return approvalToken{}, errors.New("INVALID_TOKEN")
	}
	provided, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil || !hmac.Equal(provided, store.sign(parts[0])) {
		return approvalToken{}, errors.New("INVALID_TOKEN_SIGNATURE")
	}
	data, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return approvalToken{}, errors.New("INVALID_TOKEN")
	}
	var payload approvalToken
	if json.Unmarshal(data, &payload) != nil || payload.ID == "" || payload.Digest == "" {
		return approvalToken{}, errors.New("INVALID_TOKEN")
	}
	return payload, nil
}

func (store *ApprovalStore[T]) sign(value string) []byte {
	mac := hmac.New(sha256.New, store.key)
	_, _ = mac.Write([]byte(value))
	return mac.Sum(nil)
}
