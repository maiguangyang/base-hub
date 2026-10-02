package productcatalog

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
)

// GeneratePackageBarcode returns an unreserved candidate; publication enforces uniqueness.
func (s *Service) GeneratePackageBarcode(ctx context.Context, principal *auth.WorkspacePrincipal) (string, error) {
	_, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessCreate)
	if auth.ErrorCode(err) == auth.CodePermissionDenied {
		_, err = headquartersID(principal, "hqProductCatalog:manage", authorization.AccessUpdate)
	}
	if err != nil {
		return "", err
	}
	return generatePackageBarcode(ctx, s.db, rand.Reader)
}

func generatePackageBarcode(ctx context.Context, db *gorm.DB, entropy io.Reader) (string, error) {
	for attempt := 0; attempt < 8; attempt++ {
		number, err := rand.Int(entropy, big.NewInt(1000000000000))
		if err != nil {
			return "", auth.NewError(auth.CodeInternalError)
		}
		code := fmt.Sprintf("KH%012d", number.Int64())
		var count int64
		if err := db.WithContext(ctx).Model(&gen.ProductPackage{}).Where("barcode = ?", code).Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			return code, nil
		}
	}
	return "", auth.NewError(auth.CodeConflict)
}
