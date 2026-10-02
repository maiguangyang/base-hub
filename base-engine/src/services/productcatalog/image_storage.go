package productcatalog

import (
	"bytes"
	"context"
	_ "golang.org/x/image/webp"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/src/services/authorization"
)

const maxProductImageBytes = 5 << 20

var safeProductIDPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{1,128}$`)

type pendingImage struct {
	accountID, sessionID, organizationID, extension string
	createdAt                                       time.Time
}

func (s *Service) SetImageRoot(root string) { s.imageRoot = root }
func (s *Service) ImageRoot() string        { return s.imageRoot }

func (s *Service) UploadMainImage(ctx context.Context, principal *auth.WorkspacePrincipal, data []byte) (string, error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessCreate)
	if err != nil {
		return "", err
	}
	ext, err := imageExtension(data)
	if err != nil {
		return "", err
	}
	id := uuid.Must(uuid.NewV4()).String()
	s.imageMu.Lock()
	defer s.imageMu.Unlock()
	now := time.Now()
	s.sweepPendingImages(now)
	if now.Sub(s.lastImageSweep) >= time.Hour {
		s.lastImageSweep = now
		_ = s.sweepDiskImages(ctx, now)
	}
	if err := os.MkdirAll(filepath.Join(s.imageRoot, "pending"), 0700); err != nil {
		return "", err
	}
	file, err := os.OpenFile(s.pendingPath(id, ext), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return "", err
	}
	if _, err := file.Write(data); err != nil {
		_ = file.Close()
		_ = os.Remove(s.pendingPath(id, ext))
		return "", err
	}
	if err := file.Close(); err != nil {
		_ = os.Remove(s.pendingPath(id, ext))
		return "", err
	}
	s.pendingImages[id] = pendingImage{accountID: principal.AccountID, sessionID: principal.SessionID,
		organizationID: hqID, extension: ext, createdAt: time.Now()}
	return id, nil
}

func imageExtension(data []byte) (string, error) {
	if len(data) == 0 || len(data) > maxProductImageBytes {
		return "", auth.NewError(auth.CodeValidationFailed)
	}
	format := detectImageFormat(data)
	if format == "" {
		return "", auth.NewError(auth.CodeValidationFailed)
	}
	if err := validateDecodedImage(data, format); err != nil {
		return "", err
	}
	if format == "jpeg" {
		return "jpg", nil
	}
	return format, nil
}

func detectImageFormat(data []byte) string {
	switch {
	case bytes.HasPrefix(data, []byte{0xff, 0xd8, 0xff}):
		return "jpeg"
	case bytes.HasPrefix(data, []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}):
		return "png"
	case len(data) >= 12 && string(data[:4]) == "RIFF" && string(data[8:12]) == "WEBP":
		return "webp"
	default:
		return ""
	}
}

func validateDecodedImage(data []byte, format string) error {
	config, decoded, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil || decoded != format || config.Width < 1 || config.Height < 1 || config.Width > 4096 || config.Height > 4096 {
		return auth.NewError(auth.CodeValidationFailed)
	}
	if _, decoded, err = image.Decode(bytes.NewReader(data)); err != nil || decoded != format {
		return auth.NewError(auth.CodeValidationFailed)
	}
	return nil
}

func (s *Service) pendingPath(id, ext string) string {
	return filepath.Join(s.imageRoot, "pending", id+"."+ext)
}

func (s *Service) sweepPendingImages(now time.Time) {
	for id, pending := range s.pendingImages {
		if now.Sub(pending.createdAt) <= time.Hour {
			continue
		}
		_ = os.Remove(s.pendingPath(id, pending.extension))
		delete(s.pendingImages, id)
	}
}

func (s *Service) ValidateMainImageAttachment(ctx context.Context, principal *auth.WorkspacePrincipal, id string) error {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessUpdate)
	if err != nil {
		return err
	}
	s.imageMu.Lock()
	defer s.imageMu.Unlock()
	pending, ok := s.pendingImages[id]
	if !ok || time.Since(pending.createdAt) > time.Hour {
		return auth.NewError(auth.CodeValidationFailed)
	}
	if pending.accountID != principal.AccountID || pending.sessionID != principal.SessionID || pending.organizationID != hqID {
		return auth.NewError(auth.CodePermissionDenied)
	}
	if _, err := os.Stat(s.pendingPath(id, pending.extension)); err != nil {
		return auth.NewError(auth.CodeValidationFailed)
	}
	return nil
}

func managedImagePath(root, productID, imageURL string) (string, bool) {
	if !safeProductPathID(productID) {
		return "", false
	}
	prefix := "/uploads/products/" + productID + "/"
	if !strings.HasPrefix(imageURL, prefix) {
		return "", false
	}
	name := strings.TrimPrefix(imageURL, prefix)
	if strings.ContainsAny(name, "/\\") || name == "" {
		return "", false
	}
	base := strings.TrimSuffix(strings.TrimSuffix(strings.TrimSuffix(name, ".png"), ".jpg"), ".webp")
	if base == name {
		return "", false
	}
	if _, err := uuid.FromString(base); err != nil {
		return "", false
	}
	return filepath.Join(root, "products", productID, name), true
}

func safeProductPathID(value string) bool {
	return safeProductIDPattern.MatchString(value)
}
