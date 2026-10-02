package store

import (
	"bytes"
	"context"
	_ "golang.org/x/image/webp"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/src/services/authorization"
)

const maxStoreDocumentBytes = 5 << 20

type pendingDocument struct {
	accountID, sessionID, organizationID, extension string
	createdAt                                       time.Time
}

func (s *Service) SetDocumentRoot(root string) { s.documentRoot = root }

func (s *Service) UploadDocument(_ context.Context, p *auth.WorkspacePrincipal, data []byte) (string, error) {
	orgID, err := documentUploadOrganization(p)
	if err != nil {
		return "", err
	}
	ext, err := storeDocumentExtension(data)
	if err != nil {
		return "", err
	}
	s.documentMu.Lock()
	defer s.documentMu.Unlock()
	s.sweepPendingDocuments(time.Now())
	id := uuid.Must(uuid.NewV4()).String()
	dir := filepath.Join(s.documentRoot, "store-documents", "pending")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	path := filepath.Join(dir, id+"."+ext)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return "", err
	}
	s.pendingDocuments[id] = pendingDocument{accountID: p.AccountID, sessionID: p.SessionID, organizationID: orgID, extension: ext, createdAt: time.Now()}
	return id, nil
}

func documentUploadOrganization(p *auth.WorkspacePrincipal) (string, error) {
	if p == nil || p.OrganizationID == nil {
		return "", auth.NewError(auth.CodeAuthRequired)
	}
	type candidate struct {
		action string
		mode   authorization.AccessMode
	}
	var actions []candidate
	switch p.WorkspaceType {
	case auth.WorkspaceTypeHeadquarters:
		actions = []candidate{{"hqStore:create", authorization.AccessCreate}, {"hqStore:update", authorization.AccessUpdate}}
	case auth.WorkspaceTypeFranchise:
		actions = []candidate{{"store:create", authorization.AccessCreate}, {"store:update", authorization.AccessUpdate}}
	default:
		return "", auth.NewError(auth.CodeWorkspaceForbidden)
	}
	for _, action := range actions {
		if !p.Has(action.action) {
			continue
		}
		if err := authorization.Authorize(p, authorization.Intent{Action: action.action, Mode: action.mode, ResourceOrganizationID: p.OrganizationID}); err == nil {
			return *p.OrganizationID, nil
		}
	}
	return "", auth.NewError(auth.CodePermissionDenied)
}

func storeDocumentExtension(data []byte) (string, error) {
	if len(data) == 0 || len(data) > maxStoreDocumentBytes {
		return "", auth.NewError(auth.CodeValidationFailed)
	}
	config, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil || !validDocumentDimensions(config.Width, config.Height) {
		return "", auth.NewError(auth.CodeValidationFailed)
	}
	if !validDocumentFormat(format) {
		return "", auth.NewError(auth.CodeValidationFailed)
	}
	if _, decoded, err := image.Decode(bytes.NewReader(data)); err != nil || decoded != format {
		return "", auth.NewError(auth.CodeValidationFailed)
	}
	if format == "jpeg" {
		return "jpg", nil
	}
	return format, nil
}

func validDocumentDimensions(width, height int) bool {
	return width >= 1 && height >= 1 && width <= 4096 && height <= 4096
}

func validDocumentFormat(format string) bool {
	return format == "jpeg" || format == "png" || format == "webp"
}

func (s *Service) pendingDocumentPath(id, ext string) string {
	return filepath.Join(s.documentRoot, "store-documents", "pending", id+"."+ext)
}

func (s *Service) storeDocumentPath(storeID, filename string) string {
	return filepath.Join(s.documentRoot, "store-documents", "stores", storeID, filename)
}

func (s *Service) sweepPendingDocuments(now time.Time) {
	for id, item := range s.pendingDocuments {
		if now.Sub(item.createdAt) <= time.Hour {
			continue
		}
		_ = os.Remove(s.pendingDocumentPath(id, item.extension))
		delete(s.pendingDocuments, id)
	}
}

func validDocumentFilename(value string) bool {
	if strings.ContainsAny(value, "/\\") {
		return false
	}
	ext := filepath.Ext(value)
	if ext != ".jpg" && ext != ".png" && ext != ".webp" {
		return false
	}
	_, err := uuid.FromString(strings.TrimSuffix(value, ext))
	return err == nil
}

func validDocumentStoreID(value string) bool {
	_, err := uuid.FromString(value)
	return err == nil || (value != "" && !strings.ContainsAny(value, "/\\."))
}
