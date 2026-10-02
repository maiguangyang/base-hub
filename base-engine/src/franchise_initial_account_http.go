package src

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"base-engine/auth"
	"base-engine/config"
	"base-engine/src/services/ai"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
)

const initialAccountBodyLimit = 2048

type initialAccountHTTPInput struct {
	OrganizationID    string `json:"organizationId"`
	AccountID         string `json:"accountId"`
	EvidenceReference string `json:"evidenceReference"`
	Attested          bool   `json:"attested"`
}

type initialAccountHTTPResult struct {
	OrganizationID string `json:"organizationId"`
	AccountID      string `json:"accountId"`
}

// RegisterFranchiseInitialAccountRoute exposes the evidence-backed command without changing generated GraphQL code.
func RegisterFranchiseInitialAccountRoute(router *mux.Router, db *gorm.DB, security config.SecurityConfig) {
	router.HandleFunc("/api/franchise-initial-account", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Cache-Control", "no-store")
		origin := strings.TrimSpace(request.Header.Get("Origin"))
		if origin == "" || !security.OriginAllowed(origin) {
			writeInitialAccountError(writer, http.StatusForbidden, auth.CodePermissionDenied)
			return
		}
		if !isJSONContentType(request.Header.Get("Content-Type")) {
			writeInitialAccountError(writer, http.StatusUnsupportedMediaType, auth.CodeValidationFailed)
			return
		}
		principal, err := auth.RequirePrincipal(request.Context())
		if err != nil {
			writeInitialAccountError(writer, http.StatusUnauthorized, auth.ErrorCode(err))
			return
		}
		input, err := decodeInitialAccountInput(writer, request)
		if err != nil || !input.Attested {
			writeInitialAccountError(writer, http.StatusBadRequest, auth.CodeValidationFailed)
			return
		}
		if err := authorization.SetFranchiseInitialAccount(db.WithContext(request.Context()), principal, input.OrganizationID, input.AccountID, input.EvidenceReference); err != nil {
			status, code := initialAccountHTTPError(err)
			writeInitialAccountError(writer, status, code)
			return
		}
		writeInitializationJSON(writer, http.StatusOK, initialAccountHTTPResult{OrganizationID: input.OrganizationID, AccountID: input.AccountID})
	}).Methods(http.MethodPost).Name(ai.HTTPRouteName("http.franchiseInitialAccount.set", initialAccountHTTPInput{}, initialAccountHTTPResult{}))
}

func initialAccountHTTPError(err error) (int, auth.Code) {
	code := auth.ErrorCode(err)
	switch code {
	case auth.CodeValidationFailed:
		return http.StatusBadRequest, code
	case auth.CodeConflict, auth.CodeOpeningRecordNumberConflict:
		return http.StatusConflict, code
	case "", auth.CodeInternalError:
		return http.StatusInternalServerError, auth.CodeInternalError
	default:
		return http.StatusForbidden, code
	}
}

func decodeInitialAccountInput(writer http.ResponseWriter, request *http.Request) (initialAccountHTTPInput, error) {
	input := initialAccountHTTPInput{}
	decoder := json.NewDecoder(http.MaxBytesReader(writer, request.Body, initialAccountBodyLimit))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		return input, err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return input, errors.New("INVALID_JSON")
	}
	return input, nil
}

func writeInitialAccountError(writer http.ResponseWriter, status int, code auth.Code) {
	writeInitializationJSON(writer, status, map[string]auth.Code{"code": code})
}
