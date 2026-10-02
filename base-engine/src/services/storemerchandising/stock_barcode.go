package storemerchandising

import (
	"context"
	"errors"
	"math"
	"strings"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
)

func (s *Service) StockPackageByBarcode(ctx context.Context, principal *auth.WorkspacePrincipal, storeID, barcode string) (*gen.FranchiseStockPackageRecognition, error) {
	var result *gen.FranchiseStockPackageRecognition
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if _, err := storeScope(tx, principal, storeID, "franchiseStock:read", authorization.AccessRead); err != nil {
			return err
		}
		barcode = strings.TrimSpace(barcode)
		if barcode == "" || len(barcode) > 64 {
			return auth.NewError(auth.CodeValidationFailed)
		}
		var matches []gen.ProductPackage
		if err := tx.Where("active_barcode = ?", barcode).Limit(2).Find(&matches).Error; err != nil {
			return err
		}
		switch len(matches) {
		case 0:
			result = barcodeRecognition(gen.FranchiseStockPackageRecognitionStatusNotFound)
		case 1:
			var err error
			result, err = recognizeStockPackage(tx, storeID, &matches[0])
			return err
		default:
			result = barcodeRecognition(gen.FranchiseStockPackageRecognitionStatusConflict)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
func barcodeRecognition(status gen.FranchiseStockPackageRecognitionStatus) *gen.FranchiseStockPackageRecognition {
	return &gen.FranchiseStockPackageRecognition{Status: status}
}
func recognizeStockPackage(tx *gorm.DB, storeID string, pack *gen.ProductPackage) (*gen.FranchiseStockPackageRecognition, error) {
	unavailable := barcodeRecognition(gen.FranchiseStockPackageRecognitionStatusUnavailable)
	var sku gen.ProductSku
	if err := recognitionSku(tx, pack.SkuID, &sku); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return unavailable, nil
		}
		return nil, err
	}
	if sku.PublishedPackageSetVersion != pack.PackageSetVersion || !pack.Enabled {
		return unavailable, nil
	}
	var product gen.Product
	if err := tx.First(&product, "id = ?", sku.ProductID).Error; err != nil {
		return nil, err
	}
	var listing gen.StoreListing
	err := tx.Where("store_id = ? AND sku_id = ?", storeID, sku.ID).First(&listing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return barcodeRecognition(gen.FranchiseStockPackageRecognitionStatusNotSelected), nil
	}
	if err != nil {
		return nil, err
	}
	if !listing.Enabled {
		return unavailable, nil
	}
	return stockBarcodeTarget(tx, &product, &sku, &listing, pack)
}
func stockBarcodeTarget(tx *gorm.DB, product *gen.Product, sku *gen.ProductSku, listing *gen.StoreListing, pack *gen.ProductPackage) (*gen.FranchiseStockPackageRecognition, error) {
	unavailable := barcodeRecognition(gen.FranchiseStockPackageRecognitionStatusUnavailable)
	chain, err := stockBarcodeChain(tx, pack)
	if err != nil {
		if auth.ErrorCode(err) == auth.CodeConflict {
			return unavailable, nil
		}
		return nil, err
	}
	return &gen.FranchiseStockPackageRecognition{Status: gen.FranchiseStockPackageRecognitionStatusReady,
		Target: &gen.FranchiseStockPackageTarget{ProductName: product.Name, SkuID: sku.ID, SkuName: sku.Name,
			ListingID: listing.ID, PackageID: pack.ID, PackageName: pack.Name, Barcode: *pack.Barcode, PackageSetVersion: int(pack.PackageSetVersion), ConversionChain: chain}}, nil
}
func stockBarcodeChain(tx *gorm.DB, pack *gen.ProductPackage) ([]*gen.FranchiseStockPackageConversion, error) {
	chain := []*gen.FranchiseStockPackageConversion{}
	seen := map[string]bool{}
	current := *pack
	for {
		if seen[current.ID] {
			return nil, auth.NewError(auth.CodeConflict)
		}
		seen[current.ID] = true
		step, err := stockBarcodeChainStep(&current, pack)
		if err != nil {
			return nil, err
		}
		chain = append(chain, step)
		if current.ContainsPackageID == nil {
			break
		}
		var next gen.ProductPackage
		if err := tx.First(&next, "id = ?", *current.ContainsPackageID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, auth.NewError(auth.CodeConflict)
			}
			return nil, err
		}
		current = next
	}
	return baseFirstBarcodeChain(chain)
}
func baseFirstBarcodeChain(chain []*gen.FranchiseStockPackageConversion) ([]*gen.FranchiseStockPackageConversion, error) {
	for left, right := 0, len(chain)-1; left < right; left, right = left+1, right-1 {
		chain[left], chain[right] = chain[right], chain[left]
	}
	factor := 1
	for _, item := range chain {
		if factor > math.MaxInt32/item.ContainsQuantity {
			return nil, auth.NewError(auth.CodeConflict)
		}
		factor *= item.ContainsQuantity
		item.BaseQuantity = factor
	}
	return chain, nil
}
func validateReceiptBarcode(tx *gorm.DB, pack *gen.ProductPackage, barcode *string) error {
	if barcode == nil {
		return nil
	}
	value := strings.TrimSpace(*barcode)
	if value == "" || len(value) > 64 {
		return auth.NewError(auth.CodeValidationFailed)
	}
	if pack.Barcode == nil || value != *pack.Barcode {
		return auth.NewError(auth.CodeConflict)
	}
	var sku gen.ProductSku
	if err := tx.First(&sku, "id = ?", pack.SkuID).Error; err != nil {
		return err
	}
	if sku.PublishedPackageSetVersion != pack.PackageSetVersion {
		return auth.NewError(auth.CodeConflict)
	}
	return nil
}

func stockBarcodeChainStep(current, pack *gen.ProductPackage) (*gen.FranchiseStockPackageConversion, error) {
	if current.SkuID != pack.SkuID || current.PackageSetVersion != pack.PackageSetVersion || !current.Enabled {
		return nil, auth.NewError(auth.CodeConflict)
	}
	quantity := int64(1)
	if current.ContainsPackageID != nil {
		if current.ContainsQuantity == nil || *current.ContainsQuantity < 2 || *current.ContainsQuantity > math.MaxInt32 {
			return nil, auth.NewError(auth.CodeConflict)
		}
		quantity = *current.ContainsQuantity
	}
	return &gen.FranchiseStockPackageConversion{PackageID: current.ID, Name: current.Name, ContainsQuantity: int(quantity)}, nil
}

func recognitionSku(tx *gorm.DB, id string, sku *gen.ProductSku) error {
	return tx.Joins("JOIN products ON products.id = product_skus.product_id").
		Joins("JOIN organizations ON organizations.id = products.organization_id").
		Where("product_skus.id = ? AND product_skus.enabled = ? AND products.enabled = ? AND organizations.type = ? AND organizations.status = ?",
			id, true, true, gen.OrganizationTypeHeadquarters, gen.OrganizationStatusActive).First(sku).Error
}
