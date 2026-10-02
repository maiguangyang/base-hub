package ai

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
)

// injectTrustedFields runs only after the approved business arguments are consumed.
// A stable plan-bound value lets the existing business idempotency checks recognize
// an uncertain result without allowing the model to choose or reuse the key.
func injectTrustedFields(spec ToolSpec, planID string, raw json.RawMessage) (json.RawMessage, error) {
	if spec.RequestKeyPath == "" && spec.GeneratedCodePath == "" {
		return raw, nil
	}
	if planID == "" || spec.Mode != ModeWrite {
		return nil, errors.New("TRUSTED_FIELD_PLAN_REQUIRED")
	}
	var args map[string]any
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	if err := decoder.Decode(&args); err != nil || args == nil {
		return nil, errors.New("INVALID_APPROVED_ARGUMENTS")
	}
	hash := sha256.Sum256(append([]byte(planID+"\x00"+spec.ID+"\x00"), raw...))
	value := hex.EncodeToString(hash[:])
	if err := applyTrustedFields(args, spec, value); err != nil {
		return nil, err
	}
	return json.Marshal(args)
}

func applyTrustedFields(args map[string]any, spec ToolSpec, value string) error {
	if spec.RequestKeyPath != "" {
		if err := setTrustedField(args, spec.RequestKeyPath, value); err != nil {
			return err
		}
	}
	if spec.GeneratedCodePath != "" {
		if err := setTrustedField(args, spec.GeneratedCodePath, "CC"+value[:24]); err != nil {
			return err
		}
	}
	return nil
}

func setTrustedField(args map[string]any, path, value string) error {
	parts := strings.Split(path, ".")
	for _, part := range parts[:len(parts)-1] {
		child, ok := args[part].(map[string]any)
		if !ok {
			return errors.New("INVALID_TRUSTED_FIELD_PATH")
		}
		args = child
	}
	key := parts[len(parts)-1]
	if key == "" {
		return errors.New("INVALID_TRUSTED_FIELD_PATH")
	}
	if _, exists := args[key]; exists {
		return errors.New("MODEL_TRUSTED_FIELD_FORBIDDEN")
	}
	args[key] = value
	return nil
}
