package ai

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const aiRequestBodyLimit = 32 * 1024

func decodeAIRequest(w http.ResponseWriter, r *http.Request, value any) error {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, aiRequestBodyLimit))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(value); err != nil {
		return errAIBadRequest
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errAIBadRequest
	}
	return nil
}
